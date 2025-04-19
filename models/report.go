package models

type ReportResponse struct {
	TransactionHash      string `json:"transactionHash" csv:"Transaction Hash"`
	DateTime             string `json:"dateTime" csv:"Date Time"`
	FromAddress          string `json:"fromAddress" csv:"From Address"`
	ToAddress            string `json:"toAddress" csv:"To Address"`
	TransactionType      string `json:"transactionType" csv:"Transaction Type"`
	AssetContractAddress string `json:"assetContractAddress" csv:"Asset Contract Address"`
	AssetSymbolName      string `json:"assetSymbolName" csv:"Asset Symbol Name"`
	TokenID              string `json:"tokenID" csv:"Token ID"`
	ValueAmount          string `json:"valueAmount" csv:"Value Amount"`
	GasFeeEth            string `json:"gasFeeEth" csv:"Gas Fee (ETH)"`
}
