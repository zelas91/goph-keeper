package main

import (
	"fmt"
	"github.com/zelas91/goph-keeper/internal/client"
	"github.com/zelas91/goph-keeper/internal/utils/crypto"
)

var (
	secretKey   = ""
	buildCommit = "N/A"
	buildDate   = "N/A"
)

func main() {
	fmt.Printf("client build data (%s) version (%s) ---(%s)---\n", buildDate, buildCommit, secretKey)
	crypt, err := crypto.NewEncrypt(secretKey)
	fmt.Println(crypt, err)
	client.NewClient("localhost:8080").Start()
}
