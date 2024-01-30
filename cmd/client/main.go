package main

import (
	"github.com/zelas91/goph-keeper/internal/client/commands"
	"github.com/zelas91/goph-keeper/internal/client/request"
)

func main() {
	commands.New(nil, commands.WithUserCommand(request.NewClient("localhost:8080"))).Start()
}
