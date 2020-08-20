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
	PermissionID null.Int    `db:"permission_id" json:"permissionID"`
}

// List returns all positions
func List() (*[]Position, error) {
	p := []Position{}
	err := utils.DB.Select(&p,
		`SELECT position_id, name, description, admin, permission_id
		FROM event.positions;`)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// New creates a position
func New(p *Position) error {
	_, err := utils.DB.Exec(
		`INSERT INTO event.positions (name, description, admin, permission_id)
		VALUES ($1, $2, $3, $4);`, &p.Name, &p.Description, &p.Admin, &p.PermissionID)
	return err
}

// Update a position, uses the ID from the token
func Update(p *Position) error {
	_, err := utils.DB.Exec(
		`UPDATE event.positions
		SET name=$1, description=$2, admin=$3, permission_id=$4
		WHERE position_id = $5`,
		&p.Name, &p.Description, &p.Admin, &p.PermissionID, &p.PositionID)
	return err
}
