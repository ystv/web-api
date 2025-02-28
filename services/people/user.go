package people

import (
	"context"
	//nolint:gosec
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"

	"github.com/ystv/web-api/utils"
)

// CountUsersAll will get the number of total users
func (s *Store) CountUsersAll(ctx context.Context) (CountUsers, error) {
	var countUsers CountUsers

	err := s.db.GetContext(ctx, &countUsers,
		`SELECT
		(SELECT COUNT(*) FROM people.users) as total_users,
		(SELECT COUNT(*) FROM people.users WHERE enabled = true AND deleted_by IS NULL AND deleted_at IS NULL)
		    AS active_users,
		(SELECT COUNT(*) FROM people.users WHERE last_login > TO_TIMESTAMP($1, 'YYYY-MM-DD HH24:MI:SS'))
		    AS active_users_past_24_hours,
		(SELECT COUNT(*) FROM people.users WHERE last_login > TO_TIMESTAMP($2, 'YYYY-MM-DD HH24:MI:SS'))
		    AS active_users_past_year;`,
		time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05"),
		time.Now().AddDate(-1, 0, 0).Format("2006-01-02 15:04:05"))

	if err != nil {
		return countUsers, fmt.Errorf("failed to count users all from db: %w", err)
	}

	return countUsers, nil
}

// GetUserFull will return all user information to be used for profile and management.
func (s *Store) GetUserFull(ctx context.Context, userID int) (UserFull, error) {
	var u UserFull
	//nolint:musttag
	err := s.db.GetContext(ctx, &u, `
		SELECT user_id, username, email, first_name, last_name, nickname,
		avatar, use_gravatar, last_login, created_at, created_by, updated_at, updated_by,
		deleted_at, deleted_by
		FROM people.users
		WHERE user_id = $1
		LIMIT 1;`, userID)
	if err != nil {
		return UserFull{}, fmt.Errorf("failed to get user meta: %w", err)
	}

	//nolint:musttag
	err = s.db.SelectContext(ctx, &u.Roles, `
		SELECT r.role_id, r.name, r.description
		FROM people.roles r
		INNER JOIN people.role_members rm ON rm.role_id = r.role_id
		WHERE user_id = $1;`, userID)
	if err != nil {
		return UserFull{}, fmt.Errorf("failed to get roles: %w", err)
	}

	err = s.db.SelectContext(ctx, &u.Permissions, `
		SELECT DISTINCT p.permission_id, p.name, p.description
		FROM people.permissions p
		INNER JOIN people.role_permissions rp on rp.permission_id = p.permission_id
		INNER JOIN people.role_members rm ON rm.role_id = rp.role_id
		WHERE rm.user_id = $1
		ORDER BY name;`, userID)
	if err != nil {
		return UserFull{}, fmt.Errorf("failed to get permissions: %w", err)
	}

	switch avatar := u.Avatar; {
	case u.UseGravatar:
		//nolint:gosec
		hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(u.Email))))
		u.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
	case avatar == "":
		u.Avatar = "https://auth.ystv.co.uk/public/ystv-colour-white.png"
	case strings.Contains(avatar, s.cdnEndpoint):
	case strings.Contains(avatar, fmt.Sprintf("%d.", u.UserID)):
		u.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
	default:
		log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", u.UserID, len(u.Avatar), u.Avatar)
		u.Avatar = ""
	}

	return u, nil
}

// GetUserByEmailFull will return all user information to be used for profile and management.
func (s *Store) GetUserByEmailFull(ctx context.Context, email string) (UserFull, error) {
	var u UserFull
	//nolint:musttag
	err := s.db.GetContext(ctx, &u,
		`SELECT user_id, username, email, first_name, last_name, nickname,
		avatar, use_gravatar, last_login, created_at, created_by, updated_at, updated_by,
		deleted_at, deleted_by
		FROM people.users
		WHERE email = $1
		LIMIT 1;`, email)
	if err != nil {
		return UserFull{}, fmt.Errorf("failed to get user meta: %w", err)
	}

	//nolint:musttag
	err = s.db.SelectContext(ctx, &u.Roles,
		`SELECT r.role_id, r.name, r.description
	FROM people.roles r
	INNER JOIN people.role_members rm ON rm.role_id = r.role_id
	WHERE user_id = $1;`, u.UserID)
	if err != nil {
		return UserFull{}, fmt.Errorf("failed to get roles: %w", err)
	}

	err = s.db.SelectContext(ctx, &u.Permissions, `
		SELECT DISTINCT p.permission_id, p.name, p.description
		FROM people.permissions p
		INNER JOIN people.role_permissions rp on rp.permission_id = p.permission_id
		INNER JOIN people.role_members rm ON rm.role_id = rp.role_id
		WHERE rm.user_id = $1
		ORDER BY name;`, u.UserID)
	if err != nil {
		return UserFull{}, fmt.Errorf("failed to get permissions: %w", err)
	}

	switch avatar := u.Avatar; {
	case u.UseGravatar:
		//nolint:gosec
		hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(u.Email))))
		u.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
	case avatar == "":
		u.Avatar = "https://auth.ystv.co.uk/public/ystv-colour-white.png"
	case strings.Contains(avatar, s.cdnEndpoint):
	case strings.Contains(avatar, fmt.Sprintf("%d.", u.UserID)):
		u.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
	default:
		log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", u.UserID, len(u.Avatar), u.Avatar)
		u.Avatar = ""
	}

	return u, nil
}

// GetUser returns basic user information to be used for other services.
func (s *Store) GetUser(ctx context.Context, userID int) (User, error) {
	var u User
	//nolint:musttag
	err := s.db.GetContext(ctx, &u,
		`SELECT user_id, username, email, first_name, last_name, nickname, avatar, use_gravatar
		FROM people.users
		WHERE user_id = $1;`, userID)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user meta: %w", err)
	}

	err = s.db.SelectContext(ctx, &u.Permissions,
		`SELECT p.permission_id, p.name
		FROM people.permissions p
		INNER JOIN people.role_permissions rp ON rp.permission_id = p.permission_id
		INNER JOIN people.role_members rm ON rm.role_id = rp.role_id
		WHERE rm.user_id = $1;`, userID)
	if err != nil {
		return User{}, fmt.Errorf("failed to get permissions: %w", err)
	}

	switch avatar := u.Avatar; {
	case u.UseGravatar:
		//nolint:gosec
		hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(u.Email))))
		u.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
	case avatar == "":
		u.Avatar = "https://auth.ystv.co.uk/public/ystv-colour-white.png"
	case strings.Contains(avatar, s.cdnEndpoint):
	case strings.Contains(avatar, fmt.Sprintf("%d.", u.UserID)):
		u.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
	default:
		log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", u.UserID, len(u.Avatar), u.Avatar)
		u.Avatar = ""
	}

	return u, nil
}

// GetUserByEmail returns basic user information to be used for other services.
func (s *Store) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	//nolint:musttag
	err := s.db.GetContext(ctx, &u,
		`SELECT user_id, username, email, first_name, last_name, nickname, avatar, use_gravatar
		FROM people.users
		WHERE email = $1;`, email)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user meta: %w", err)
	}

	err = s.db.SelectContext(ctx, &u.Permissions,
		`SELECT p.permission_id, p.name
		FROM people.permissions p
		INNER JOIN people.role_permissions rp ON rp.permission_id = p.permission_id
		INNER JOIN people.role_members rm ON rm.role_id = rp.role_id
		WHERE rm.user_id = $1;`, u.UserID)
	if err != nil {
		return User{}, fmt.Errorf("failed to get permissions: %w", err)
	}

	switch avatar := u.Avatar; {
	case u.UseGravatar:
		//nolint:gosec
		hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(u.Email))))
		u.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
	case avatar == "":
		u.Avatar = "https://auth.ystv.co.uk/public/ystv-colour-white.png"
	case strings.Contains(avatar, s.cdnEndpoint):
	case strings.Contains(avatar, fmt.Sprintf("%d.", u.UserID)):
		u.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
	default:
		log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", u.UserID, len(u.Avatar), u.Avatar)
		u.Avatar = ""
	}

	return u, nil
}

// ListAllUsers returns all users.
// It doesn't return the full User object.
// Returns user_id, avatar, nickname, first_name, last_name.
//
// There will likely be modifications to include the other fields
// but will need to add a filter at the web handler first or offer
// a different function.
func (s *Store) ListAllUsers(ctx context.Context) ([]User, error) {
	var u []User
	//nolint:musttag
	err := s.db.SelectContext(ctx, &u,
		`SELECT user_id, avatar, use_gravatar, nickname, first_name, last_name
		FROM people.users;`)
	if err != nil {
		return nil, fmt.Errorf("fialed to list all users: %w", err)
	}

	return u, nil
}

// GetUsersPagination will get users search with sorting with size and page, enabled and deleted
// Use the parameter direction for determining of the sorting will be ascending(asc) or descending(desc)
func (s *Store) GetUsersPagination(ctx context.Context, size, page int, search, sortBy, direction, enabled,
	deleted string) ([]UserFull, int, error) {
	var u []UserFull

	var count int

	builder, err := s._getUsersBuilder(size, page, search, sortBy, direction, enabled, deleted)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to build sql for getUsers: %w", err)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for GetUsersPagination: %w", err))
	}

	rows, err := s.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to get db users: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	type tempStruct struct {
		UserFull
		Count int `db:"full_count" json:"fullCount"`
	}

	for rows.Next() {
		var u1 UserFull

		var temp tempStruct

		//nolint:musttag
		err = rows.StructScan(&temp)
		if err != nil {
			return nil, -1, fmt.Errorf("failed to get db users: %w", err)
		}

		count = temp.Count

		err = copier.Copy(&u1, &temp)
		if err != nil {
			return nil, -1, fmt.Errorf("failed to copy struct: %w", err)
		}

		u = append(u, u1)
	}

	return u, count, nil
}

func (s *Store) _getUsersBuilder(size, page int, search, sortBy, direction, enabled,
	deleted string) (*sq.SelectBuilder, error) {
	builder := utils.PSQL().Select("user_id", "username", "first_name", "nickname", "last_name", "email",
		"last_login", "enabled", "created_at", "created_by", "updated_at", "updated_by", "deleted_at", "deleted_by",
		"count(*) OVER() AS full_count").From("people.users")

	if len(search) > 0 {
		builder = builder.Where(
			"(CAST(user_id AS TEXT) LIKE '%' || ? || '%' "+
				"OR LOWER(username) LIKE LOWER('%' || ? || '%') "+
				"OR LOWER(nickname) LIKE LOWER('%' || ? || '%') "+
				"OR LOWER(first_name) LIKE LOWER('%' || ? || '%') "+
				"OR LOWER(last_name) LIKE LOWER('%' || ? || '%') "+
				"OR LOWER(email) LIKE LOWER('%' || ? || '%') "+
				"OR LOWER(first_name || ' ' || last_name) LIKE LOWER('%' || ? || '%'))",
			search, search, search, search, search, search, search)
	}

	switch enabled {
	case "enabled":
		builder = builder.Where(sq.Eq{"enabled": true})
	case "disabled":
		builder = builder.Where(sq.Eq{"enabled": false})
	}

	switch deleted {
	case "not_deleted":
		builder = builder.Where(sq.Eq{"deleted_by": nil})
	case "deleted":
		builder = builder.Where(sq.NotEq{"deleted_by": nil})
	}

	if len(sortBy) > 0 && len(direction) > 0 {
		switch direction {
		case "asc":
			builder = builder.OrderByClause(
				"CASE WHEN ? = 'userId' THEN user_id END ASC, "+
					"CASE WHEN ? = 'name' THEN first_name END ASC, "+
					"CASE WHEN ? = 'name' THEN last_name END ASC, "+
					"CASE WHEN ? = 'username' THEN username END ASC, "+
					"CASE WHEN ? = 'email' THEN email END ASC, "+
					"CASE WHEN ? = 'lastLogin' THEN last_login END ASC NULLS FIRST",
				sortBy, sortBy, sortBy, sortBy, sortBy, sortBy)
		case "desc":
			builder = builder.OrderByClause(
				"CASE WHEN ? = 'userId' THEN user_id END DESC, "+
					"CASE WHEN ? = 'name' THEN first_name END DESC, "+
					"CASE WHEN ? = 'name' THEN last_name END DESC, "+
					"CASE WHEN ? = 'username' THEN username END DESC, "+
					"CASE WHEN ? = 'email' THEN email END DESC, "+
					"CASE WHEN ? = 'lastLogin' THEN last_login END DESC NULLS LAST",
				sortBy, sortBy, sortBy, sortBy, sortBy, sortBy)
		default:
			return nil, fmt.Errorf(`invalid sorting direction, entered "%s" of length %d, but expected either 
"direction" or "desc"`, direction, len(direction))
		}
	}

	if page >= 1 && size >= 5 && size <= 100 {
		parsed1, err := strconv.ParseUint(strconv.Itoa(size), 10, 64)
		if err != nil {
			return nil, fmt.Errorf(`invalid value for size in direction "%s"`, direction)
		}
		parsed2, err := strconv.ParseUint(strconv.Itoa(size*(page-1)), 10, 64)
		if err != nil {
			return nil, fmt.Errorf(`invalid value for page in direction "%s"`, direction)
		}
		builder = builder.Limit(parsed1).Offset(parsed2)
	}

	return &builder, nil
}
