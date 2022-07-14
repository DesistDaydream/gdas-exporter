package main

import (
	"fmt"

	core "github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/core/v1"
	"github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/services"
)

func main() {
	client := core.NewClient("https://172.38.30.192:8003")
	sClient := services.NewLoginService(client)

	login, err := sClient.PostLogin(&core.PostLogin{
		Username: "system",
		Password: "HLgd@S123",
	})
	if err != nil {
		panic(err)
	}

	client.Token = login.Token

	nClient := services.NewNodeService(client)
	fmt.Println(client.Token)
	nodes, err := nClient.GetNode()
	if err != nil {
		panic(err)
	}

	fmt.Println(nodes)
}
