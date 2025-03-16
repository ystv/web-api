package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ystv/web-api/services/people"
)

func (c *Client) ListPermissions(ctx context.Context, apiKey string) ([]people.Permission, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/people/permissions")
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var res []people.Permission
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetPermission(ctx context.Context, apiKey string, permissionID int) (people.Permission, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/permission/%d", c.BaseURL, permissionID))
	if err != nil {
		return people.Permission{}, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return people.Permission{}, err
	}

	var res people.Permission
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.Permission{}, err
	}

	return res, nil
}
