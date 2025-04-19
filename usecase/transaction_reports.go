package usecase

import (
	"encoding/json"
	"fmt"
	"path/filepath"

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
	for key, value := range actionTagMap {
		err = GenerateReports(dataProvider, config.WalletAddress, value, key)
		if err != nil {
			fmt.Printf("[%s] Error generating transaction report: %v\n", key, err)
			return err
		}
	}

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
