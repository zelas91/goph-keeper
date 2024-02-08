package main

import (
	"fmt"
	"github.com/zelas91/goph-keeper/internal/client"
)

var (
	buildCommit = "N/A"
	buildDate   = "N/A"
)

func main() {
	fmt.Printf("client build data (%s) version (%s)\n", buildDate, buildCommit)
	client.NewClient("localhost:8080").Start()
}
