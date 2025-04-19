package models

import (
	"encoding/json"
)

type (
	EtherscanBaseResponse struct {
		Status  string          `json:"status"` // "1" for success, "0" for error
		Message string          `json:"message"`
		Result  json.RawMessage `json:"result"` // Use RawMessage to defer parsing
	}

	// Structure for a single normal/external transaction result
	ExternalTransaction struct {
		BlockNumber       string `json:"blockNumber"`
		TimeStamp         string `json:"timeStamp"`
		Hash              string `json:"hash"`
		Nonce             string `json:"nonce"`
		BlockHash         string `json:"blockHash"`
		TransactionIndex  string `json:"transactionIndex"`
		From              string `json:"from"`
		To                string `json:"to"`
		Value             string `json:"value"` // Value in Wei
		Gas               string `json:"gas"`
		GasPrice          string `json:"gasPrice"`
		IsError           string `json:"isError"`          // "0" for success, "1" for error
		TxReceiptStatus   string `json:"txreceipt_status"` // "1" for success, "0" for failure
		Input             string `json:"input"`
		ContractAddress   string `json:"contractAddress"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		GasUsed           string `json:"gasUsed"`
		Confirmations     string `json:"confirmations"`
		MethodId          string `json:"methodId"`
		FunctionName      string `json:"functionName"`
	}

	// Structure for a single internal transaction result
	InternalTransaction struct {
		BlockNumber     string `json:"blockNumber"`
		TimeStamp       string `json:"timeStamp"`
		Hash            string `json:"hash"` // Hash of the parent transaction
		From            string `json:"from"`
		To              string `json:"to"`
		Value           string `json:"value"` // Value in Wei
		ContractAddress string `json:"contractAddress"`
		Input           string `json:"input"`
		Type            string `json:"type"` // e.g., "call", "create"
		Gas             string `json:"gas"`
		GasUsed         string `json:"gasUsed"`
		TraceId         string `json:"traceId"`
		IsError         string `json:"isError"` // "0" for success, "1" for error
		ErrCode         string `json:"errCode"`
	}

	// Structure for a single ERC-20 token transfer result
	TokenTransaction struct {
		BlockNumber       string `json:"blockNumber"`
		TimeStamp         string `json:"timeStamp"`
		Hash              string `json:"hash"`
		Nonce             string `json:"nonce"`
		BlockHash         string `json:"blockHash"`
		From              string `json:"from"`
		ContractAddress   string `json:"contractAddress"`
		To                string `json:"to"`
		Value             string `json:"value"` // Token amount (consider decimals)
		TokenName         string `json:"tokenName"`
		TokenSymbol       string `json:"tokenSymbol"`
		TokenDecimal      string `json:"tokenDecimal"`
		TransactionIndex  string `json:"transactionIndex"`
		Gas               string `json:"gas"`
		GasPrice          string `json:"gasPrice"`
		GasUsed           string `json:"gasUsed"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		Input             string `json:"input"` // Typically "deprecated" for token transfers
		Confirmations     string `json:"confirmations"`
	}

	// Structure for a single ERC-721 / ERC-1155 NFT transfer result
	NftTransaction struct {
		BlockNumber       string `json:"blockNumber"`
		TimeStamp         string `json:"timeStamp"`
		Hash              string `json:"hash"`
		Nonce             string `json:"nonce"`
		BlockHash         string `json:"blockHash"`
		From              string `json:"from"`
		ContractAddress   string `json:"contractAddress"`
		To                string `json:"to"`
		TokenID           string `json:"tokenID"` // The specific NFT ID
		TokenName         string `json:"tokenName"`
		TokenSymbol       string `json:"tokenSymbol"`
		TokenDecimal      string `json:"tokenDecimal"` // Usually "0" for NFTs
		TransactionIndex  string `json:"transactionIndex"`
		Gas               string `json:"gas"`
		GasPrice          string `json:"gasPrice"`
		GasUsed           string `json:"gasUsed"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		Input             string `json:"input"` // Typically "deprecated"
		Confirmations     string `json:"confirmations"`
		// Etherscan might sometimes include TokenType for ERC-1155, add if needed
		// TokenType         string `json:"tokenType"` // e.g., "ERC-721", "ERC-1155"
	}
)
