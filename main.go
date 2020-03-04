package main

func main() {
	config := SSHConfig{
		User:          "herbertzz",
		// Password:      "henyi",
		PrivateKeyPath:	"/Users/herbertzz/.ssh/id_rsa",
		Host:          "149.129.103.98",
		// Port:          22,
	}
	openSSH(config)
}
