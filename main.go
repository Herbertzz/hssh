package main

import (
	"fmt"
	"github.com/urfave/cli"
	"hssh/common"
	"hssh/config"
	"hssh/ssh"
	"os"
)

func main() {
	app()
}

func app()  {
	app := &cli.App{
		Name:    config.ProjectName,
		Usage:   "manage ssh sessions",
		Version: config.Version,
		Action: func(c *cli.Context) error {
			if c.Args().First() != "" {
				configs, success := config.ReadYamlConfig()
				sessions := configs.Servers
				if !success {
					fmt.Printf("please execute command `%s h` for help\n", config.ProjectName)
					os.Exit(0)
				}
				session, ok := sessions[c.Args().First()]
				if !ok {
					fmt.Printf("do not find session named: %s\n", c.Args().First())
					os.Exit(0)
				}

				key := ""
				if session.AuthMethod == "key" {
					key, ok = configs.Keys[session.PrivateKey]
					if !ok {
						fmt.Printf("%s not exist in keys\n", session.PrivateKey)
						os.Exit(0)
					}
				}

				ssh.OpenSSH(session, key)
				return nil
			}
			fmt.Printf("please execute command `%s h` for help\n", config.ProjectName)
			return nil
		},
		Commands: []*cli.Command{
			// 添加
			{
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
						Name:    "auth",
						Usage:   "auth `method`: password or key",
						DefaultText: "password",
						Value: "password",
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
				Action: func(c *cli.Context) error {
					if arg := c.Args().First(); arg != "" {
						configs, _ := config.ReadYamlConfig()
						sessions := configs.Servers
						if len(sessions) == 0 {
							sessions = make(map[string]config.Server)
						}
						_, ok := sessions[arg]
						if ok {
							fmt.Printf("%s is already in the list\n", arg)
							os.Exit(0)
						}

						session := config.Server{}
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
						// 认证方式
						authMethod := c.String("auth")
						if authMethod == "password" {
							if c.String("password") == "" {
								fmt.Println("Error: auth method is password, cannot be empty, Try adding '--pass'")
								return nil
							}
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
						config.WriteYamlConfig(sessions)
						return nil
					}
					fmt.Println("alias is not set")
					return nil
				},
			},
			// 删除
			{
				Name:  "rm",
				Usage: "remove a ssh session to the list",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "all",
						Usage:    "delete all session",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("all") {
						sessions := make(map[string]config.Server, 0)
						config.WriteYamlConfig(sessions)
						return nil
					}
					if arg := c.Args().First(); arg != "" {
						configs, success := config.ReadYamlConfig()
						sessions := configs.Servers
						if !success {
							fmt.Printf("list is empty, please execute command `%s add` first\n", config.ProjectName)
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
			// 查看
			{
				Name: "ls",
				Usage: "show session list",
				Action: func(c *cli.Context) error {
					configs, success := config.ReadYamlConfig()
					sessions := configs.Servers
					if !success {
						fmt.Printf("list is empty, please execute command `%s add` first\n", config.ProjectName)
						os.Exit(0)
					}

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
				},
			},
			// 编辑
			{
				Name:  "edit",
				Usage: "modify a ssh session to the list",
				Flags: []cli.Flag{
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
						Name:    "auth",
						Usage:   "auth `method`: password or key",
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
				Action: func(c *cli.Context) error {
					if arg := c.Args().First(); arg != "" {
						configs, _ := config.ReadYamlConfig()
						sessions := configs.Servers
						_, ok := sessions[arg]
						if !ok {
							fmt.Printf("%s is not exist in the list\n", arg)
							os.Exit(0)
						}

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
							if c.String("password") == "" {
								fmt.Println("Error: auth method is password, cannot be empty, Try adding '--pass'")
								return nil
							}
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

						sessions[arg] = session
						config.WriteYamlConfig(sessions)
						return nil
					}
					fmt.Println("alias is not set")
					return nil
				},
			},
			// 卸载
			{
				Name: "uninstall",
				Usage: "unistall the app",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "all",
						Usage:    "Delete with the configuration file",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("all") {
						config.DelYamlFile()
					}
					common.DelCurrentApp()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	common.CheckErr(err)
}
