package main

import (
	"github.com/reed/cmd/node/command"
	"github.com/reed/log"
)

func main() {
	log.Init()
	command.Execute()
}
