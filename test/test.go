package main

import (
	"fmt"
	"time"

	"github.com/DesistDaydream/gdas-exporter/pkg/gdassdk"
)

func main() {
	password := "XXXXXX"
	token, err := gdassdk.GetToken("https://172.38.30.192:8003", "system", password)
	if err != nil {
		panic(err)
	}
	fmt.Println(token)

	client := gdassdk.NewServices("https://172.38.30.192:8003", token, time.Second*10)

	fmt.Println(client.Client.Token)

	auth, err := client.Auth.GetAuthorize()
	if err != nil {
		panic(err)
	}

	fmt.Println(auth)
}
