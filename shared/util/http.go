package util

import (
	"fmt"
	"io"
	"net/http"
)

func TriggerHttpRequest(requestMethod, requestUrl, tag string, client *http.Client) (string, error) {

	req, err := http.NewRequest(requestMethod, requestUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create etherscan request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute etherscan request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Read body for context even on error
		return "", fmt.Errorf("etherscan API request failed with status %s: %s", resp.Status, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read etherscan response body: %w", err)
	}

	return string(bodyBytes), nil
}
