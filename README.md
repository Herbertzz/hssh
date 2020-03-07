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
$ hssh add -a alias -i host [-u nameuser] [-p port] [-pass password] [--key private-key] [--key-pass key-passphrase]
# -a alias: 保存的别名，必要
# -i host: 服务器的 IP 地址，必要
# [-u nameuser]: 用户名(默认值: root), 可选
# [-p port]: 端口(默认值: 22), 可选
# [-pass password]: 密码认证, 可选, 密码认证和密钥认证必须设置其一
# [--key private-key]: 密钥认证, 可选, 密码认证和密钥认证必须设置其一
# [--key-pass key-passphrase]: 密钥密码, 可选

# 查看保存的服务器列表
$ hssh ls

# 删除指定的 alias 的服务器信息
$ hssh rm [alias]
# 删除配置文件
# hssh rm --all

# 修改指定 alias 的服务器信息
$ hssh edit -d alias [-i host] [-u nameuser] [-p port] [-pass password] [--key private-key] [--key-pass key-passphrase] [-h]
# -d alias: 要修改的别名
# 其他同 hssh add 的参数
```

## 配置文件位置及格式
配置文件路径: `~/.hssh.yaml`, 存储格式:
```yaml
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