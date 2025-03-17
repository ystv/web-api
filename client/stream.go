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
