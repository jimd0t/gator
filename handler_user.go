package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jimd0t/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	username := cmd.Args[0]
	_, err := s.queries.GetUser(context.Background(), username)
	if err != nil {
		return err
	}
	err = s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("User %s switched successfully!", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	username := cmd.Args[0]

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	user, err := s.queries.CreateUser(context.Background(), newUser)
	if err != nil {
		return err
	}
	fmt.Printf("Created user %s (%v) successfully!\n", user.Name, user.ID)
	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	return err
}

func handlerReset(s *state, cmd command) error {
	err := s.queries.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Successfully erased all user information from the database")
	return err
}

func handlerList(s *state, cmd command) error {
	users, err := s.queries.GetUsers(context.Background())
	if err != nil {
		return err
	}
	currentUser := s.config.CurrentUserName
	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == currentUser {
			fmt.Println(" (current)")
		} else {
			fmt.Println()
		}
	}
	return nil
}
