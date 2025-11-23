package main

import (
	"errors"
	"os"

	"github.com/jimd0t/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}
	s := state{config: &cfg}
	cmds := commands{registeredCommands: make(map[string]func(*state, command) error)}

	err = cmds.register("login", handlerLogin)
	if err != nil {
		panic(err)
	}

	args := os.Args
	if len(args) < 2 {
		panic(errors.New("expected some arguments but got none"))
	}

	cmdName := args[1]
	cmd := command{
		Name: cmdName,
		Args: args[2:],
	}
	err = cmds.run(&s, cmd)
	if err != nil {
		panic(err)
	}
}
