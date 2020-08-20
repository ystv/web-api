package event

import (
	"time"

	"github.com/ystv/web-api/utils"
	"gopkg.in/guregu/null.v4"
)

type (
	// Event represents a group of signups
	Event struct {
		EventID     int         `db:"event_id" json:"eventID"`
		EventType   string      `db:"event_type" json:"eventType"`
		Name        string      `db:"name" json:"name"`
		StartDate   time.Time   `db:"start_date" json:"startDate"`
		EndDate     time.Time   `db:"end_date" json:"endDate"`
		Description null.String `db:"description" json:"description"`
		Location    null.String `db:"location" json:"location"`
		Status      string      `db:"status" json:"status"`
		Signups     []Signup    `json:"signups,omitempty"`
	}
	// Signup represents a signup sheet which contains a group of roles
	Signup struct {
		SignupID    int            `db:"signup_id" json:"signupID"`
		Title       string         `db:"title" json:"title"`
		Description null.String    `db:"description" json:"description"`
		UnlockDate  null.Time      `db:"unlock_date" json:"unlockDate"`
		StartTime   null.Int       `db:"start_time" json:"startDate"`
		EndTime     null.Int       `db:"end_time" json:"endDate"`
		Crew        []CrewPosition `json:"crew"`
	}
	// CrewPosition represents a role for a signup sheet
	CrewPosition struct {
		CrewID   int `db:"crew_id" json:"crewID"`
		User     `json:"user"`
		Locked   bool `db:"locked" json:"locked"`
		Ordering int  `db:"ordering" json:"ordering,omitempty"`
		Position
	}
	// Position is a role people can signup too
	Position struct {
		PositionID   int         `db:"position_id" json:"positionID"`
		Name         string      `db:"name" json:"name"`
		Description  null.String `db:"description" json:"description"`
		Admin        bool        `db:"admin" json:"admin"`
		Credible     bool        `db:"credible" json:"credible"`
		PermissionID null.Int    `db:"permission_id" json:"permissionID"`
	}
	// User a basic representation of a user
	User struct {
		UserID    int    `db:"member_id" json:"userID"`
		Nickname  string `db:"nickname" json:"nickname"`
		FirstName string `db:"first_name" json:"firstName"`
		LastName  string `db:"last_name" json:"lastName"`
	}
)

// ListMonth Lists all event meta's for a month
func ListMonth(year, month int) (*[]Event, error) {
	e := []Event{}
	err := utils.DB.Select(&e,
		`SELECT event_id, event_type, name, start_date, end_date, description, location, status
		FROM event.events
		WHERE EXTRACT(YEAR FROM start_date) = $1 AND
		EXTRACT(MONTH FROM start_date) = $2`, year, month)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// Get returns an event including the signup sheets
func Get(eventID int) (*Event, error) {
	e := Event{}
	err := utils.DB.Get(&e,
		`SELECT event_id, event_type, name, start_date, end_date, description, location, status
		FROM event.events
		WHERE event_id = $1;`, eventID)
	if err != nil {
		return nil, err
	}
	err = utils.DB.Select(&e.Signups,
		`SELECT signup_id, title, description, unlock_date, start_time, end_time
		FROM event.signups
		WHERE event_id = $1;`, eventID)
	if err != nil {
		return nil, err
	}
	for i := range e.Signups {
		err := utils.DB.Select(&e.Signups[i].Crew,
			`SELECT crew_id, member_id, nickname, first_name, last_name, locked, event.positions.position_id, name, description, admin, credible, permission_id
			FROM event.crew
			INNER JOIN event.positions ON event.positions.position_id = event.crew.position_id
			INNER JOIN people.users ON member_id = user_id
			WHERE signup_id = $1;`, e.Signups[i].SignupID) //TODO update crew to crews and sort out consistency
		if err != nil {
			return nil, err
		}
	}
	return &e, nil
}
