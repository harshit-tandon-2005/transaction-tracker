package usecase

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/coin-tracker/transaction-tracker/models"
	"github.com/coin-tracker/transaction-tracker/shared/constants"
	"github.com/coin-tracker/transaction-tracker/shared/util"
	thirdparty "github.com/coin-tracker/transaction-tracker/third-party"
)

func GenerateTransactionReports(providerType string, config models.Config) error {

	dataProvider, err := thirdparty.NewDataProvider(providerType, config)
	if err != nil {
		fmt.Printf("Error creating data provider: %v\n", err)
		return err
	}

	/*
		Result from txlist -> External Transaction
		Result from txlistinternal -> Internal Transaction
		Result from tokentx -> ERC-20 Token Transfer
		Result from tokennfttx -> ERC-721/ERC-1155 (NFT) Token Transfer
	*/
	actionTagMap := map[string]string{
		constants.EXTERNAL_REPORT: constants.EXTERNAL_REPORT_ACTION,
		constants.INTERNAL_REPORT: constants.INTERNAL_REPORT_ACTION,
		constants.ERC20_REPORT:    constants.ERC20_REPORT_ACTION,
		constants.ERC721_REPORT:   constants.ERC721_REPORT_ACTION,
	}
	numTasks := len(actionTagMap)
	// Create a buffered channel to receive potential errors.
	// Buffer size equals the number of tasks to prevent goroutines from blocking on send.
	errChan := make(chan error, numTasks)

	// Use a WaitGroup to wait for all goroutines to finish.
	var wg sync.WaitGroup

	fmt.Printf("Starting concurrent generation of %d reports...\n", numTasks)
	for key, value := range actionTagMap {
		// Increment the WaitGroup counter for each goroutine we are about to launch.
		wg.Add(1)
		go func(k, v string) {
			// Decrement the counter when the goroutine finishes, regardless of success or failure.
			defer wg.Done()

			fmt.Printf("[%s] Starting report generation...\n", k)
			err = GenerateReports(dataProvider, config.WalletAddress, v, k)
			if err != nil {
				fmt.Printf("[%s] Error generating report: %v\n", k, err)
				// Send the error to the error channel. Wrap it for context.
				errChan <- fmt.Errorf("report generation failed for key '%s': %w", k, err)
			}
		}(key, value)
	}

	// Wait for all goroutines launched in the loop to finish.
	fmt.Println("Waiting for report generation tasks to complete...")
	wg.Wait()
	fmt.Println("All report generation tasks finished.")

	// Close the error channel *after* all goroutines are guaranteed to be done.
	// This signals to the reading loop below that no more errors will be sent.
	close(errChan)

	// Check if any errors were sent to the channel.
	// We'll collect the first error encountered. If multiple goroutines fail,
	// only the first error received here will be returned.
	var firstError error
	for err := range errChan {
		if firstError == nil {
			firstError = err
		}
	}

	// Return the first error encountered, or nil if all succeeded.
	if firstError != nil {
		fmt.Printf("One or more report generation tasks failed. Returning first error: %v\n", firstError)
		return firstError
	}

	fmt.Println("All reports generated successfully.")

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

	resp := models.EtherscanBaseResponse{}
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		fmt.Printf("Error unmarshalling transaction data: %v\n", err)
		return err
	}

	switch tag {
	case constants.EXTERNAL_REPORT:
		err = ExternalReport(resp.Result, walletAddress)
	case constants.INTERNAL_REPORT:
		err = InternalReport(resp.Result, walletAddress)
	case constants.ERC20_REPORT:
		err = Erc20Report(resp.Result, walletAddress)
	case constants.ERC721_REPORT:
		err = Erc721Report(resp.Result, walletAddress)
	}

	if err != nil {
		fmt.Printf("Error generating transaction report: %v\n", err)
		return err
	}

	return nil
}

func ExternalReport(res json.RawMessage, walletAddress string) error {

	txList := []models.ExternalTransaction{}
	err := json.Unmarshal(res, &txList)
	if err != nil {
		fmt.Printf("Error unmarshalling transaction data: %v\n", err)
		return err
	}

	if len(txList) == 0 {
		fmt.Printf("No External transactions found for wallet address: %s\n", walletAddress)
		return nil
	}

	csvResp := []models.ReportResponse{}
	for _, tx := range txList {

		dateTime, _ := util.FormatUnixTimestampString(tx.TimeStamp)
		csvResp = append(csvResp, models.ReportResponse{
			TransactionHash:      tx.Hash,
			DateTime:             dateTime,
			FromAddress:          tx.From,
			ToAddress:            tx.To,
			TransactionType:      constants.TRANSACTION_TYPE_ETH_TRANSFER,
			AssetContractAddress: tx.ContractAddress,
			AssetSymbolName:      constants.TOKEN_SYMBOL_ETH,
			TokenID:              "",
			ValueAmount:          tx.Value,
			GasFeeEth:            tx.Gas,
		})
	}

	dir, err := util.GetCurrentWorkingDirectory()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return err
	}

	filePath := filepath.Join(dir, "/files/reports", walletAddress+"_external_report.csv")
	err = util.WriteCSV(filePath, csvResp)
	if err != nil {
		fmt.Printf("Error writing external report to file: %v\n", err)
		return err
	}

	return nil

}

func InternalReport(res json.RawMessage, walletAddress string) error {
	txList := []models.InternalTransaction{}
	err := json.Unmarshal(res, &txList)
	if err != nil {
		fmt.Printf("Error unmarshalling transaction data: %v\n", err)
		return err
	}

	if len(txList) == 0 {
		fmt.Printf("No Internal transactions found for wallet address: %s\n", walletAddress)
		return nil
	}

	csvResp := []models.ReportResponse{}
	for _, tx := range txList {
		dateTime, _ := util.FormatUnixTimestampString(tx.TimeStamp)
		csvResp = append(csvResp, models.ReportResponse{
			TransactionHash:      tx.Hash,
			DateTime:             dateTime,
			FromAddress:          tx.From,
			ToAddress:            tx.To,
			TransactionType:      constants.TRANSACTION_TYPE_INTERNAL_TRANSFER,
			AssetContractAddress: tx.ContractAddress,
			AssetSymbolName:      constants.TOKEN_SYMBOL_ETH,
			TokenID:              "",
			ValueAmount:          tx.Value,
			GasFeeEth:            tx.Gas,
		})
	}

	dir, err := util.GetCurrentWorkingDirectory()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return err
	}

	filePath := filepath.Join(dir, "/files/reports", walletAddress+"_internal_report.csv")
	err = util.WriteCSV(filePath, csvResp)
	if err != nil {
		fmt.Printf("Error writing external report to file: %v\n", err)
		return err
	}
	return nil
}

func Erc20Report(res json.RawMessage, walletAddress string) error {
	txList := []models.TokenTransaction{}
	err := json.Unmarshal(res, &txList)
	if err != nil {
		fmt.Printf("Error unmarshalling transaction data: %v\n", err)
		return err
	}

	if len(txList) == 0 {
		fmt.Printf("No ERC-20 transactions found for wallet address: %s\n", walletAddress)
		return nil
	}

	csvResp := []models.ReportResponse{}
	for _, tx := range txList {
		dateTime, _ := util.FormatUnixTimestampString(tx.TimeStamp)
		csvResp = append(csvResp, models.ReportResponse{
			TransactionHash:      tx.Hash,
			DateTime:             dateTime,
			FromAddress:          tx.From,
			ToAddress:            tx.To,
			TransactionType:      constants.TRANSACTION_TYPE_ERC20_TRANSFER,
			AssetContractAddress: tx.ContractAddress,
			AssetSymbolName:      tx.TokenSymbol + " " + tx.TokenName,
			TokenID:              "",
			ValueAmount:          tx.Value,
			GasFeeEth:            tx.Gas,
		})
	}

	dir, err := util.GetCurrentWorkingDirectory()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return err
	}

	filePath := filepath.Join(dir, "/files/reports", walletAddress+"_erc-20_report.csv")
	err = util.WriteCSV(filePath, csvResp)
	if err != nil {
		fmt.Printf("Error writing external report to file: %v\n", err)
		return err
	}

	return nil
}

func Erc721Report(res json.RawMessage, walletAddress string) error {
	txList := []models.NftTransaction{}
	err := json.Unmarshal(res, &txList)
	if err != nil {
		fmt.Printf("Error unmarshalling transaction data: %v\n", err)
		return err
	}

	if len(txList) == 0 {
		fmt.Printf("No ERC-721 transactions found for wallet address: %s\n", walletAddress)
		return nil
	}

	csvResp := []models.ReportResponse{}
	for _, tx := range txList {
		dateTime, _ := util.FormatUnixTimestampString(tx.TimeStamp)
		csvResp = append(csvResp, models.ReportResponse{
			TransactionHash:      tx.Hash,
			DateTime:             dateTime,
			FromAddress:          tx.From,
			ToAddress:            tx.To,
			TransactionType:      constants.TRANSACTION_TYPE_ERC721_TRANSFER,
			AssetContractAddress: tx.ContractAddress,
			AssetSymbolName:      tx.TokenSymbol + " " + tx.TokenName,
			TokenID:              tx.TokenID,
			ValueAmount:          tx.TransactionIndex,
			GasFeeEth:            tx.Gas,
		})
	}

	dir, err := util.GetCurrentWorkingDirectory()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return err
	}

	filePath := filepath.Join(dir, "/files/reports", walletAddress+"_erc-721_report.csv")
	err = util.WriteCSV(filePath, csvResp)
	if err != nil {
		fmt.Printf("Error writing external report to file: %v\n", err)
		return err
	}

	return nil
}
