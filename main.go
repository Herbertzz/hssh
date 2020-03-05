package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"hssh/common"
	"hssh/config"
	"hssh/ssh"
	"os"
)

func main() {
	app := &cli.App{
		Name:    common.ProjectName,
		Usage:   "manage ssh sessions",
		Version: "0.0.8-beta",
		Action: func(c *cli.Context) error {
			if c.Args().First() != "" {
				sessions, success := config.ReadYamlConfig()
				if !success {
					fmt.Printf("please execute command `%s h` for help\n", common.ProjectName)
					os.Exit(0)
				}
				session, ok := sessions[c.Args().First()]
				if !ok {
					fmt.Printf("do not find session named: %s\n", c.Args().First())
					os.Exit(0)
				}

				ssh.OpenSSH(session)
				return nil
			}
			fmt.Printf("please execute command `%s h` for help\n", common.ProjectName)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add a ssh session to the list",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "alias",
						Aliases:  []string{"a"},
						Usage:    "ssh config alias",
						Required: true,
					},
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
						Name:    "password",
						Aliases: []string{"pass"},
						Usage:   "password auth",
					},
					&cli.StringFlag{
						Name:    "private-key",
						Aliases: []string{"key"},
						Usage:   "sshkey auth: Non-absolute path, will join home path + private key `path`",
					},
					&cli.StringFlag{
						Name:    "key-passphrase",
						Aliases: []string{"key-pass"},
						Usage:   "private key password",
					},
				},
				Action: func(c *cli.Context) error {
					alias := c.String("alias")
					sessions, _ := config.ReadYamlConfig()
					_, ok := sessions[alias]
					if ok {
						fmt.Printf("%s is already in the list\n", alias)
						os.Exit(0)
					}

					session := ssh.Config{}
					session.Host = c.String("host")
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
					if c.String("password") != "" {
						session.Password = c.String("password")
					} else if c.String("private-key") != "" {
						session.PrivateKeyPath = common.PrivateKeyPath(c.String("private-key"))
						if c.String("key-passphrase") != "" {
							session.KeyPassphrase = c.String("key-passphrase")
						}
					} else {
						common.CheckErr(errors.New("密码认证和密钥认证必须设置其一"))
					}
					sessions[alias] = session
					config.WriteYamlConfig(sessions)
					return nil
				},
			},
			{
				Name:  "rm",
				Usage: "remove a ssh session to the list",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "all",
						Usage:    "delete all",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("all") {
						config.DelYamlFile()
						return nil
					}
					if arg := c.Args().First(); arg != "" {
						sessions, success := config.ReadYamlConfig()
						if !success {
							fmt.Printf("list is empty, please execute command `%s add` first\n", common.ProjectName)
							os.Exit(0)
						}
						_, ok := sessions[arg]
						if !ok {
							fmt.Printf("%s not found\n", arg)
							os.Exit(0)
						}
						delete(sessions, arg)
						config.WriteYamlConfig(sessions)
						fmt.Println("success")
						return nil
					}
					fmt.Println("alias argument not set")
					return nil
				},
			},
			{
				Name: "ls",
				Usage: "show session list",
				Action: func(c *cli.Context) error {
					sessions, success := config.ReadYamlConfig()
					if !success {
						fmt.Printf("list is empty, please execute command `%s add` first\n", common.ProjectName)
						os.Exit(0)
					}

					var password string
					for k, v := range sessions {
						if v.Password != "" {
							password = "password: " + v.Password
						} else if v.PrivateKeyPath != "" {
							password = "private key path: " + v.PrivateKeyPath
						} else {
							password = "auth none"
						}
						fmt.Printf("%s: %s@%s:%d(%s)\n", k, v.User, v.Host, v.Port, password)
					}
					return nil
				},
			},
			{
				Name:  "edit",
				Usage: "modify a ssh session to the list",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dst",
						Aliases:  []string{"d"},
						Usage:    "`alias` to be modified",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "host",
						Aliases:  []string{"i"},
						Usage:    "ip address or host",
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
						Name:    "password",
						Aliases: []string{"pass"},
						Usage:   "password auth",
					},
					&cli.StringFlag{
						Name:    "private-key",
						Aliases: []string{"key"},
						Usage:   "sshkey auth: Non-absolute path, will join home path + private key `path`",
					},
					&cli.StringFlag{
						Name:    "key-passphrase",
						Aliases: []string{"key-pass"},
						Usage:   "private key password",
					},
				},
				Action: func(c *cli.Context) error {
					dst := c.String("dst")
					sessions, _ := config.ReadYamlConfig()
					_, ok := sessions[dst]
					if !ok {
						fmt.Printf("%s is not exist in the list\n", dst)
						os.Exit(0)
					}

					session := sessions[dst]
					if c.String("host") != "" {
						session.Host = c.String("host")
					}
					if c.String("username") != "" {
						session.User = c.String("username")
					}
					if c.Int("port") != 0 {
						session.Port = c.Int("port")
					}
					if c.String("password") != "" {
						session.Password = c.String("password")
					}
					if c.String("private-key") != "" {
						session.PrivateKeyPath = common.PrivateKeyPath(c.String("private-key"))
					}
					if c.String("key-passphrase") != "" {
						session.KeyPassphrase = c.String("key-passphrase")
					}
					sessions[dst] = session
					config.WriteYamlConfig(sessions)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	common.CheckErr(err)
}
