package misc

import (
	"context"
	"fmt"

	"github.com/ystv/web-api/services/people"
)

type (
	// List is a mailing list which people can subscribe to
	List struct {
		ListID       int          `db:"list_id" json:"listID"`
		Name         string       `db:"name" json:"name"`
		Description  string       `db:"description" json:"description"`
		Alias        string       `db:"alias" json:"alias"`
		PermissionID *int         `db:"permission_id" json:"permissionID"`
		IsSubscribed bool         `db:"is_subscribed" json:"isSubscribed"`
		Subscribers  []Subscriber `json:"subscribers,omitempty"`
	}
	// Subscriber is an individual user on a mailing list
	Subscriber struct {
		SubscribeID string `db:"subscribe_id" json:"subscribeID"`
		people.User `json:"user"`
	}
)

var _ ListRepo = &Store{}

// GetLists returns all available mailing lists
// Doesn't include individual subscribers
func (m *Store) GetLists(ctx context.Context) ([]List, error) {
	var l []List
	err := m.db.SelectContext(ctx, &l, `
		SELECT list_id, name, description, alias, permission_id
		FROM mail.lists;`)
	if err != nil {
		return nil, fmt.Errorf("failed to get mailing lists: %w", err)
	}
	return l, nil
}

// GetListsByUserID returns all available and currently subscribed to
// mailing lists for a user don't include individual subscribers
func (m *Store) GetListsByUserID(ctx context.Context, userID int) ([]List, error) {
	var l []List
	err := m.db.SelectContext(ctx, &l, `
		SELECT DISTINCT list.list_id, name, description, alias, permission_id,
		CASE WHEN sub.user_id = $1 THEN true ELSE false END AS is_subscribed
		FROM mail.lists list
		INNER JOIN mail.subscribers sub ON list.list_id = sub.list_id;`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get mailing lists by user ID: %w", err)
	}
	return l, nil
}

// GetList returns a list including all subscribers
func (m *Store) GetList(ctx context.Context, listID int) (List, error) {
	l := List{}
	err := m.db.GetContext(ctx, &l, `
		SELECT list_id, name, description, alias, permission_id
		FROM mail.lists
		WHERE list_id = $1
		LIMIT 1;`, listID)
	if err != nil {
		return l, fmt.Errorf("failed to get list meta: %w", err)
	}
	l.Subscribers, err = m.GetSubscribers(ctx, listID)
	if err != nil {
		return l, fmt.Errorf("failed to get list subscribers: %w", err)
	}
	return l, nil
}

// GetSubscribers returns all subscribers of a list
func (m *Store) GetSubscribers(ctx context.Context, listID int) ([]Subscriber, error) {
	var s []Subscriber
	err := m.db.SelectContext(ctx, &s, `
		SELECT subscribe_id, sub.user_id, username, email, first_name, last_name, nickname, avatar
		FROM mail.subscribers sub
		INNER JOIN people.users u ON sub.user_id = u.user_id;`)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscribers: %w", err)
	}
	return s, nil
}

// Subscribe adds a user to a mailing list
func (m *Store) Subscribe(ctx context.Context, userID, listID int) error {
	_, err := m.db.ExecContext(ctx, `INSERT INTO mail.subscribers(list_id, user_id VALUES ($1, $2);`, listID, userID)
	if err != nil {
		return fmt.Errorf("failed to subscribe to mailing list \"%d\": %w", listID, err)
	}
	return nil
}

// UnsubscribeByID removes a user from a mailing list by userID and listID
func (m *Store) UnsubscribeByID(ctx context.Context, userID, listID int) error {
	_, err := m.db.ExecContext(ctx, `
		DELETE FROM mail.subscribers
		WHERE list_id = $1 AND user_id = $2`, listID, userID)
	if err != nil {
		return fmt.Errorf("failed to unscribe user \"%d\" from list \"%d\" %w", userID, listID, err)
	}
	return nil
}

// UnsubscribeByUUID removes a user from a mailing list by the subscriber UUID
func (m *Store) UnsubscribeByUUID(ctx context.Context, uuid string) error {
	_, err := m.db.ExecContext(ctx, `
		DELETE FROM mail.subscribers
		WHERE subscriber_id = $1`, uuid)
	if err != nil {
		return fmt.Errorf("failed to unscribe user by uuid \"%s\": %w", uuid, err)
	}
	return nil
}
