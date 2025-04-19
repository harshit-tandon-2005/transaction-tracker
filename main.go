package main

import (
	"fmt"
	"os"

	"github.com/coin-tracker/transaction-tracker/models"
	"github.com/coin-tracker/transaction-tracker/shared/constants"
	usecase "github.com/coin-tracker/transaction-tracker/usecase"
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

	err = usecase.GenerateTransactionReports(constants.PROVIDER_ETHERSCAN, config)
	if err != nil {
		fmt.Printf("Error generating transaction reports: %v\n", err)
		os.Exit(1)
	}

	// --- Data Fetching ---
	// Use the provider interface to fetch data

	// --- Output ---
	fmt.Println("\n--- Fetched Data ---")
	fmt.Println("--------------------")
	fmt.Println("Operation completed.")
}
