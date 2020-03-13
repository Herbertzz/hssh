package models

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"hssh/conf"
	"os"
	"os/exec"
	"path/filepath"
)

// StartAPP 运行程序
func StartAPP() {
	app := &cli.App{
		Name:    conf.PROJECTNAME,
		Usage:   "manage ssh sessions",
		Version: conf.VERSION,
		Action:  actionOfOpenSSH,
		Commands: []*cli.Command{
			 commandOfAdd(),
			 commandOfRm(),
			 commandOfLs(),
			 commandOfEdit(),
			 commandOfUninstall(),
			 commandOfKeys(),
		},
	}

	err := app.Run(os.Args)
	CheckErr(err)
}

// 动作: 打开ssh会话
func actionOfOpenSSH(c *cli.Context) error {
	if c.Args().First() != "" {
		// 读取配置文件
		configs, err := conf.ReadYamlConfig()
		CheckErr(err)
		// 检查配置文件中服务器列表
		if len(configs.Servers) == 0 {
			fmt.Printf("list is empty, please execute command `%s add` first\n", conf.PROJECTNAME)
			return nil
		}
		session, ok := configs.Servers[c.Args().First()]
		if !ok {
			fmt.Printf("do not find session named: %s\n", c.Args().First())
			return nil
		}

		key := ""
		if session.AuthMethod == "key" {
			key, ok = configs.Keys[session.PrivateKey]
			if !ok {
				fmt.Printf("%s not exist in keys\n", session.PrivateKey)
				os.Exit(0)
			}
		}

		OpenSSH(session, key)
		return nil
	}
	fmt.Printf("please execute command `%s h` for help\n", conf.PROJECTNAME)
	return nil
}

// add 命令
func commandOfAdd() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "add a ssh session to the list",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Aliases:  []string{"i"},
				Usage:    "ip address or host",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "username",
				Aliases:     []string{"u"},
				Usage:       "username",
				DefaultText: "root",
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "`port`",
				DefaultText: "22",
			},
			&cli.StringFlag{
				Name:        "auth",
				Usage:       "auth `method`: password or key",
				DefaultText: "password",
				Value:       "password",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"pass"},
				Usage:   "password auth",
			},
			&cli.StringFlag{
				Name:    "private-key",
				Aliases: []string{"key"},
				Usage:   "The value of the Keys list",
			},
			&cli.StringFlag{
				Name:    "key-passphrase",
				Aliases: []string{"key-pass"},
				Usage:   "private key password",
			},
		},
		Action: actionOfAdd,
	}
}

// 动作: add
func actionOfAdd(c *cli.Context) error {
	arg := c.Args().First()
	// 检查alias是否存在
	if arg == "" {
		fmt.Println("alias is not set")
		return nil
	}
	// 读取配置文件
	configs, err := conf.ReadYamlConfig()
	CheckErr(err)
	sessions := configs.Servers
	if len(sessions) == 0 {
		sessions = make(map[string]conf.Server)
	}
	_, ok := sessions[arg]
	if ok {
		fmt.Printf("%s is already in the list\n", arg)
		return nil
	}
	// 生成服务器配置
	session := conf.Server{
		Host: c.String("host"),
	}
	if c.String("username") != "" {
		session.User = c.String("username")
	} else {
		session.User = "root"
	}
	if c.Int("port") != 0 {
		session.Port = c.Int("port")
	} else {
		session.Port = 22
	}
	// 认证方式
	authMethod := c.String("auth")
	if authMethod == "password" {
		session.Password = c.String("password")
	} else if authMethod == "key" {
		if c.String("private-key") == "" {
			session.PrivateKey = "default"
		} else {
			key := c.String("private-key")
			_, ok = configs.Keys[key]
			if !ok {
				fmt.Printf("%s does not exist in keys\n", key)
				return nil
			}
			session.PrivateKey = key
		}
		// 密钥密码
		if c.String("key-passphrase") != "" {
			session.KeyPassphrase = c.String("key-passphrase")
		}
	} else {
		fmt.Println("'--auth' only supports password and key")
		return nil
	}
	session.AuthMethod = authMethod

	sessions[arg] = session
	configs.Servers = sessions
	err = conf.WriteYamlConfig(configs)
	CheckErr(err)
	return nil
}

// rm 命令
func commandOfRm() *cli.Command {
	return &cli.Command{
		Name:  "rm",
		Usage: "remove a ssh session to the list",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "all",
				Usage: "delete all session",
			},
		},
		Action: actionOfRm,
	}
}

// 动作: rm
func actionOfRm(c *cli.Context) error {
	// 读取配置
	configs, err := conf.ReadYamlConfig()
	CheckErr(err)
	// 清空服务器列表
	if c.Bool("all") {
		configs.Servers = make(map[string]conf.Server, 0)
		err := conf.WriteYamlConfig(configs)
		CheckErr(err)
		return nil
	}
	// 检查alias参数是否存在
	arg := c.Args().First()
	if arg == "" {
		fmt.Println("alias argument not set")
		return nil
	}
	// 检查指定alias的服务器是否存在
	_, ok := configs.Servers[arg]
	if !ok {
		fmt.Printf("%s not found\n", arg)
		return nil
	}
	delete(configs.Servers, arg)

	err = conf.WriteYamlConfig(configs)
	CheckErr(err)
	fmt.Println("success")
	return nil
}

// ls 命令
func commandOfLs() *cli.Command {
	return &cli.Command{
		Name:   "ls",
		Usage:  "show session list",
		Action: actionOfLs,
	}
}

// 动作: ls
func actionOfLs(c *cli.Context) error {
	// 读取配置
	configs, err := conf.ReadYamlConfig()
	CheckErr(err)
	// 检查是否存在服务器配置
	sessions := configs.Servers
	if len(sessions) == 0 {
		fmt.Printf("list is empty, please execute command `%s add` first\n", conf.PROJECTNAME)
		return nil
	}
	// 打印服务器列表
	var authMethod string
	index := 1
	for k, v := range sessions {
		if v.AuthMethod == "password" {
			authMethod = "Password: " + v.Password
		} else if v.AuthMethod == "key" {
			authMethod = "Key: " + v.PrivateKey
		} else {
			authMethod = "undefined"
		}
		fmt.Printf("%02d. %s: %s@%s:%d(%s)\n", index, k, v.User, v.Host, v.Port, authMethod)
		index++
	}
	return nil
}

// edit 命令
func commandOfEdit() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "modify a ssh session to the list",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"i"},
				Usage:   "ip address or host",
			},
			&cli.StringFlag{
				Name:        "username",
				Aliases:     []string{"u"},
				Usage:       "username",
				DefaultText: "root",
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "`port`",
				DefaultText: "22",
			},
			&cli.StringFlag{
				Name:        "auth",
				Usage:       "auth `method`: password or key",
				DefaultText: "password",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"pass"},
				Usage:   "password auth",
			},
			&cli.StringFlag{
				Name:    "private-key",
				Aliases: []string{"key"},
				Usage:   "The value of the Keys list",
			},
			&cli.StringFlag{
				Name:    "key-passphrase",
				Aliases: []string{"key-pass"},
				Usage:   "private key password",
			},
		},
		Action: actionOfEdit,
	}
}

// 动作: edit
func actionOfEdit(c *cli.Context) error {
	arg := c.Args().First()
	if arg == "" {
		fmt.Println("alias is not set")
		return nil
	}
	// 读取配置
	configs, _ := conf.ReadYamlConfig()
	sessions := configs.Servers
	_, ok := sessions[arg]
	if !ok {
		fmt.Printf("%s is not exist in the list\n", arg)
		os.Exit(0)
	}
	// 生成服务器配置
	session := sessions[arg]
	if c.String("host") != "" {
		session.Host = c.String("host")
	}
	if c.String("username") != "" {
		session.User = c.String("username")
	}
	if c.Int("port") != 0 {
		session.Port = c.Int("port")
	}
	// 认证方式
	authMethod := c.String("auth")
	if authMethod == "password" {
		session.Password = c.String("password")
		session.AuthMethod = authMethod
		// 清除key认证的值
		session.PrivateKey = ""
		session.KeyPassphrase = ""
	} else if authMethod == "key" {
		if c.String("private-key") == "" {
			session.PrivateKey = "default"
		} else {
			key := c.String("private-key")
			_, ok = configs.Keys[key]
			if !ok {
				fmt.Printf("%s does not exist in keys\n", key)
				return nil
			}
			session.PrivateKey = key
		}
		// 密钥密码
		if c.String("key-passphrase") != "" {
			session.KeyPassphrase = c.String("key-passphrase")
		}
		session.AuthMethod = authMethod
		session.Password = ""
	}
	// 更新配置文件
	sessions[arg] = session
	configs.Servers = sessions
	err := conf.WriteYamlConfig(configs)
	CheckErr(err)
	return nil
}

// uninstall 命令
func commandOfUninstall() *cli.Command {
	return &cli.Command{
		Name:  "uninstall",
		Usage: "unistall the app",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "all",
				Usage: "Delete with the configuration file",
			},
		},
		Action: actionOfUninstall,
	}
}

// 动作: uninstall
func actionOfUninstall(c *cli.Context) error {
	if c.Bool("all") {
		// 删除配置文件
		err := conf.DelYamlFile()
		CheckErr(err)
	}
	// 获取当前执行程序的绝对路径并删除该程序
	file, err := exec.LookPath(os.Args[0])
	CheckErr(err)
	path, err := filepath.Abs(file)
	CheckErr(err)
	err = os.Remove(path)
	CheckErr(err)
	return nil
}

// keys 管理
func commandOfKeys() *cli.Command {
	return &cli.Command{
		Name:  "keys",
		Usage: "private keys manager",
		Action: func(c *cli.Context) error {
			// 读取配置文件
			configs, err := conf.ReadYamlConfig()
			CheckErr(err)

			// 显示key详情
			if arg := c.Args().First(); arg != "" {
				v, ok := configs.Keys[arg]
				if !ok {
					fmt.Printf("%s not found\n", arg)
					os.Exit(0)
				}
				fmt.Printf("%s: %s", arg, v)
				return nil
			}

			// 显示keys列表
			ShowKeys(configs)
			return nil
		},
		Subcommands: []*cli.Command{
			keysSubcommandsOfAdd(),
			keysSubcommandsOfEdit(),
			keysSubcommandsOfRm(),
		},
	}
}

// keys subcommand: add 添加
func keysSubcommandsOfAdd() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "add a key to the keys",
		Action: func(c *cli.Context) error {
			key := c.Args().Get(0)
			value := c.Args().Get(1)
			if key == "" {
				fmt.Printf("key is empty\n")
				os.Exit(0)
			}
			if value == "" {
				fmt.Printf("private key path is empty\n")
				os.Exit(0)
			}

			configs, err := conf.ReadYamlConfig()
			CheckErr(err)

			_, ok := configs.Keys[key]
			if ok {
				fmt.Printf("%s is already in key\n", key)
				fmt.Printf("Modify the key if necessary, Please execute: %s keys edit %s %s\n", conf.PROJECTNAME, key, value)
				os.Exit(0)
			}
			path, err := PrivateKeyPath(value)
			CheckErr(err)
			configs.Keys[key] = path
			err = conf.WriteYamlConfig(configs)
			CheckErr(err)
			// 显示keys列表
			ShowKeys(configs)
			return nil
		},
	}
}

// keys subcommand: rm 删除
func keysSubcommandsOfRm() *cli.Command {
	return &cli.Command{
		Name:  "rm",
		Usage: "remove a key to the keys",
		Action: func(c *cli.Context) error {
			if arg := c.Args().First(); arg != "" {
				configs, err := conf.ReadYamlConfig()
				CheckErr(err)

				_, ok := configs.Keys[arg]
				if !ok {
					fmt.Printf("%s not found\n", arg)
					os.Exit(0)
				}
				delete(configs.Keys, arg)
				err = conf.WriteYamlConfig(configs)
				CheckErr(err)
				ShowKeys(configs)
				return nil
			}
			fmt.Println("key not set")
			return nil
		},
	}
}

// keys subcommand: edit 编辑
func keysSubcommandsOfEdit() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "modify a key to the keys",
		Action: func(c *cli.Context) error {
			key := c.Args().Get(0)
			value := c.Args().Get(1)
			if key == "" {
				fmt.Printf("key is empty\n")
				os.Exit(0)
			}
			if value == "" {
				fmt.Printf("private key path is empty\n")
				os.Exit(0)
			}

			configs, err := conf.ReadYamlConfig()
			CheckErr(err)

			_, ok := configs.Keys[key]
			if !ok {
				fmt.Printf("add the key if necessary, Please execute: %s keys add %s %s\n", conf.PROJECTNAME, key, value)
				os.Exit(0)
			}

			path, err := PrivateKeyPath(value)
			CheckErr(err)
			configs.Keys[key] = path
			err = conf.WriteYamlConfig(configs)
			CheckErr(err)
			// 显示keys列表
			ShowKeys(configs)
			return nil
		},
	}
}
