package client

import (
	"bytes"
	"context"
	"encoding/json"
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

func (c *Client) ListPermissionsWithRolesCount(ctx context.Context, apiKey string) ([]people.PermissionWithRolesCount, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/people/permissions/count")
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var res []people.PermissionWithRolesCount
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) ListPermissionMembersByID(ctx context.Context, apiKey string, permissionID int) ([]people.User, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/permission/%d/members", c.BaseURL, permissionID))
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var res []people.User
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

func (c *Client) GetPermissionWithRolesCount(ctx context.Context, apiKey string, permissionID int) (people.PermissionWithRolesCount, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/permission/%d/count", c.BaseURL, permissionID))
	if err != nil {
		return people.PermissionWithRolesCount{}, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return people.PermissionWithRolesCount{}, err
	}

	var res people.PermissionWithRolesCount
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.PermissionWithRolesCount{}, err
	}

	return res, nil
}

func (c *Client) NewPermission(ctx context.Context, apiKey string, options people.PermissionAddEditDTO) (people.Permission, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/people/permission")
	if err != nil {
		return people.Permission{}, fmt.Errorf("invalid base URL: %w", err)
	}

	optionBytes, err := json.Marshal(options)
	if err != nil {
		return people.Permission{}, err
	}

	buf := bytes.NewBuffer(optionBytes)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), buf)
	if err != nil {
		return people.Permission{}, err
	}

	var res people.Permission
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.Permission{}, err
	}

	return res, nil
}

func (c *Client) EditPermission(ctx context.Context, apiKey string, permissionID int, options people.PermissionAddEditDTO) (people.Permission, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/permission/%d", c.BaseURL, permissionID))
	if err != nil {
		return people.Permission{}, fmt.Errorf("invalid base URL: %w", err)
	}

	optionBytes, err := json.Marshal(options)
	if err != nil {
		return people.Permission{}, err
	}

	buf := bytes.NewBuffer(optionBytes)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), buf)
	if err != nil {
		return people.Permission{}, err
	}

	var res people.Permission
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.Permission{}, err
	}

	return res, nil
}

func (c *Client) DeletePermission(ctx context.Context, apiKey string, permissionID int) error {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/permission/%d", c.BaseURL, permissionID))
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
