package models

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"hssh/conf"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type terminalInfo struct {
	Session *ssh.Session
	exitMsg string
	stdout  io.Reader
	stdin   io.Writer
	stderr  io.Reader
}

// 解决当本地调整了终端大小后，远程终端毫无反应的问题
// 解决方案：启动一个 goroutine 在后台不断监听窗口改变事件，然后调用 WindowChange 即可
// PS：*ssh.Session 上有一个 WindowChange 方法，用于向远端发送窗口调整事件
func (t *terminalInfo) updateTerminalSize() {
	go func() {
		// 监听窗口变更事件
		sigwinchCh := make(chan os.Signal, 1)
		signal.Notify(sigwinchCh, syscall.SIGWINCH)

		fd := int(os.Stdin.Fd())
		termWidth, termHeight, err := terminal.GetSize(fd)
		if err != nil {
			fmt.Println(err)
		}

		for {
			select {
			// 堵塞读取
			case sigwinch := <-sigwinchCh:
				if sigwinch == nil {
					return
				}
				currTermWidth, currTermHeight, err := terminal.GetSize(fd)

				// 判断以下窗口尺寸是否有改变
				if currTermHeight == termHeight && currTermWidth == termWidth {
					continue
				}

				// 更新远端大小
				err = t.Session.WindowChange(currTermHeight, currTermWidth)
				if err != nil {
					fmt.Printf("Unable to send window-change reqest: %s.", err)
					continue
				}

				termWidth, termHeight = currTermWidth, currTermHeight
			}
		}
	}()
}

// 交互的 session
func (t *terminalInfo) interactiveSession() error {
	defer func() {
		if t.exitMsg == "" {
			_, _ = fmt.Fprintln(os.Stdout, "the connection was closed on the remote side on ", time.Now().Format(time.RFC822))
		} else {
			_, _ = fmt.Fprintln(os.Stdout, t.exitMsg)
		}
	}()

	// 拿到当前终端文件描述符
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, state)

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	// request pty
	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}
	err = t.Session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{})
	if err != nil {
		return nil
	}

	t.updateTerminalSize()

	// 解决Tmux 标题以及回显不换行的问题
	// 解决方案：启动一个异步的管道式复制行为，并且带有 buf 的发送
	t.stdin, err = t.Session.StdinPipe()
	if err != nil {
		return nil
	}
	t.stdout, err = t.Session.StdoutPipe()
	if err != nil {
		return nil
	}
	t.stderr, err = t.Session.StderrPipe()
	if err != nil {
		return nil
	}

	go io.Copy(os.Stderr, t.stderr)
	go io.Copy(os.Stdout, t.stdout)
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			if n > 0 {
				_, err = t.stdin.Write(buf[:n])
				if err != nil {
					fmt.Println(err)
					t.exitMsg = err.Error()
					return
				}
			}
		}
	}()

	err = t.Session.Shell()
	if err != nil {
		return err
	}
	err = t.Session.Wait()
	if err != nil {
		return err
	}
	return nil
}

// 询问用户输入密码
func getPass() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password: ")
	text, _ := reader.ReadString('\n')
	// 去除换行符
	text = strings.Replace(text, "\n", "", -1)
	return text, nil
}

// OpenSSH 运行一个 ssh 会话
func OpenSSH(c conf.Server, key conf.Key) {
	var (
		user      string
		port      int
		addr      string
		auth      []ssh.AuthMethod
		config      ssh.Config
		sshConfig *ssh.ClientConfig
		client    *ssh.Client
		err       error
	)

	if c.Host == "" {
		conf.CheckErr(errors.New("host not exist"))
	}

	if c.User == "" {
		user = "root"
	} else {
		user = c.User
	}

	if c.Port == 0 {
		port = 22
	} else {
		port = c.Port
	}

	auth = make([]ssh.AuthMethod, 0)
	if c.AuthMethod == "password" {
		if c.Password == "" {
			auth = append(auth, ssh.RetryableAuthMethod(ssh.PasswordCallback(getPass), 3))
		} else {
			auth = append(auth, ssh.Password(c.Password))
		}
	} else if c.AuthMethod == "key" {
		pemBytes, err := ioutil.ReadFile(key.Path)
		conf.CheckErr(err)

		var signer ssh.Signer
		if key.Passphrase == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(key.Passphrase))
		}
		conf.CheckErr(err)
		auth = append(auth, ssh.PublicKeys(signer))
	} else {
		conf.CheckErr(errors.New("auth method only supports password and key"))
	}

	config = ssh.Config{
		Ciphers: []string{
			"aes128-ctr",
			"aes192-ctr",
			"aes256-ctr",
			"aes128-gcm@openssh.com",
			"arcfour256",
			"arcfour128",
			"aes128-cbc",
			"3des-cbc",
			"aes192-cbc",
			"aes256-cbc",
		},
	}

	// 创建 ssh 配置
	sshConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Config:          config,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// 创建 client
	addr = fmt.Sprintf("%s:%d", c.Host, port)
	client, err = ssh.Dial("tcp", addr, sshConfig)
	conf.CheckErr(err)
	defer client.Close()

	// 获取 session
	err = newSession(client)
	conf.CheckErr(err)
}

// 创建一个新的交互式 session
func newSession(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	s := terminalInfo{
		Session: session,
	}
	return s.interactiveSession()
}
