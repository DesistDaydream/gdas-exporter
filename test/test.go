package main

import (
	"fmt"

	"github.com/DesistDaydream/gdas-exporter/pkg/gdassdk"
)

func main() {

	// sClient := services.NewLoginService(client)

	// login, err := sClient.PostLogin(&core.PostLogin{
	// 	Username: "system",
	// 	Password: "HLgd@S123",
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// client.Token = login.Token

	token, err := gdassdk.GetToken("https://172.38.30.192:8003", "system", "HLgd@S123")
	if err != nil {
		panic(err)
	}
	fmt.Println(token)

	s := gdassdk.NewService("https://172.38.30.192:8003", token)

	fmt.Println(s.Client.Token)

	nodes, err := s.Node.GetNode()
	if err != nil {
		panic(err)
	}

	fmt.Println(nodes)
}
