package people

import (
	"context"
	"fmt"
	"net/http"

	sq "github.com/Masterminds/squirrel"

	"github.com/ystv/web-api/utils"
)

func (s *Store) ListAllRolesWithPermissions(ctx context.Context) ([]RoleWithPermissions, error) {
	var temp []RoleWithPermissions
	r := make([]RoleWithPermissions, 0)
	//nolint:musttag
	err := s.db.SelectContext(ctx, &temp, `
		SELECT role_id, name, description
		FROM people.roles;`)
	if err != nil {
		return nil, fmt.Errorf("failed to select roles: %w", err)
	}

	for _, role := range temp {
		err = s.db.SelectContext(ctx, &role.Permissions, `
			SELECT perm.permission_id, name, description
			FROM people.permissions perm
			INNER JOIN people.role_permissions role ON perm.permission_id = role.permission_id
			WHERE role_id = $1;`, role.RoleID)
		if err != nil {
			return nil, fmt.Errorf("failed to get permissions for role \"%d\": %w", role.RoleID, err)
		}

		r = append(r, role)
	}

	return r, nil
}

func (s *Store) ListAllRolesWithCount(ctx context.Context) ([]RoleWithCount, error) {
	var r []RoleWithCount
	//nolint:musttag
	err := s.db.SelectContext(ctx, &r, `SELECT r.*, COUNT(DISTINCT rm.user_id) AS users, 
        COUNT(DISTINCT rp.permission_id) AS permissions
		FROM people.roles r
		LEFT JOIN people.role_members rm ON r.role_id = rm.role_id
		LEFT JOIN people.role_permissions rp ON r.role_id = rp.role_id
		GROUP BY r, r.role_id, name, description
		ORDER BY r.name`)
	if err != nil {
		return nil, fmt.Errorf("failed to select roles: %w", err)
	}

	return r, nil
}

func (s *Store) GetRole(ctx context.Context, roleGetDTO RoleGetDTO) (Role, error) {
	var r Role

	builder := utils.PSQL().Select("*").
		From("people.roles").
		Where(sq.Or{
			sq.Eq{"role_id": roleGetDTO.RoleID},
			sq.And{
				sq.Eq{"name": roleGetDTO.Name},
				sq.NotEq{"name": ""},
			},
		}).
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getRole: %w", err))
	}

	err = s.db.GetContext(ctx, &r, sql, args...)
	if err != nil {
		return Role{}, fmt.Errorf("failed to get role: %w", err)
	}

	return r, nil
}

// GetRoleFull returns all users and permissions of a certain role.
// It doesn't return the full User object.
// Returns user_id, avatar, nickname, first_name, last_name.
//
// There will likely be modifications to include the other fields
// but will need to add a filter at the web handler first or offer
// a different function.
func (s *Store) GetRoleFull(ctx context.Context, roleID int) (RoleFull, error) {
	var r RoleFull
	//nolint:musttag
	err := s.db.GetContext(ctx, &r, `
		SELECT role_id, name, description
		FROM people.roles
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return RoleFull{}, fmt.Errorf("failed to select roles: %w", err)
	}

	err = s.db.SelectContext(ctx, &r.Permissions, `
		SELECT perm.permission_id, name, description
		FROM people.permissions perm
		INNER JOIN people.role_permissions role ON perm.permission_id = role.permission_id
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return RoleFull{}, fmt.Errorf("failed to get permissions for role \"%d\": %w", r.RoleID, err)
	}

	//nolint:musttag
	err = s.db.SelectContext(ctx, &r.Users, `
		SELECT u.user_id, avatar, nickname, first_name, last_name
		FROM people.users u
		INNER JOIN people.role_members member ON u.user_id = member.user_id
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return RoleFull{}, fmt.Errorf("failed to get users for role: %w", err)
	}

	return r, nil
}

// ListRoleMembersByID returns all users who have a certain role.
// It doesn't return the full User object.
// Returns user_id, avatar, nickname, first_name, last_name.
//
// There will likely be modifications to include the other fields
// but will need to add a filter at the web handler first or offer
// a different function.
func (s *Store) ListRoleMembersByID(ctx context.Context, roleID int) ([]User, int, error) {
	err := s.validateRole(ctx, roleID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("failed to find role: %w", err)
	}

	var u []User
	//nolint:musttag
	err = s.db.SelectContext(ctx, &u,
		`SELECT u.user_id, avatar, nickname, first_name, last_name
		FROM people.users u
		INNER JOIN people.role_members rm ON u.user_id = rm.user_id
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to list role users: %w", err)
	}

	return u, http.StatusOK, nil
}

func (s *Store) ListRolePermissionsByID(ctx context.Context, roleID int) ([]Permission, error) {
	var p []Permission

	err := s.db.SelectContext(ctx, &p, `
		SELECT perms.permission_id, perms.name, perms.description
		FROM people.permissions perms
		INNER JOIN people.role_permissions rp ON perms.permission_id = rp.permission_id
		WHERE rp.role_id = $1`, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions for role \"%d\": %w", roleID, err)
	}

	return p, nil
}

func (s *Store) validateRole(ctx context.Context, roleID int) error {
	var r RoleFull
	//nolint:musttag
	err := s.db.GetContext(ctx, &r, `
		SELECT role_id, name, description
		FROM people.roles
		WHERE role_id = $1;`, roleID)
	if err != nil {
		return fmt.Errorf("failed to get roles: %w", err)
	}

	return nil
}

func (s *Store) AddRole(ctx context.Context, roleAdd RoleAddEditDTO) (Role, error) {
	builder := utils.PSQL().Insert("people.roles").
		Columns("name", "description").
		Values(roleAdd.Name, roleAdd.Description).
		Suffix("RETURNING role_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for addRole: %w", err))
	}

	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return Role{}, fmt.Errorf("failed to add role: %w", err)
	}

	defer stmt.Close()

	var roleID int

	err = stmt.QueryRow(args...).Scan(&roleID)
	if err != nil {
		return Role{}, fmt.Errorf("failed to add role: %w", err)
	}

	return s.GetRole(ctx, RoleGetDTO{RoleID: roleID})
}

func (s *Store) EditRole(ctx context.Context, roleID int, roleEdit RoleAddEditDTO) (Role, error) {
	builder := utils.PSQL().Update("people.roles").
		SetMap(map[string]interface{}{
			"name":        roleEdit.Name,
			"description": roleEdit.Description,
		}).
		Where(sq.Eq{"role_id": roleID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for editRole: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return Role{}, fmt.Errorf("failed to edit role: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return Role{}, fmt.Errorf("failed to edit role: %w", err)
	}

	if rows < 1 {
		return Role{}, fmt.Errorf("failed to edit role: invalid rows affected: %d", rows)
	}

	return s.GetRole(ctx, RoleGetDTO{RoleID: roleID})
}

func (s *Store) DeleteRole(ctx context.Context, roleID int) error {
	builder := utils.PSQL().Delete("people.roles").
		Where(sq.Eq{"role_id": roleID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for deleteRole: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// RemoveRoleForPermissions deletes links between a Role and Permissions
func (s *Store) RemoveRoleForPermissions(ctx context.Context, roleID int) error {
	builder := utils.PSQL().Delete("people.role_permissions").
		Where(sq.Eq{"role_id": roleID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for removeRoleForPermissions: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete rolePermission: %w", err)
	}

	return nil
}

// RemoveRoleForUsers deletes links between a Role and Users
func (s *Store) RemoveRoleForUsers(ctx context.Context, roleID int) error {
	builder := utils.PSQL().Delete("people.role_members").
		Where(sq.Eq{"role_id": roleID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for removeRoleForUsers: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete roleUser: %w", err)
	}

	return nil
}

// GetRoleUser returns a role user - moved here for cycle import reasons
func (s *Store) GetRoleUser(ctx context.Context, ru1 RoleUser) (RoleUser, error) {
	var ru RoleUser

	builder := utils.PSQL().Select("*").
		From("people.role_members").
		Where(sq.And{
			sq.Eq{"role_id": ru1.RoleID},
			sq.Eq{"user_id": ru1.UserID},
		}).
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getRoleUser: %w", err))
	}

	err = s.db.GetContext(ctx, &ru, sql, args...)
	if err != nil {
		return RoleUser{}, fmt.Errorf("failed to get role user: %w", err)
	}

	return ru, nil
}

// GetUsersNotInRole returns all the users not currently in the role.Role to be added
func (s *Store) GetUsersNotInRole(ctx context.Context, roleID int) ([]User, error) {
	var u []User

	subQuery := utils.PSQL().Select("u.user_id").
		From("people.users u").
		LeftJoin("people.role_members ru on u.user_id = ru.user_id").
		Where(sq.Eq{"ru.role_id": roleID})

	builder := utils.PSQL().Select("u.*").
		Distinct().
		From("people.users u").
		Where(sq.And{
			utils.NotIn("user_id", subQuery),
			sq.Eq{"deleted_by": nil},
		}).
		OrderBy("first_name", "last_name")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getRoles: %w", err))
	}

	//nolint:musttag
	err = s.db.SelectContext(ctx, &u, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	return u, nil
}

// AddRoleUser creates a link between a role.Role and User
func (s *Store) AddRoleUser(ctx context.Context, ru1 RoleUser) (RoleUser, error) {
	var ru RoleUser

	builder := utils.PSQL().Insert("people.role_members").
		Columns("role_id", "user_id").
		Values(ru1.RoleID, ru1.UserID).
		Suffix("RETURNING role_id, user_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for addRoleUser: %w", err))
	}

	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return RoleUser{}, fmt.Errorf("failed to add role user: %w", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(args...).Scan(&ru.RoleID, &ru.UserID)
	if err != nil {
		return RoleUser{}, fmt.Errorf("failed to add role user: %w", err)
	}

	return ru, nil
}

// RemoveRoleUser removes a link between a role.Role and User
func (s *Store) RemoveRoleUser(ctx context.Context, ru RoleUser) error {
	builder := utils.PSQL().Delete("people.role_members").
		Where(sq.And{
			sq.Eq{"role_id": ru.RoleID},
			sq.Eq{"user_id": ru.UserID},
		})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for removeRoleUser: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to remove role user: %w", err)
	}

	return nil
}

// RemoveUserForRoles removes all links between role.Role and a User
func (s *Store) RemoveUserForRoles(ctx context.Context, userID int) error {
	builder := utils.PSQL().Delete("people.role_members").
		Where(sq.Eq{"user_id": userID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for removeUserForRoles: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to remove user for roles: %w", err)
	}

	return nil
}

// GetPermissionsForRole returns all permissions for a role - moved here for cycle import reasons
func (s *Store) GetPermissionsForRole(ctx context.Context, roleID int) ([]Permission, error) {
	var p []Permission

	builder := utils.PSQL().Select("p.*").
		From("people.permissions p").
		LeftJoin("people.role_permissions rp ON p.permission_id = rp.permission_id").
		Where(sq.Eq{"rp.role_id": roleID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getPermissionsForRole: %w", err))
	}

	err = s.db.SelectContext(ctx, &p, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions for role: %w", err)
	}

	return p, nil
}

// GetRolesForPermission returns all roles for a permission - moved here for cycle import reasons
func (s *Store) GetRolesForPermission(ctx context.Context, permissionID int) ([]Role, error) {
	var r []Role

	builder := utils.PSQL().Select("r.*").
		From("people.roles r").
		LeftJoin("people.role_permissions rp ON r.role_id = rp.role_id").
		Where(sq.Eq{"rp.permission_id": permissionID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getRolesForPermission: %w", err))
	}

	err = s.db.SelectContext(ctx, &r, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for permission: %w", err)
	}

	return r, nil
}

// GetRolePermission returns a role permission - moved here for cycle import reasons
func (s *Store) GetRolePermission(ctx context.Context, rp1 RolePermission) (RolePermission, error) {
	var rp RolePermission

	builder := utils.PSQL().Select("*").
		From("people.role_permissions").
		Where(sq.And{
			sq.Eq{"role_id": rp1.RoleID},
			sq.Eq{"permission_id": rp1.PermissionID},
		}).
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getRolePermission: %w", err))
	}

	err = s.db.GetContext(ctx, &rp, sql, args...)
	if err != nil {
		return RolePermission{}, fmt.Errorf("failed to get role permission: %w", err)
	}

	return rp, nil
}

// GetPermissionsNotInRole returns all the permissions not currently in the role.Role to be added
func (s *Store) GetPermissionsNotInRole(ctx context.Context, roleID int) ([]Permission, error) {
	var p []Permission

	subQuery := utils.PSQL().Select("p.permission_id").
		From("people.permissions p").
		LeftJoin("people.role_permissions rp on p.permission_id = rp.permission_id").
		Where(sq.Eq{"rp.role_id": roleID})

	builder := utils.PSQL().Select("p.*").
		Distinct().
		From("people.permissions p").
		Where(utils.NotIn("permission_id", subQuery)).
		OrderBy("name")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getPermissionsNotInRole: %w", err))
	}

	//nolint:asasalint
	err = s.db.SelectContext(ctx, &p, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions not in role: %w", err)
	}

	return p, nil
}

// AddRolePermission creates a link between a role.Role and permission.Permission
func (s *Store) AddRolePermission(ctx context.Context, rp1 RolePermission) (RolePermission, error) {
	var rp RolePermission

	builder := utils.PSQL().Insert("people.role_permissions").
		Columns("role_id ", "permission_id").
		Values(rp1.RoleID, rp1.PermissionID).
		Suffix("RETURNING role_id, permission_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for addRolePermission: %w", err))
	}

	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return RolePermission{}, fmt.Errorf("failed to add rolePermission: %w", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(args...).Scan(&rp.RoleID, &rp.PermissionID)
	if err != nil {
		return RolePermission{}, fmt.Errorf("failed to add rolePermission: %w", err)
	}

	return rp, nil
}

// RemoveRolePermission removes a link between a role.Role and permission.Permission
func (s *Store) RemoveRolePermission(ctx context.Context, rp RolePermission) error {
	builder := utils.PSQL().Delete("people.role_permissions").
		Where(sq.And{sq.Eq{"role_id": rp.RoleID}, sq.Eq{"permission_id": rp.PermissionID}})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for removeRolePermission: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete rolePermission: %w", err)
	}

	return nil
}
