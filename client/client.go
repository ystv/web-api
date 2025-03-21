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
		return fmt.Errorf("request failed: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode > http.StatusNoContent && res.StatusCode <= http.StatusBadGateway {
		errRes := types.ErrorResponse{
			Code: res.StatusCode,
		}
		if err = json.NewDecoder(res.Body).Decode(&errRes.Message); err == nil {
			return fmt.Errorf("invalid status code: %d, body: %s", errRes.Code, errRes.Message)
		}

		switch res.StatusCode {
		case http.StatusUnauthorized:
			return errors.New("unauthorised")
		case http.StatusNotFound:
			return errors.New("the requested url is not currently found: " + req.URL.Path)
		case http.StatusInternalServerError:
			return errors.New("internal server error: " + errRes.Message)
		case http.StatusNotImplemented:
			return errors.New("the requested url is not currently implemented: " + req.URL.Path)
		case http.StatusBadGateway:
			return errors.New("bad gateway, web api may be down: " + req.URL.Path)
		}

		return fmt.Errorf("unknown error, status code: %d, message: %s", res.StatusCode, errRes.Message)
	}

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	return nil
}
