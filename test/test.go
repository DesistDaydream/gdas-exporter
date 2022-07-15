package main

import (
	"fmt"
	"time"

	"github.com/DesistDaydream/gdas-exporter/pkg/gdassdk"
)

func main() {
	password := "HLgd@S123"
	token, err := gdassdk.GetToken("https://172.38.30.192:8003", "system", password)
	if err != nil {
		panic(err)
	}

	client := gdassdk.NewServices("https://172.38.30.192:8003", token, time.Second*10)

	data, err := client.Das.GetDas()
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
}
