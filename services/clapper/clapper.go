package clapper

import (
	"context"
	"time"
)

type (
	// New objects

	// NewEvent represents the necessary fields to create a new event.
	NewEvent struct {
		EventType   string    `db:"event_type" json:"eventType"`
		Name        string    `db:"name" json:"name"`
		StartDate   time.Time `db:"start_date" json:"startDate"`
		EndDate     time.Time `db:"end_date" json:"endDate"`
		Description string    `db:"description" json:"description"`
		Location    string    `db:"location" json:"location"`
		IsPrivate   bool      `db:"is_private" json:"isPrivate"`
		IsCancelled bool      `db:"is_cancelled" json:"isCancelled"`
		IsTentative bool      `db:"is_tentative" json:"isTentative"`
	}

	// NewSignup required fields to add
	NewSignup struct {
		EventID     int        `db:"event_id" json:"eventID"`
		Title       string     `db:"title" json:"title"`
		Description string     `db:"description" json:"description"`
		UnlockDate  *time.Time `db:"unlock_date" json:"unlockDate"`
		ArrivalTime *time.Time `db:"arrival_time" json:"arrivalTime"`
		StartTime   *time.Time `db:"start_time" json:"startTime"`
		EndTime     *time.Time `db:"end_time" json:"endTime"`
		Crew        []NewCrew  `json:"crew"`
	}
	// NewCrew required fields to add crew to a signup sheet
	NewCrew struct {
		PositionID int  `db:"position_id" json:"positionID"`
		Locked     bool `db:"locked" json:"locked"`
		Credited   bool `db:"credited" json:"credited"`
		Ordering   int  `db:"ordering" json:"ordering,omitempty"`
	}

	// Normal objects

	// Event represents a group of signups
	Event struct {
		EventID     int        `db:"event_id" json:"eventID"`
		EventType   string     `db:"event_type" json:"eventType"`
		Name        string     `db:"name" json:"name"`
		StartDate   time.Time  `db:"start_date" json:"startDate"`
		EndDate     time.Time  `db:"end_date" json:"endDate"`
		Description string     `db:"description" json:"description"`
		Location    string     `db:"location" json:"location"`
		IsPrivate   bool       `db:"is_private" json:"isPrivate"`
		IsCancelled bool       `db:"is_cancelled" json:"isCancelled"`
		IsTentative bool       `db:"is_tentative" json:"isTentative"`
		Signups     []Signup   `json:"signups,omitempty"`   // Used for shows
		Attendees   []Attendee `json:"attendees,omitempty"` // Used for social, meet and other. This would be an XOR with Signups
	}
	// Signup represents a signup sheet which contains a group of roles
	Signup struct {
		SignupID    int            `db:"signup_id" json:"signupID"`
		Title       string         `db:"title" json:"title"`
		Description string         `db:"description" json:"description"`
		UnlockDate  *time.Time     `db:"unlock_date" json:"unlockDate"`
		ArrivalTime *time.Time     `db:"arrival_time" json:"arrivalTime"`
		StartTime   *time.Time     `db:"start_time" json:"startTime"`
		EndTime     *time.Time     `db:"end_time" json:"endTime"`
		Crew        []CrewPosition `json:"crew"`
	}

	// Position is a role people can sign up too
	Position struct {
		PositionID   int    `db:"position_id" json:"positionID"`
		Name         string `db:"name" json:"name"`
		Description  string `db:"description" json:"description"`
		Admin        bool   `db:"admin" json:"admin"`
		PermissionID *int   `db:"permission_id" json:"permissionID"`
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
	// Attendee represents a person's attendance for a meeting, social or other
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

type (
	// EventRepo defines all event interactions
	EventRepo interface {
		ListMonth(ctx context.Context, year, month int) (*[]Event, error)
		Get(ctx context.Context, eventID int) (*Event, error)
		New(ctx context.Context, e *NewEvent, userID int) (int, error)
		Update(ctx context.Context, e *Event, userID int) error
		Delete(ctx context.Context, eventID int) error
	}

	// SignupRepo defines all signup sheet interactions
	SignupRepo interface {
		New(ctx context.Context, eventID int, s NewSignup) (int, error)
		Update(ctx context.Context, s Signup) error
		Delete(ctx context.Context, signupID int) error
	}

	// PositionRepo defines all position interactions.
	//
	// This repo is for managing the positions, where it
	// feeds the system where the producer makes the signup sheet,
	// but not the doesn't interact with an event directly.
	PositionRepo interface {
		List(ctx context.Context) (*[]Position, error)
		New(ctx context.Context, p *Position) (int, error)
		Update(ctx context.Context, p *Position) error
		Delete(ctx context.Context, positionID int) error
	}

	// CrewRepo defines all crew interactions.
	//
	// This repo is similar to the position one except it deals
	// with each unique role on a signup sheet. Providing the
	// facilities for users to use the signup sheet.
	CrewRepo interface {
		New(ctx context.Context, signupID, positionID int) error
		Get(ctx context.Context, crewID int) (*CrewPosition, error)
		UpdateUser(ctx context.Context, crewID, userID int) error
		UpdateUserAndVerify(ctx context.Context, crewID, userID int) error
		DeleteUser(ctx context.Context, crewID int) error
		Delete(ctx context.Context, crewID int) error
	}
)
