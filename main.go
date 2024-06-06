package main

import (
	"github.com/Terracode-Dev/terraui-back/api"
)

func main() {
	server := api.NewServer(":8080")
	server.Run()
}
