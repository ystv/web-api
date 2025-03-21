package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ystv/web-api/client/types"
	"github.com/ystv/web-api/services/people"
)

func (c *Client) ListUsersPagination(ctx context.Context, apiKey string, options types.ListUsersPaginationOptions) (people.UserFullPagination, error) {
	u, err := url.Parse(c.BaseURL + "/v1/internal/people/users/pagination")
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return people.UserFullPagination{}, err
	}

	var res people.UserFullPagination
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.UserFullPagination{}, err
	}

	return res, nil
}

func (c *Client) GetUser(ctx context.Context, apiKey string, userID int) (people.User, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v1/internal/people/user/%d", c.BaseURL, userID))
	if err != nil {
		return people.User{}, fmt.Errorf("invalid base URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return people.User{}, err
	}

	var res people.User
	if err = c.sendRequest(req, apiKey, &res); err != nil {
		return people.User{}, err
	}

	return res, nil
}
