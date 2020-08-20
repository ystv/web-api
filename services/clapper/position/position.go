package position

import (
	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

// Position is a role people can signup too
type Position struct {
	PositionID   int         `db:"position_id" json:"positionID"`
	Name         string      `db:"name" json:"name"`
	Description  null.String `db:"description" json:"description"`
	Admin        bool        `db:"admin" json:"admin"`
	Credible     bool        `db:"credible" json:"credible"` // TODO move credible to crew
	PermissionID null.Int    `db:"permission_id" json:"permissionID"`
}

// List returns all positions
func List() (*[]Position, error) {
	p := []Position{}
	err := utils.DB.Select(&p,
		`SELECT position_id, name, description, admin, credible, permission_id
		FROM event.positions;`) // TODO fix permissions
	if err != nil {
		return nil, err
	}
	return &p, nil
}
