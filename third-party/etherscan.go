package thirdparty

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/coin-tracker/transaction-tracker/models"
	"github.com/coin-tracker/transaction-tracker/shared/util"
)

type EtherscanProvider struct {
	ApiKey     string
	BaseURL    string
	Client     *http.Client
	ListParams map[string]string
}

// NewEtherscanProvider creates a new Etherscan provider instance.
// Typically called by the factory, but can be used directly.
func NewEtherscanProvider(config models.ThirdPartyApiConfig, client *http.Client) *EtherscanProvider {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.etherscan.io/api"
	}

	// Default list params, some values if required can be taken from a config file or some input
	listParams := map[string]string{
		"startblock": "0",
		"endblock":   "99999999",
		"sort":       "asc",
		"module":     "account",
	}

	return &EtherscanProvider{
		ApiKey:     config.ApiKey,
		BaseURL:    config.BaseURL,
		Client:     client,
		ListParams: listParams,
	}
}

func (p *EtherscanProvider) BuildRequestURL(action, walletAddress string) string {
	// Base URL should already be set in the provider, default to official if needed
	baseURL := p.BaseURL

	// Prepare query parameters
	queryParams := url.Values{}
	queryParams.Set("address", walletAddress)
	queryParams.Set("action", action)
	queryParams.Set("apikey", p.ApiKey) // Add API key automatically

	// Add other specific parameters from the map
	if p.ListParams != nil {
		for key, value := range p.ListParams {
			queryParams.Set(key, value)
		}
	}

	// Construct the full URL
	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
	return fullURL
}

// FetchTransactionData implements the BlockchainDataProvider interface for Etherscan.
func (p *EtherscanProvider) FetchTransactionData(url, tag string) (string, error) {

	res, err := util.TriggerHttpRequest(http.MethodGet, url, tag, p.Client)
	if err != nil {
		return "", err
	}
	return res, nil
}
