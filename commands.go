package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type CommandRegister struct {
	Name string
	F    func(*state, command) error
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("no command found")
	}
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.registeredCommands[name] = f
	return nil
}

func getCommands() []CommandRegister {
	commandRegisters := []CommandRegister{
		{
			Name: "login",
			F:    handlerLogin,
		},
		{
			Name: "register",
			F:    handlerRegister,
		},
		{
			Name: "reset",
			F:    handlerReset,
		},
		{
			Name: "users",
			F:    handlerList,
		},
	}
	return commandRegisters
}
