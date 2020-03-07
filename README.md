# SSH 连接管理工具
# 简介
个人需要管理的服务器比较多，自带的 ssh 只能设置别名，没法快速查询我设置了哪些别名，每次都需要打开配置文件进行查看，很是麻烦。
故开发这个工具，提升下效率
## 使用
```bash
# 查看帮助
$ hssh h

# ssh 连接指定 alias 的服务器
$ hssh [alias]

# 添加一个服务器信息
$ hssh add -i host [-u nameuser] [-p port] [--auth method] [-pass password] [--key private-key] [--key-pass key-passphrase] alias
# -i host: 服务器的 IP 地址，必要
# [-u nameuser]: 用户名(默认值: root), 可选
# [-p port]: 端口(默认值: 22), 可选
# [--auth method]: 认证方式(默认值: password), 可选, 只支持 password 和 key
# [-pass password]: 密码认证, --auth 为 password 且该字段为空时，会在登录时询问用户输入密码
# [--key private-key]: 密钥认证(默认值: default), 可选, 对应配置文件 keys 字段
# [--key-pass key-passphrase]: 密钥密码, 可选
# alias: 保存的别名，必要

# 查看保存的服务器列表
$ hssh ls

# 删除指定的 alias 的服务器信息
$ hssh rm [alias]
# 删除配置文件
# hssh rm --all

# 修改指定 alias 的服务器信息
$ hssh edit [-i host] [-u nameuser] [-p port] [--auth method] [-pass password] [--key private-key] [--key-pass key-passphrase] alias
# 同 hssh add 的参数

# 卸载
$ hssh uninstall [--all]
# --all: 连同配置文件一起删除，不带此参数则保留配置文件


# ---- keys 管理 ----
# 查看 keys 列表
# hssh keys
# 查看指定 key 的值
# hssh keys key
# 添加一个 key 到 kyes
# hssh keys add key path
# 编辑 key 的值
# hssh keys edit key path
# 删除指定的 key
# hssh keys rm key
```

## 配置文件位置及格式
配置文件路径: `~/.hssh.yaml`, 存储格式:
```yaml
keys:
  default: /Users/herbertzz/.ssh/id_rsa  # 程序自动生成, 默认为 {当前用户主目录路径}/.ssh/id_rsa
servers:
  alias:
    username: root
    host: 192.168.1.3
    port: 22
    password: ""
    private_key: /Users/herbertzz/.ssh/id_rsa
    key_passphrase: ""

```

## 待开发
* `hssh ls` 支持模糊匹配功能
* `hssh rm {index}` 支持使用序号删除