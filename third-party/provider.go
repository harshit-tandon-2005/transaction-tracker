package thirdparty

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coin-tracker/transaction-tracker/models"
	"github.com/coin-tracker/transaction-tracker/shared/constants"
)

// BlockchainDataProvider defines the interface for fetching data from blockchain explorers.
type BlockchainDataProvider interface {
	FetchTransactionData(url, tag string) (string, error)

	// builds a request URL for the provider and if in future if a provider has a requets body another method can be defined to build requets body
	BuildRequestURL(action, walletAddress string) string
}

func NewDataProvider(providerType string, config models.Config) (BlockchainDataProvider, error) {
	// Use a shared HTTP client with a reasonable timeout
	httpClient := &http.Client{Timeout: 15 * time.Second}
	var err error

	switch strings.ToLower(providerType) {
	case constants.PROVIDER_ETHERSCAN:
		// Use default URL if not provided in config
		return NewEtherscanProvider(config.Etherscan, httpClient), nil

	// TODO: Add support Blockscout provider
	case constants.PROVIDER_BLOCKSCOUT:
		// Blockscout API key usage is optional/depends on instance
		// provider := NewBlockscoutProvider(config.Blockscout.Key, config.Blockscout.BaseURL, httpClient)
		// return provider, nil

	default:
		// Return a wrapped error for better context
		err = fmt.Errorf("unknown provider: %s", providerType)
	}
	return nil, err
}
