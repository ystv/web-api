package public

import (
	"context"
	"fmt"
)

type (
	// Team an organisational group
	Team struct {
		TeamID           int          `json:"id" db:"team_id"`
		Name             string       `json:"name" db:"name"`
		EmailAlias       string       `json:"emailAlias" db:"email_alias"`
		ShortDescription string       `json:"shortDescription" db:"short_description"`
		LongDescripton   string       `json:"longDescription,omitempty" db:"full_description"`
		Members          []TeamMember `json:"members,omitempty"`
	}
	// TeamMember a position within a group
	TeamMember struct {
		UserID             int    `json:"userID" db:"user_id"`
		UserName           string `json:"userName" db:"user_name"`
		Avatar             string `json:"avatar" db:"avatar"`
		OfficerID          int    `json:"officerID" db:"officer_id"`
		EmailAlias         string `json:"emailAlias" db:"email_alias"`
		OfficerName        string `json:"officerName" db:"officer_name"`
		OfficerDescription string `json:"officerDescription" db:"officer_description"`
		HistoryWikiURL     string `json:"historywikiURL" db:"historywiki_url"`
	}
)

var _ TeamRepo = &Store{}

// ListTeams returns a list of the ystv teams and their current members.
func (s *Store) ListTeams(ctx context.Context) ([]Team, error) {
	t := []Team{}
	err := s.db.SelectContext(ctx, &t, `
		SELECT team_id, name, email_alias, short_description, full_description
		FROM people.officership_teams
		ORDER BY name;`)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}
	return t, nil
}

// GetTeam returns a single team including it's members
func (s *Store) GetTeam(ctx context.Context, teamID int) (Team, error) {
	t := Team{}
	err := s.db.GetContext(ctx, &t, `
		SELECT team_id, name, email_alias, short_description, full_description
		FROM people.officership_teams
		WHERE team_id = $1;`, teamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team: %w", err)
	}
	t.Members, err = s.ListTeamMembers(ctx, teamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team members: %w", err)
	}
	return t, nil
}

// GetTeamByYear returns a team by a calendar year
func (s *Store) GetTeamByYear(ctx context.Context, teamID, year int) (Team, error) {
	t := Team{}
	err := s.db.GetContext(ctx, &t, `
		SELECT team_id, name, email_alias, short_description, full_description
		FROM people.officership_teams
		WHERE team_id = $1;`, teamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team by year: %w", err)
	}
	err = s.db.SelectContext(ctx, &t.Members, `
		SELECT u.user_id, CONCAT(first_name, ' ', last_name) AS user_name, COALESCE(avatar, '') AS avatar, officer.officer_id,
		email_alias, officer.name AS officer_name, officer.description AS officer_description,
		historywiki_url
		FROM people.officerships officer
		INNER JOIN people.officership_members off_mem ON officer.officer_id = off_mem.officer_id
		INNER JOIN people.users u ON off_mem.user_id = u.user_id
		INNER JOIN people.officership_team_members tm ON officer.officer_id = tm.officer_id
		WHERE EXTRACT(year FROM start_date) = $1 OR EXTRACT(year FROM end_date) = $1 AND
		team_id = $2;`, year, teamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by year: %w", err)
	}
	return t, nil
}

// ListTeamMembers returns a list of TeamMembers who are part of a team
func (s *Store) ListTeamMembers(ctx context.Context, teamID int) ([]TeamMember, error) {
	m := []TeamMember{}
	err := s.db.SelectContext(ctx, &m, `
		SELECT u.user_id, CONCAT(first_name, ' ', last_name) AS user_name, COALESCE(avatar, '') AS avatar, officer.officer_id,
		email_alias, officer.name AS officer_name, officer.description AS officer_description,
		historywiki_url
		FROM people.officerships officer
		INNER JOIN people.officership_members off_mem ON officer.officer_id = off_mem.officer_id
		INNER JOIN people.users u ON off_mem.user_id = u.user_id
		INNER JOIN people.officership_team_members tm ON officer.officer_id = tm.officer_id
		WHERE start_date < NOW() AND (end_date > NOW() OR end_date IS NULL) AND
		team_id = $1
		ORDER BY officer_id;`, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}
	return m, nil
}

// ListOfficers returns the list of officers of the current officers
func (s *Store) ListOfficers(ctx context.Context) ([]TeamMember, error) {
	m := []TeamMember{}
	err := s.db.SelectContext(ctx, &m, `
		SELECT u.user_id, CONCAT(first_name, ' ', last_name) AS user_name, COALESCE(avatar, '') AS avatar, officer.officer_id,
		email_alias, officer.name AS officer_name, officer.description AS officer_description,
		historywiki_url
		FROM people.officerships officer
		INNER JOIN people.officership_members off_mem ON officer.officer_id = off_mem.officer_id
		INNER JOIN people.users u ON off_mem.user_id = u.user_id
		WHERE start_date < NOW() AND (end_date > NOW() OR end_date IS NULL);`)
	if err != nil {
		return nil, fmt.Errorf("failed to list all officers: %w", err)
	}
	return m, nil
}
