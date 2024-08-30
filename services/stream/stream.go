package stream

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"

	"github.com/ystv/web-api/utils"
)

type (
	// StreamRepo represents all stream endpoint interactions
	StreamRepo interface {
		GetEndpoints(ctx context.Context) ([]Endpoint, error)
		GetEndpointByID(ctx context.Context, endpointID string) (Endpoint, error)
		GetEndpointByApplicationNamePwd(ctx context.Context, application, name, pwd string) (Endpoint, error)
		SetEndpointActiveByID(ctx context.Context, endpointID string) error
		SetEndpointInactiveByApplicationNamePwd(ctx context.Context, application, name, pwd string) error
	}

	// Endpoint stores a stream endpoint value
	Endpoint struct {
		EndpointID  string      `json:"endpoint_id" db:"endpoint_id"`
		Application string      `json:"application" db:"application"`
		Name        string      `json:"name" db:"name"`
		Pwd         string      `json:"pwd" db:"pwd"`
		StartValid  null.Time   `json:"start_valid" db:"start_valid"`
		EndValid    null.Time   `json:"end_valid" db:"end_valid"`
		Notes       null.String `json:"notes" db:"notes"`
		Active      bool        `json:"active" db:"active"`
		Blocked     bool        `json:"blocked" db:"blocked"`
	}

	// Store encapsulates our dependency
	Store struct {
		db *sqlx.DB
	}
)

// NewStore creates our data store
func NewStore(db *sqlx.DB) *Store {
	return &Store{db}
}

func (s *Store) GetEndpoints(ctx context.Context) ([]Endpoint, error) {
	var e []Endpoint

	builder := utils.PSQL().Select("*").
		From("web_api.stream_endpoints").
		OrderBy("endpoint_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getEndpoints: %w", err))
	}

	err = s.db.SelectContext(ctx, &e, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %w", err)
	}

	return e, nil
}

func (s *Store) GetEndpointByID(ctx context.Context, endpointID string) (Endpoint, error) {
	var e Endpoint

	builder := utils.PSQL().Select("*").
		From("web_api.stream_endpoints").
		Where(sq.Eq{"endpoint_id": endpointID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getEndpointByID: %w", err))
	}

	err = s.db.GetContext(ctx, &e, sql, args...)
	if err != nil {
		return Endpoint{}, fmt.Errorf("failed to get endpoint by id: %w", err)
	}

	return e, nil
}

func (s *Store) GetEndpointByApplicationNamePwd(ctx context.Context, app, name, pwd string) (Endpoint, error) {
	var e Endpoint

	builder := utils.PSQL().Select("*").
		From("web_api.stream_endpoints").
		Where(sq.And{
			sq.Eq{"application": app},
			sq.Eq{"name": name},
			sq.Eq{"pwd": pwd},
		}).
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getEndpointByID: %w", err))
	}

	err = s.db.GetContext(ctx, &e, sql, args...)
	if err != nil {
		return Endpoint{}, fmt.Errorf("failed to get endpoint by id: %w", err)
	}

	return e, nil
}

func (s *Store) SetEndpointActiveByID(ctx context.Context, endpointID string) error {
	builder := utils.PSQL().Update("web_api.stream_endpoints").
		Set("active", true).
		Where(sq.Eq{"endpoint_id": endpointID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for setEndpointActiveByID: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to set endpoint active by id: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to set endpoint active by id: %w", err)
	}

	if rows < 1 {
		return fmt.Errorf("failed to set endpoint active by id: invalid rows affected: %d", rows)
	}

	return nil
}

func (s *Store) SetEndpointInactiveByApplicationNamePwd(ctx context.Context, application, name, pwd string) error {
	builder := utils.PSQL().Update("web_api.stream_endpoints").
		Set("active", false).
		Where(sq.And{
			sq.Eq{"application": application},
			sq.Eq{"name": name},
			sq.Eq{"pwd": pwd},
		})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for setEndpointInactiveByApplicationNamePwd: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to set endpoint inactive by application name pwd: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to set endpoint active by application name pwd: %w", err)
	}

	if rows < 1 {
		return fmt.Errorf("failed to set endpoint inactive by application name pwd: invalid rows affected: %d", rows)
	}

	return nil
}
