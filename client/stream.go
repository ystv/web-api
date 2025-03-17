package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ystv/web-api/client/types"
	"github.com/ystv/web-api/services/stream"
)

func (c *Client) ListStreamEndpoints(ctx context.Context, apiKey string) ([]stream.Endpoint, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/streams")
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var res []stream.Endpoint
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) FindStreamEndpoint(ctx context.Context, apiKey string, options types.FindStreamEndpointOptions) (stream.Endpoint, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/streams/find")
	if err != nil {
		return stream.Endpoint{}, fmt.Errorf("invalid base URL: %w", err)
	}

	optionBytes, err := json.Marshal(options)
	if err != nil {
		return stream.Endpoint{}, err
	}

	buf := bytes.NewBuffer(optionBytes)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), buf)
	if err != nil {
		return stream.Endpoint{}, err
	}

	var res stream.Endpoint
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return stream.Endpoint{}, err
	}

	return res, nil
}

func (c *Client) AddStreamEndpoint(ctx context.Context, apiKey string, options stream.EndpointAddEditDTO) (stream.Endpoint, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/streams")
	if err != nil {
		return stream.Endpoint{}, fmt.Errorf("invalid base URL: %w", err)
	}

	optionBytes, err := json.Marshal(options)
	if err != nil {
		return stream.Endpoint{}, err
	}

	buf := bytes.NewBuffer(optionBytes)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), buf)
	if err != nil {
		return stream.Endpoint{}, err
	}

	var res stream.Endpoint
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return stream.Endpoint{}, err
	}

	return res, nil
}

func (c *Client) EditStreamEndpoint(ctx context.Context, apiKey string, endpointID int, options stream.EndpointAddEditDTO) (stream.Endpoint, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/streams/%d", c.BaseURL, endpointID))
	if err != nil {
		return stream.Endpoint{}, fmt.Errorf("invalid base URL: %w", err)
	}

	optionBytes, err := json.Marshal(options)
	if err != nil {
		return stream.Endpoint{}, err
	}

	buf := bytes.NewBuffer(optionBytes)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), buf)
	if err != nil {
		return stream.Endpoint{}, err
	}

	var res stream.Endpoint
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return stream.Endpoint{}, err
	}

	return res, nil
}

func (c *Client) DeleteStreamEndpoint(ctx context.Context, apiKey string, endpointID int) error {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/streams/%d", c.BaseURL, endpointID))
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	if err = c.sendRequest(req, apiKey, nil); err != nil {
		return err
	}

	return nil
}
