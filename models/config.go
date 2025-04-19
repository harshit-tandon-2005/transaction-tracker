package models

type (
	ThirdPartyApiConfig struct {
		BaseURL string `yaml:"BASE_URL"`
		ApiKey  string `yaml:"API_KEY"`
		Retries int    `yaml:"RETRIES"`
	}
	Config struct {
		Etherscan     ThirdPartyApiConfig `yaml:"ETHERSCAN"`
		Blockscout    ThirdPartyApiConfig `yaml:"BLOCKSCOUT"`
		WalletAddress string              `yaml:"WALLET_ADDRESS"`
	}
)
