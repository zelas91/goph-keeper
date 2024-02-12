package main

import (
	"flag"
	"fmt"

	"github.com/zelas91/goph-keeper/internal/client"
)

var (
	buildCommit = "N/A"
	buildDate   = "N/A"
)

var addr *string

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
}
func main() {
	fmt.Printf("client build data (%s) version (%s)\n", buildDate, buildCommit)
	flag.Parse()
	client.NewClient(*addr).Start()
}
