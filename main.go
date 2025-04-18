package main

import (
	"fmt"
	"os"

	"github.com/coin-tracker/transaction-tracker/models"

	"gopkg.in/yaml.v3"
)

// Config struct to hold the configuration from config.yml

func main() {
	// Read the config file
	configFile, err := os.ReadFile("config.yml")
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	// Define a variable of type Config
	config := models.Config{}

	// Unmarshal the YAML data into the config struct
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Printf("Error unmarshalling config data: %v\n", err)
		os.Exit(1)
	}

	// Print a value from the config to verify
	fmt.Printf("Successfully loaded config:\n %+v", config)

}
