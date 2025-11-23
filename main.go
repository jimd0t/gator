package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/jimd0t/gator/internal/config"
	"github.com/jimd0t/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	config  *config.Config
	queries *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		panic(fmt.Errorf("error trying to open db connection: %e", err))
	}

	dbQueries := database.New(db)

	s := state{config: &cfg, queries: dbQueries}
	cmds := commands{registeredCommands: make(map[string]func(*state, command) error)}
	commandRegisters := getCommands()
	for _, register := range commandRegisters {
		err = cmds.register(register.Name, register.F)
		if err != nil {
			panic(err)
		}
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
