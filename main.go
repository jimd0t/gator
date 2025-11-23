package main

import (
	"fmt"

	"github.com/jimd0t/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.DbURL)
	fmt.Println(cfg.CurrentUserName)

	err = cfg.SetUser("jim")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.DbURL)
	fmt.Println(cfg.CurrentUserName)
}
