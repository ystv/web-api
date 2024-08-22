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
		LongDescription  string       `json:"longDescription,omitempty" db:"full_description"`
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

//var _ TeamRepo = &Store{}

// ListTeams returns a list of the ystv teams and their current members.
func (s *Store) ListTeams(ctx context.Context) ([]Team, error) {
	var t []Team
	err := s.db.SelectContext(ctx, &t, `
		SELECT team_id, name, email_alias, short_description, full_description
		FROM people.officership_teams
		ORDER BY name;`)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}
	return t, nil
}

// GetTeamByEmail returns a single team including its members
func (s *Store) GetTeamByEmail(ctx context.Context, emailAlias string) (Team, error) {
	t, err := s.getTeamByEmail(ctx, emailAlias)
	if err != nil {
		return t, fmt.Errorf("failed to get team by email: %w", err)
	}
	t.Members, err = s.ListTeamMembers(ctx, t.TeamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by email: %w", err)
	}
	return t, nil
}

// GetTeamById returns a single team including its members
func (s *Store) GetTeamById(ctx context.Context, teamId int) (Team, error) {
	t, err := s.getTeamById(ctx, teamId)
	if err != nil {
		return t, fmt.Errorf("failed to get team by id: %w", err)
	}
	t.Members, err = s.ListTeamMembers(ctx, t.TeamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by id: %w", err)
	}
	return t, nil
}

// GetTeamByYearByEmail returns a team by a calendar year
func (s *Store) GetTeamByYearByEmail(ctx context.Context, emailAlias string, year int) (Team, error) {
	t, err := s.getTeamByEmail(ctx, emailAlias)
	if err != nil {
		return t, fmt.Errorf("failed to get team by year by email: %w", err)
	}
	err = s.db.SelectContext(ctx, &t.Members, `
		SELECT users.user_id, CONCAT(users.first_name, ' ', users.last_name) AS user_name, COALESCE(users.avatar, '') AS avatar, officer.officer_id,
		officer.email_alias, officer.name AS officer_name, officer.description AS officer_description,
		officer.historywiki_url
		FROM people.officership_teams teams
		INNER JOIN people.officership_team_members teamMembers ON teams.team_id = teamMembers.team_id
		INNER JOIN people.officerships officer ON teamMembers.officer_id = officer.officer_id
		INNER JOIN people.officership_members officerTeamMembers ON officerTeamMembers.officer_id = teamMembers.officer_id
		INNER JOIN people.users users ON officerTeamMembers.user_id = users.user_id
		WHERE EXTRACT(year FROM officerTeamMembers.start_date) <= $1 AND (EXTRACT(year FROM officerTeamMembers.end_date) >= $1 OR officerTeamMembers.end_date IS NULL) AND
		teams.email_alias = $2
		ORDER BY start_date, CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END;`, year, emailAlias)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by year by email: %w", err)
	}
	return t, nil
}

// GetTeamByYearById returns a team by a calendar year
func (s *Store) GetTeamByYearById(ctx context.Context, teamId, year int) (Team, error) {
	t, err := s.getTeamById(ctx, teamId)
	if err != nil {
		return t, fmt.Errorf("failed to get team by year by id: %w", err)
	}
	err = s.db.SelectContext(ctx, &t.Members, `
		SELECT users.user_id, CONCAT(users.first_name, ' ', users.last_name) AS user_name, COALESCE(users.avatar, '') AS avatar, officer.officer_id,
		officer.email_alias, officer.name AS officer_name, officer.description AS officer_description,
		officer.historywiki_url
		FROM people.officership_teams teams
		INNER JOIN people.officership_team_members teamMembers ON teams.team_id = teamMembers.team_id
		INNER JOIN people.officerships officer ON teamMembers.officer_id = officer.officer_id
		INNER JOIN people.officership_members officerTeamMembers ON officerTeamMembers.officer_id = teamMembers.officer_id
		INNER JOIN people.users users ON officerTeamMembers.user_id = users.user_id
		WHERE (EXTRACT(year FROM officerTeamMembers.start_date) = $1 OR EXTRACT(year FROM officerTeamMembers.end_date) = $1) AND
		teams.team_id = $2
		ORDER BY CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END;`, year, teamId)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by year by id: %w", err)
	}
	return t, nil
}

// GetTeamByStartEndYearByEmail returns a team by an academic year
func (s *Store) GetTeamByStartEndYearByEmail(ctx context.Context, emailAlias string, startYear, endYear int) (Team, error) {
	t, err := s.getTeamByEmail(ctx, emailAlias)
	if err != nil {
		return t, fmt.Errorf("failed to get team by start end year by email: %w", err)
	}
	err = s.db.SelectContext(ctx, &t.Members, `
		SELECT users.user_id, CONCAT(users.first_name, ' ', users.last_name) AS user_name, COALESCE(users.avatar, '') AS avatar, officer.officer_id,
		officer.email_alias, officer.name AS officer_name, officer.description AS officer_description,
		officer.historywiki_url
		FROM people.officership_teams teams
		INNER JOIN people.officership_team_members teamMembers ON teams.team_id = teamMembers.team_id
		INNER JOIN people.officerships officer ON teamMembers.officer_id = officer.officer_id
		INNER JOIN people.officership_members officerTeamMembers ON officerTeamMembers.officer_id = teamMembers.officer_id
		INNER JOIN people.users users ON officerTeamMembers.user_id = users.user_id
		WHERE (EXTRACT(year FROM officerTeamMembers.start_date) = $1 OR EXTRACT(year FROM officerTeamMembers.end_date) = $2) AND
		teams.email_alias = $3
		ORDER BY CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END;`, startYear, endYear, emailAlias)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by start end year by email: %w", err)
	}
	return t, nil
}

// GetTeamByStartEndYearById returns a team by an academic year
func (s *Store) GetTeamByStartEndYearById(ctx context.Context, teamId, startYear, endYear int) (Team, error) {
	t, err := s.getTeamById(ctx, teamId)
	if err != nil {
		return t, fmt.Errorf("failed to get team by start end year by id: %w", err)
	}
	err = s.db.SelectContext(ctx, &t.Members, `
		SELECT users.user_id, CONCAT(users.first_name, ' ', users.last_name) AS user_name, COALESCE(users.avatar, '') AS avatar, officer.officer_id,
		officer.email_alias, officer.name AS officer_name, officer.description AS officer_description,
		officer.historywiki_url
		FROM people.officership_teams teams
		INNER JOIN people.officership_team_members teamMembers ON teams.team_id = teamMembers.team_id
		INNER JOIN people.officerships officer ON teamMembers.officer_id = officer.officer_id
		INNER JOIN people.officership_members officerTeamMembers ON officerTeamMembers.officer_id = teamMembers.officer_id
		INNER JOIN people.users users ON officerTeamMembers.user_id = users.user_id
		WHERE (EXTRACT(year FROM officerTeamMembers.start_date) = $1 OR EXTRACT(year FROM officerTeamMembers.end_date) = $2) AND
		teams.team_id = $3
		ORDER BY CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END;`, startYear, endYear, teamId)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by start end year by id: %w", err)
	}
	return t, nil
}

func (s *Store) getTeamByEmail(ctx context.Context, emailAlias string) (Team, error) {
	team := Team{}
	err := s.db.GetContext(ctx, &team, `
		SELECT team_id, name, email_alias, short_description, full_description
		FROM people.officership_teams
		WHERE email_alias = $1;`, emailAlias)
	if err != nil {
		return team, err
	}
	return team, nil
}

func (s *Store) getTeamById(ctx context.Context, teamId int) (Team, error) {
	team := Team{}
	err := s.db.GetContext(ctx, &team, `
		SELECT team_id, name, email_alias, short_description, full_description
		FROM people.officership_teams
		WHERE team_id = $1;`, teamId)
	if err != nil {
		return team, err
	}
	return team, nil
}

// ListTeamMembers returns a list of TeamMembers who are part of a team
func (s *Store) ListTeamMembers(ctx context.Context, teamID int) ([]TeamMember, error) {
	var m []TeamMember
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
		ORDER BY CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END;`, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}
	return m, nil
}

// ListOfficers returns the list of current officers
func (s *Store) ListOfficers(ctx context.Context) ([]TeamMember, error) {
	var m []TeamMember
	err := s.db.SelectContext(ctx, &m, `
		SELECT u.user_id, CONCAT(first_name, ' ', last_name) AS user_name, COALESCE(avatar, '') AS avatar, officer.officer_id,
		email_alias, officer.name AS officer_name, officer.description AS officer_description,
		historywiki_url
		FROM people.officerships officer
		INNER JOIN people.officership_members off_mem ON officer.officer_id = off_mem.officer_id
		INNER JOIN people.users u ON off_mem.user_id = u.user_id
		WHERE start_date < NOW() AND (end_date > NOW() OR end_date IS NULL)
		ORDER BY CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' AND officer.name NOT LIKE '%Assistant%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END, 
		    officer.name,
		    off_mem.start_date;`)
	if err != nil {
		return nil, fmt.Errorf("failed to list all officers: %w", err)
	}
	return m, nil
}
