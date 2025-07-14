package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DBUrl    string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.Username = user

	err := write(c)

	if err != nil {
		fmt.Print("Error in set user")
		return err
	}

	return nil
}

const configFileName = ".gatorconfig.json"

func write(cfg *Config) error {
	configLocation, err := getConfigFilePath()
	if err != nil {
		fmt.Print("Error extracting config location for write")
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		fmt.Print("Error marshalling config")
		return err
	}

	if err := os.WriteFile(configLocation, data, os.FileMode(0644)); err != nil {
		fmt.Print("Error wtiting config to file")
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	configLocation, err := os.UserHomeDir()
	if err != nil {
		fmt.Print("No Config location found")
		return "", err
	}

	return filepath.Join(configLocation, configFileName), nil
}

func Read() (*Config, error) {

	configLocation, err := getConfigFilePath()
	if err != nil {
		fmt.Print("Error extracting config")
		return nil, err
	}

	bytesData, err := os.ReadFile(configLocation)
	if err != nil {
		fmt.Print("Error bytes config")
		return nil, err
	}

	var config Config

	if err := json.Unmarshal(bytesData, &config); err != nil {
		fmt.Print("Error parsing config json")
		return nil, err
	}

	return &config, nil
}
