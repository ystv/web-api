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

func (c *Client) ListRolesWithPermissions(ctx context.Context, apiKey string) ([]people.RoleWithPermissions, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/people/roles")
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var res []people.RoleWithPermissions
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) ListRolesWithCount(ctx context.Context, apiKey string) ([]people.RoleWithCount, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/people/roles/count")
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var res []people.RoleWithCount
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetRole(ctx context.Context, apiKey string, roleID int) (people.Role, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d", c.BaseURL, roleID))
	if err != nil {
		return people.Role{}, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return people.Role{}, err
	}

	var res people.Role
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.Role{}, err
	}

	return res, nil
}

func (c *Client) GetRoleFull(ctx context.Context, apiKey string, roleID int) (people.RoleFull, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d/full", c.BaseURL, roleID))
	if err != nil {
		return people.RoleFull{}, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return people.RoleFull{}, err
	}

	var res people.RoleFull
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.RoleFull{}, err
	}

	return res, nil
}

func (c *Client) ListRoleMembersByID(ctx context.Context, apiKey string, roleID int) ([]people.User, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d/members", c.BaseURL, roleID))
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

func (c *Client) ListRolePermissionsByID(ctx context.Context, apiKey string, roleID int) ([]people.Permission, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d/permissions", c.BaseURL, roleID))
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

func (c *Client) AddRole(ctx context.Context, apiKey string, options people.RoleAddEditDTO) (people.Role, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/people/role")
	if err != nil {
		return people.Role{}, fmt.Errorf("invalid base URL: %w", err)
	}

	optionBytes, err := json.Marshal(options)
	if err != nil {
		return people.Role{}, err
	}

	buf := bytes.NewBuffer(optionBytes)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), buf)
	if err != nil {
		return people.Role{}, err
	}

	var res people.Role
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.Role{}, err
	}

	return res, nil
}

func (c *Client) EditRole(ctx context.Context, apiKey string, roleID int, options people.RoleAddEditDTO) (people.Role, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d", c.BaseURL, roleID))
	if err != nil {
		return people.Role{}, fmt.Errorf("invalid base URL: %w", err)
	}

	optionBytes, err := json.Marshal(options)
	if err != nil {
		return people.Role{}, err
	}

	buf := bytes.NewBuffer(optionBytes)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), buf)
	if err != nil {
		return people.Role{}, err
	}

	var res people.Role
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.Role{}, err
	}

	return res, nil
}

func (c *Client) DeleteRole(ctx context.Context, apiKey string, roleID int) error {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d", c.BaseURL, roleID))
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

func (c *Client) ListUsersNotInRole(ctx context.Context, apiKey string, roleID int) ([]people.User, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d/users/notinrole", c.BaseURL, roleID))
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

func (c *Client) RoleAddUser(ctx context.Context, apiKey string, roleUser people.RoleUser) (people.RoleUser, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d/user/%d", c.BaseURL, roleUser.RoleID, roleUser.UserID))
	if err != nil {
		return people.RoleUser{}, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return people.RoleUser{}, err
	}

	var res people.RoleUser
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.RoleUser{}, err
	}

	return res, nil
}

func (c *Client) RoleRemoveUser(ctx context.Context, apiKey string, roleUser people.RoleUser) error {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d/user/%d", c.BaseURL, roleUser.RoleID, roleUser.UserID))
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

func (c *Client) ListPermissionsNotInRole(ctx context.Context, apiKey string, roleID int) ([]people.Permission, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d/permissions/notinrole", c.BaseURL, roleID))
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

func (c *Client) RoleAddPermission(ctx context.Context, apiKey string, rolePermission people.RolePermission) (people.RolePermission, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d/permission/%d", c.BaseURL, rolePermission.RoleID, rolePermission.PermissionID))
	if err != nil {
		return people.RolePermission{}, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return people.RolePermission{}, err
	}

	var res people.RolePermission
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.RolePermission{}, err
	}

	return res, nil
}

func (c *Client) RoleRemovePermission(ctx context.Context, apiKey string, rolePermission people.RolePermission) error {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/role/%d/permission/%d", c.BaseURL, rolePermission.RoleID, rolePermission.PermissionID))
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
