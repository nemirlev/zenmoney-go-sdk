package zenmoney

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("token is not provided")
	}

	return &Client{
		BaseURL:    "https://api.zenmoney.ru/v8/",
		Token:      token,
		HTTPClient: &http.Client{},
	}, nil
}

func (c *Client) sendRequest(endpoint string, method string, body interface{}) (string, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body:", err)
		}
	}(res.Body)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(resBody), nil
}
