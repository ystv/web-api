package clapper

import (
	"context"
	"time"

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
		IsPrivate   bool        `db:"is_private" json:"isPrivate"`
		IsCancelled bool        `db:"is_cancelled" json:"isCancelled"`
		IsTentative bool        `db:"is_tentative" json:"isTentative"`
		Signups     []Signup    `json:"signups,omitempty"`   // Used for shows
		Attendees   []Attendee  `json:"attendees,omitempty"` // Used for social, meet and other. This would be a XOR with Signups
	}
	// Signup represents a signup sheet which contains a group of roles
	Signup struct {
		SignupID    int            `db:"signup_id" json:"signupID"`
		Title       string         `db:"title" json:"title"`
		Description null.String    `db:"description" json:"description"`
		UnlockDate  null.Time      `db:"unlock_date" json:"unlockDate"`
		StartTime   null.Time      `db:"start_time" json:"startTime"`
		EndTime     null.Time      `db:"end_time" json:"endTime"`
		Crew        []CrewPosition `json:"crew"`
	}
	// CrewPosition represents a role for a signup sheet
	CrewPosition struct {
		CrewID   int `db:"crew_id" json:"crewID"`
		User     `json:"user"`
		Locked   bool `db:"locked" json:"locked"`
		Credited bool `db:"credited" json:"credited"`
		Ordering int  `db:"ordering" json:"ordering,omitempty"`
		Position
	}
	// Attendee represents a persons attendance for a meeting, social or other
	Attendee struct {
		User
		AttendStatus string `db:"attend_status" json:"attendStatus"`
	}
	// User a basic representation of a user
	User struct {
		UserID    int    `db:"user_id" json:"userID"`
		Nickname  string `db:"nickname" json:"nickname"`
		FirstName string `db:"first_name" json:"firstName"`
		LastName  string `db:"last_name" json:"lastName"`
	}
)

// Position is a role people can signup too
type Position struct {
	PositionID   int      `db:"position_id" json:"positionID"`
	Name         string   `db:"name" json:"name"`
	Description  string   `db:"description" json:"description"`
	Admin        bool     `db:"admin" json:"admin"`
	PermissionID null.Int `db:"permission_id" json:"permissionID"`
}

type (
	// EventRepo defines all event interactions
	EventRepo interface {
		ListMonth(ctx context.Context, year, month int) (*[]Event, error)
		Get(ctx context.Context, eventID int) (*Event, error)
		New(ctx context.Context, e *Event, userID int) (int, error)
		Update(ctx context.Context, e *Event, userID int) error
		Delete(ctx context.Context, eventID int) error
	}

	// PositionRepo defines all position interactions
	PositionRepo interface {
		List(ctx context.Context) (*[]Position, error)
		New(ctx context.Context, p *Position) (int, error)
		Update(ctx context.Context, p *Position) error
		Delete(ctx context.Context, positionID int) error
	}
)
