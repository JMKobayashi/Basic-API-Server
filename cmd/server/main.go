package main

import "github.com/JMKobayashi/Basic-API-Server/configs"

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
