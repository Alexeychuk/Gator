package main

import (
	"fmt"

	"github.com/Alexeychuk/Gator/internal/config"
)

func main() {
	currentConfig, err := config.Read()
	if err != nil {
		return
	}

	fmt.Printf("db: %s, user: %s\n", currentConfig.DBUrl, currentConfig.Username)

	err = currentConfig.SetUser("Vova")
	if err != nil {
		return
	}

	currentConfig, err = config.Read()
	if err != nil {
		return
	}

	fmt.Printf("db: %s, user: %s\n", currentConfig.DBUrl, currentConfig.Username)

}
