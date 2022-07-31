package server_conf

import "fmt"

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func LoadServerConf(conffile string) ServerConfig {
	fmt.Println("test")
	return ServerConfig{
		Host: "0.0.0.0",
		Port: "8080",
	}
}
