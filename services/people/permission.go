package people

import (
	"context"
	//nolint:gosec
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/ystv/web-api/utils"
)

func (s *Store) ListAllPermissions(ctx context.Context) ([]Permission, error) {
	var p []Permission

	err := s.db.SelectContext(ctx, &p, `
		SELECT permission_id, name, description
		FROM people.permissions;`)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return p, nil
}

func (s *Store) ListPermissionsWithRolesCount(ctx context.Context) ([]PermissionWithRolesCount, error) {
	var p []PermissionWithRolesCount

	builder := utils.PSQL().Select("p.*", "COUNT(rp.role_id) AS roles").
		From("people.permissions p").
		LeftJoin("people.role_permissions rp on p.permission_id = rp.permission_id").
		GroupBy("p", "p.permission_id", "name", "description").
		OrderBy("p.name")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for ListAllPermissionsWithRolesCount: %w", err))
	}

	err = s.db.SelectContext(ctx, &p, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return p, nil
}

func (s *Store) ListPermissionMembersByID(ctx context.Context, permissionID int) ([]User, error) {
	var u []User
	//nolint:musttag
	err := s.db.GetContext(ctx, &u, `
		SELECT u.user_id, username, email, first_name, last_name, nickname, avatar
		FROM people.users u
		INNER JOIN people.role_members rm ON u.user_id = rm.user_id
		INNER JOIN people.role_permissions p ON rm.role_id = p.role_id
		WHERE permission_id = $1;`, permissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user meta: %w", err)
	}

	for _, user := range u {
		switch avatar := user.Avatar; {
		case user.UseGravatar:
			//nolint:gosec
			hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(user.Email))))
			user.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
		case avatar == "", strings.Contains(avatar, s.cdnEndpoint):
		case strings.Contains(avatar, fmt.Sprintf("%d.", user.UserID)):
			user.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
		default:
			log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", user.UserID, len(user.Avatar), user.Avatar)
			user.Avatar = ""
		}
	}

	return u, nil
}

// GetPermission returns a permission
func (s *Store) GetPermission(ctx context.Context, permissionID int) (Permission, error) {
	var p Permission

	builder := utils.PSQL().Select("*").
		From("people.permissions").
		Where(sq.Eq{"permission_id": permissionID}).
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for GetPermission: %w", err))
	}

	err = s.db.GetContext(ctx, &p, sql, args...)
	if err != nil {
		return Permission{}, fmt.Errorf("failed to get permission: %w", err)
	}

	return p, nil
}

// GetPermissionWithRolesCount returns a permission with a roles count
func (s *Store) GetPermissionWithRolesCount(ctx context.Context, permissionID int) (PermissionWithRolesCount, error) {
	var p PermissionWithRolesCount

	builder := utils.PSQL().Select("p.*", "COUNT(rp.role_id) AS roles").
		From("people.permissions p").
		LeftJoin("people.role_permissions rp on p.permission_id = rp.permission_id").
		Where(sq.Eq{"p.permission_id": permissionID}).
		GroupBy("p", "p.permission_id", "name", "description").
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for GetPermission: %w", err))
	}

	err = s.db.GetContext(ctx, &p, sql, args...)
	if err != nil {
		return PermissionWithRolesCount{}, fmt.Errorf("failed to get permission: %w", err)
	}

	return p, nil
}

// AddPermission adds a new permission
func (s *Store) AddPermission(ctx context.Context, permission PermissionAddEditDTO) (Permission, error) {
	var addedPermission Permission

	builder := utils.PSQL().Insert("people.permissions").
		Columns("name", "description").
		Values(permission.Name, permission.Description).
		Suffix("RETURNING permission_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for addPermission: %w", err))
	}

	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return Permission{}, fmt.Errorf("failed to add permission: %w", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(args...).Scan(&addedPermission.PermissionID)
	if err != nil {
		return Permission{}, fmt.Errorf("failed to add permission: %w", err)
	}

	addedPermission.Name = permission.Name
	addedPermission.Description = permission.Description

	return addedPermission, nil
}

// EditPermission edits an existing permission
func (s *Store) EditPermission(ctx context.Context, permissionID int, permission PermissionAddEditDTO) (Permission, error) {
	var editedPermission Permission

	builder := utils.PSQL().Update("people.permissions").
		SetMap(map[string]interface{}{
			"name":        permission.Name,
			"description": permission.Description,
		}).
		Where(sq.Eq{"permission_id": permissionID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for EditPermission: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return Permission{}, fmt.Errorf("failed to edit permission: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return Permission{}, fmt.Errorf("failed to edit permission: %w", err)
	}

	if rows < 1 {
		return Permission{}, fmt.Errorf("failed to edit permissions: invalid rows affected: %d", rows)
	}

	editedPermission = Permission{
		PermissionID: permissionID,
		Name:         permission.Name,
		Description:  permission.Description,
	}

	return editedPermission, nil
}

// DeletePermission deletes a specific permission
func (s *Store) DeletePermission(ctx context.Context, permissionID int) error {
	err := s.removePermissionForRoles(ctx, permissionID)
	if err != nil {
		return fmt.Errorf("failed to delete role permission link: %w", err)
	}
	builder := utils.PSQL().Delete("people.permissions").
		Where(sq.Eq{"permission_id": permissionID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for DeletePermission: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return nil
}

// removePermissionForRoles deletes the connection between multiple Role and a Permission
func (s *Store) removePermissionForRoles(ctx context.Context, permissionID int) error {
	builder := utils.PSQL().Delete("people.role_permissions").
		Where(sq.Eq{"permission_id": permissionID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for removePermissionForRoles: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete rolePermission: %w", err)
	}

	return nil
}
