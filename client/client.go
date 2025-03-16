package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ystv/web-api/client/types"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

const ten = 10

func NewClient(baseURL string) (*Client, error) {
	_, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: ten * time.Second,
		},
	}, nil
}

func (c *Client) sendRequest(req *http.Request, apiKey string, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes types.ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	fullResponse := types.SuccessResponse{
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}
