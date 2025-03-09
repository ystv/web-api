package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ystv/web-api/client/types"
	"github.com/ystv/web-api/services/people"
	"github.com/ystv/web-api/services/stream"
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

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

func (c *Client) GetUsersPagination(ctx context.Context, apiKey string, options types.UsersListPaginationOptions) (people.UserFullPagination, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return people.UserFullPagination{}, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = "/v1/internal/people/users/pagination"
	q := u.Query()
	if options.Size != nil {
		q.Set("size", fmt.Sprintf("%d", options.Size))
	}
	if options.Page != nil {
		q.Set("page", fmt.Sprintf("%d", options.Page))
	}
	if options.Search != nil {
		q.Set("search", *options.Search)
	}
	if options.Column != nil {
		q.Set("column", fmt.Sprintf("%s", *options.Column))
	}
	if options.Enabled != nil {
		q.Set("enabled", fmt.Sprintf("%s", *options.Enabled))
	}
	if options.Deleted != nil {
		q.Set("deleted", fmt.Sprintf("%s", *options.Deleted))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return people.UserFullPagination{}, err
	}

	var res people.UserFullPagination
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.UserFullPagination{}, err
	}

	return res, nil
}

func (c *Client) FindStreamEndpoint(ctx context.Context, apiKey string, options types.FindStreamEndpointOptions) (stream.Endpoint, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/people/users/pagination")
	if err != nil {
		return stream.Endpoint{}, fmt.Errorf("invalid base URL: %w", err)
	}

	optionBytes, err := json.Marshal(options)
	if err != nil {
		return stream.Endpoint{}, err
	}

	buf := bytes.NewBuffer(optionBytes)

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), buf)
	if err != nil {
		return stream.Endpoint{}, err
	}

	var res stream.Endpoint
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return stream.Endpoint{}, err
	}

	return res, nil
}
