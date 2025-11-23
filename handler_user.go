package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	username := cmd.Args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("User %s switched successfully!", username)
	return nil
}
