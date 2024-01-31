package main

import (
	"github.com/zelas91/goph-keeper/internal/client"
)

func main() {
	client.NewClient("localhost:8080").Start()
}
