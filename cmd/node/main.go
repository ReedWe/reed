package main

import (
	"github.com/tybc/cmd/node/command"
	"github.com/tybc/log"
)

func main() {
	log.Init()
	command.Execute()
}
