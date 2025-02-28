package people

import (
	"context"
	//nolint:gosec
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
)

var _ PermissionRepo = &Store{}

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

func (s *Store) ListPermissionMembersByID(ctx context.Context, permissionID int) ([]User, error) {
	var u []User
	//nolint:musttag
	err := s.db.GetContext(ctx, &u,
		`SELECT u.user_id, username, email, first_name, last_name, nickname, avatar
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
