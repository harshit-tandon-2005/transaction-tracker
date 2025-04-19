package usecase

import (
	"fmt"

	"github.com/coin-tracker/transaction-tracker/models"
	thirdparty "github.com/coin-tracker/transaction-tracker/third-party"
)

func GenerateTransactionReports(providerType string, config models.Config) error {

	dataProvider, err := thirdparty.NewDataProvider(providerType, config)
	if err != nil {
		fmt.Printf("Error creating data provider: %v\n", err)
		return err
	}

	actionTagMap := map[string]string{
		"EXTERNAL_REPORT": "txlist",
		"INTERNAL_REPORT": "txlistinternal",
		"ERC20_REPORT":    "tokentx",
		"ERC721_REPORT":   "tokennfttx",
	}
	for key, value := range actionTagMap {
		err = GenerateReports(dataProvider, config.WalletAddress, value, key)
		if err != nil {
			fmt.Printf("[%s] Error generating transaction report: %v\n", key, err)
			return err
		}
	}

	// err = ExternalReport(dataProvider, config.WalletAddress)
	// if err != nil {
	// 	fmt.Printf("Error generating external transacrion report: %v\n", err)
	// 	return err
	// }
	// err = InternalReport(dataProvider, config.WalletAddress)
	// if err != nil {
	// 	fmt.Printf("Error generating internal transacrion report: %v\n", err)
	// 	return err
	// }

	// err = Erc20Report(dataProvider, config.WalletAddress)
	// if err != nil {
	// 	fmt.Printf("Error generating erc20 transacrion report: %v\n", err)
	// 	return err
	// }
	// err = Erc721Report(dataProvider, config.WalletAddress)
	// if err != nil {
	// 	fmt.Printf("Error generating erc721 transacrion report: %v\n", err)
	// 	return err
	// }

	// data, err := dataProvider.FetchTransactionData(walletAddress)
	// if err != nil {
	// 	fmt.Printf("Error fetching transaction data using provider %s: %v\n", constants.PROVIDER_ETHERSCAN, err)
	// 	os.Exit(1)
	// }

	return nil
}

func GenerateReports(dataProvider thirdparty.BlockchainDataProvider, walletAddress, action, tag string) error {

	url := dataProvider.BuildRequestURL(action, walletAddress) // Use the helper

	fmt.Printf("Request URL for [%s]- %s", tag, url)

	res, err := dataProvider.FetchTransactionData(url, tag)
	if err != nil {
		fmt.Printf("Error fetching transaction data: %v\n", err)
		return err
	}

	fmt.Printf("Transaction data: %s\n", res)

	return nil
}

/*
Result from txlist -> External Transaction
Result from txlistinternal -> Internal Transaction
Result from tokentx -> ERC-20 Token Transfer
Result from tokennfttx -> ERC-721/ERC-1155 (NFT) Token Transfer
*/

// func ExternalReport(dataProvider thirdparty.BlockchainDataProvider, walletAddress string) error {

// 	tag := "EXTERNAL_REPORT"

// 	url := dataProvider.BuildRequestURL("txlist", walletAddress) // Use the helper

// 	fmt.Printf("Request URL for [%s]- %s", tag, url)

// 	res, err := dataProvider.FetchTransactionData(url, tag)
// 	if err != nil {
// 		fmt.Printf("Error fetching transaction data: %v\n", err)
// 		return err
// 	}

// 	fmt.Printf("Transaction data: %s\n", res)

// 	return nil

// }

// func InternalReport(dataProvider thirdparty.BlockchainDataProvider, walletAddress string) error {
// 	tag := "INTERNAL_REPORT"
// 	url := dataProvider.BuildRequestURL("txlistinternal", walletAddress) // Use the helper

// 	fmt.Printf("Request URL for [%s]- %s", tag, url)

// 	res, err := dataProvider.FetchTransactionData(url, tag)
// 	if err != nil {
// 		fmt.Printf("Error fetching transaction data: %v\n", err)
// 		return err
// 	}

// 	fmt.Printf("Transaction data: %s\n", res)

// 	return nil
// }

// func Erc20Report(dataProvider thirdparty.BlockchainDataProvider, walletAddress string) error {
// 	tag := "ERC20_REPORT"
// 	url := dataProvider.BuildRequestURL("tokentx", walletAddress) // Use the helper

// 	fmt.Printf("Request URL for [%s]- %s", tag, url)

// 	res, err := dataProvider.FetchTransactionData(url, tag)
// 	if err != nil {
// 		fmt.Printf("Error fetching transaction data: %v\n", err)
// 		return err
// 	}

// 	fmt.Printf("Transaction data: %s\n", res)

// 	return nil
// }

// func Erc721Report(dataProvider thirdparty.BlockchainDataProvider, walletAddress string) error {
// 	tag := "ERC721_REPORT"
// 	url := dataProvider.BuildRequestURL("tokennfttx", walletAddress) // Use the helper

// 	fmt.Printf("Request URL for External Transactions- %s", url)

// 	res, err := dataProvider.FetchTransactionData(url, tag)
// 	if err != nil {
// 		fmt.Printf("Error fetching transaction data: %v\n", err)
// 		return err
// 	}

// 	fmt.Printf("Transaction data: %s\n", res)

// 	return nil
// }
