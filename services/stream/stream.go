package stream

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"

	"github.com/ystv/web-api/utils"
)

type (
	// Repo represents all stream endpoint interactions
	Repo interface {
		ListEndpoints(ctx context.Context) ([]EndpointDB, error)
		FindEndpoint(ctx context.Context, findEndpoint FindEndpoint) (EndpointDB, error)
		GetEndpointByID(ctx context.Context, endpointID int) (EndpointDB, error)
		GetEndpointByApplicationNamePwd(ctx context.Context, application, name, pwd string) (EndpointDB, error)

		SetEndpointActiveByID(ctx context.Context, endpointID int) error
		SetEndpointInactiveByApplicationNamePwd(ctx context.Context, application, name, pwd string) error

		AddEndpoint(ctx context.Context, endpointNew EndpointAddEditDTO) (EndpointDB, error)
		EditEndpoint(ctx context.Context, endpointID int, endpointEdit EndpointAddEditDTO) (EndpointDB, error)
		DeleteEndpoint(ctx context.Context, endpointID int) error
	}

	// EndpointDB stores a stream endpoint value
	EndpointDB struct {
		EndpointID  int         `json:"endpointId" db:"endpoint_id"`
		Application string      `json:"application" db:"application"`
		Name        string      `json:"name" db:"name"`
		Pwd         null.String `json:"pwd" db:"pwd"`
		StartValid  null.Time   `json:"startValid" db:"start_valid"`
		EndValid    null.Time   `json:"endValid" db:"end_valid"`
		Notes       null.String `json:"notes" db:"notes"`
		Active      bool        `json:"active" db:"active"`
		Blocked     bool        `json:"blocked" db:"blocked"`
		AutoRemove  bool        `json:"autoRemove" db:"auto_remove"`
	}

	// Endpoint returned endpoint value
	Endpoint struct {
		EndpointID int `json:"endpointId" db:"endpoint_id"`
		// Application defines which RTMP application this is valid for
		Application string `json:"application"`
		// Name is the unique name given in an application
		Name string `json:"name"`
		// Pwd defines an extra layer of security for authentication
		Pwd *string `json:"pwd,omitempty"`
		// StartValid defines the optional start time that this endpoint becomes valid
		StartValid *time.Time `json:"startValid,omitempty"`
		// EndValid defines the optional end time that this endpoint stops being valid
		EndValid *time.Time `json:"endValid,omitempty"`
		// Notes is an optional internal note for the endpoint
		Notes *string `json:"notes,omitempty"`
		// Active indicates if this endpoint is currently being used
		Active bool `json:"active"`
		// Blocked prevents the endpoint from going live
		Blocked bool `json:"blocked"`
		// AutoRemove indicates that this endpoint can be automatically removed when the end valid time comes, optional
		AutoRemove bool `json:"autoRemove,omitempty"`
	}

	// FindEndpoint used to find an endpoint
	FindEndpoint struct {
		// EndpointID is the unique database id of the stream
		EndpointID *int `json:"endpointId,omitempty"`
		// Application defines which RTMP application this is valid for
		Application *string `json:"application,omitempty"`
		// Name is the unique name given in an application
		Name *string `json:"name,omitempty"`
		// Pwd defines an extra layer of security for authentication
		Pwd *string `json:"pwd,omitempty"`
	}

	// EndpointAddEditDTO encapsulates the creation of a stream endpoint
	EndpointAddEditDTO struct {
		// Application defines which RTMP application this is valid for
		Application string `json:"application"`
		// Name is the unique name given in an application
		Name string `json:"name"`
		// Pwd defines an extra layer of security for authentication
		Pwd *string `json:"pwd,omitempty"`
		// StartValid defines the optional start time that this endpoint becomes valid, RFC3339
		StartValid *time.Time `json:"startValid,omitempty"`
		// EndValid defines the optional end time that this endpoint stops being valid, RFC3339
		EndValid *time.Time `json:"endValid,omitempty"`
		// Notes is an optional internal note for the endpoint
		Notes *string `json:"notes,omitempty"`
		// Blocked prevents the endpoint from going live, optional defaults to false
		Blocked bool `json:"blocked,omitempty"`
		// AutoRemove indicates that this endpoint can be automatically removed when the end valid time comes, optional
		AutoRemove bool `json:"autoRemove,omitempty"`
	}

	// Store encapsulates our dependency
	Store struct {
		db *sqlx.DB
	}
)

// NewStore creates our data store
func NewStore(db *sqlx.DB) Repo {
	return &Store{db}
}

func (s *Store) ListEndpoints(ctx context.Context) ([]EndpointDB, error) {
	var e []EndpointDB

	builder := utils.PSQL().Select("*").
		From("web_api.stream_endpoints").
		OrderBy("endpoint_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for GetEndpoints: %w", err))
	}

	err = s.db.SelectContext(ctx, &e, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %w", err)
	}

	return e, nil
}

func (s *Store) FindEndpoint(ctx context.Context, findEndpoint FindEndpoint) (EndpointDB, error) {
	var e EndpointDB

	builder := utils.PSQL().Select("*").
		From("web_api.stream_endpoints").
		Where(sq.Or{
			sq.Eq{"endpoint_id": findEndpoint.EndpointID},
			sq.And{
				sq.Eq{"application": findEndpoint.Application},
				sq.Eq{"name": findEndpoint.Name},
				sq.Eq{"pwd": findEndpoint.Pwd},
			},
		})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for FindEndpoint: %w", err))
	}

	err = s.db.GetContext(ctx, &e, sql, args...)
	if err != nil {
		return EndpointDB{}, fmt.Errorf("failed to find endpoint: %w", err)
	}

	findEndpoint.EndpointID = &e.EndpointID
	findEndpoint.Application = &e.Application
	findEndpoint.Name = &e.Name

	return e, nil
}

func (s *Store) GetEndpointByID(ctx context.Context, endpointID int) (EndpointDB, error) {
	var e EndpointDB

	builder := utils.PSQL().Select("*").
		From("web_api.stream_endpoints").
		Where(sq.Eq{"endpoint_id": endpointID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for GetEndpointByID: %w", err))
	}

	err = s.db.GetContext(ctx, &e, sql, args...)
	if err != nil {
		return EndpointDB{}, fmt.Errorf("failed to get endpoint by id: %w", err)
	}

	return e, nil
}

func (s *Store) GetEndpointByApplicationNamePwd(ctx context.Context, app, name, pwd string) (EndpointDB, error) {
	var e EndpointDB

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
		panic(fmt.Errorf("failed to build sql for GetEndpointByID: %w", err))
	}

	err = s.db.GetContext(ctx, &e, sql, args...)
	if err != nil {
		return EndpointDB{}, fmt.Errorf("failed to get endpoint by application name pwd: %w", err)
	}

	return e, nil
}

func (s *Store) SetEndpointActiveByID(ctx context.Context, endpointID int) error {
	builder := utils.PSQL().Update("web_api.stream_endpoints").
		Set("active", true).
		Where(sq.Eq{"endpoint_id": endpointID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for SetEndpointActiveByID: %w", err))
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
		panic(fmt.Errorf("failed to build sql for SetEndpointInactiveByApplicationNamePwd: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to set endpoint inactive by application name pwd: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to set endpoint inactive by application name pwd: %w", err)
	}

	if rows < 1 {
		return fmt.Errorf("failed to set endpoint inactive by application name pwd: invalid rows affected: %d", rows)
	}

	return nil
}

func (s *Store) AddEndpoint(ctx context.Context, endpointNew EndpointAddEditDTO) (EndpointDB, error) {
	builder := utils.PSQL().Insert("web_api.stream_endpoints").
		Columns("application", "name", "pwd", "start_valid", "end_valid", "notes", "active", "blocked",
			"auto_remove").
		Values(endpointNew.Application, endpointNew.Name, endpointNew.Pwd, endpointNew.StartValid, endpointNew.EndValid, endpointNew.Notes, false, endpointNew.Blocked, endpointNew.AutoRemove).
		Suffix("RETURNING endpoint_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for EndpointAddEditDTO: %w", err))
	}

	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return EndpointDB{}, fmt.Errorf("failed to add stream endpoint: %w", err)
	}

	defer stmt.Close()

	var endpointID int

	err = stmt.QueryRow(args...).Scan(&endpointID)
	if err != nil {
		return EndpointDB{}, fmt.Errorf("failed to add stream endpoint: %w", err)
	}

	return s.GetEndpointByID(ctx, endpointID)
}

func (s *Store) EditEndpoint(ctx context.Context, endpointID int, endpointEdit EndpointAddEditDTO) (EndpointDB, error) {
	builder := utils.PSQL().Update("web_api.stream_endpoints").
		SetMap(map[string]interface{}{
			"application": endpointEdit.Application,
			"name":        endpointEdit.Name,
			"pwd":         endpointEdit.Pwd,
			"start_valid": endpointEdit.StartValid,
			"end_valid":   endpointEdit.EndValid,
			"notes":       endpointEdit.Notes,
			"blocked":     endpointEdit.Blocked,
			"auto_remove": endpointEdit.AutoRemove,
		}).
		Where(sq.Eq{"endpoint_id": endpointID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for EditEndpoint: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return EndpointDB{}, fmt.Errorf("failed to edit endpoint: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return EndpointDB{}, fmt.Errorf("failed to edit endpoint: %w", err)
	}

	if rows < 1 {
		return EndpointDB{}, fmt.Errorf("failed to edit endpoint: invalid rows affected: %d, this endpoint may not exist: %d",
			rows, endpointID)
	}

	return s.GetEndpointByID(ctx, endpointID)
}

func (s *Store) DeleteEndpoint(ctx context.Context, endpointID int) error {
	builder := utils.PSQL().Delete("web_api.stream_endpoints").
		Where(sq.Eq{"endpoint_id": endpointID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for DeleteEndpoint: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete stream endpoint: %w", err)
	}
	return nil
}
