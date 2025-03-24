package public

import (
	"context"
	//nolint:gosec
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/guregu/null.v4"
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
		UserName           string     `json:"userName"`
		Avatar             string     `json:"avatar"`
		EmailAlias         string     `json:"emailAlias"`
		OfficerName        string     `json:"officerName"`
		OfficerDescription string     `json:"officerDescription"`
		HistoryWikiURL     string     `json:"historywikiURL"`
		StartDate          *time.Time `json:"startDate,omitempty"`
		EndDate            *time.Time `json:"endDate,omitempty"`
	}

	// TeamMemberDB a position within a group
	TeamMemberDB struct {
		UserID             int       `db:"user_id"`
		UserName           string    `db:"user_name"`
		UserEmail          string    `db:"user_email"`
		Avatar             string    `db:"avatar"`
		UseGravatar        bool      `db:"use_gravatar"`
		EmailAlias         string    `db:"email_alias"`
		OfficerName        string    `db:"officer_name"`
		OfficerDescription string    `db:"officer_description"`
		HistoryWikiURL     string    `db:"historywiki_url"`
		StartDate          null.Time `db:"start_date"`
		EndDate            null.Time `db:"end_date"`
	}
)

func (s *Store) TeamMemberDBToTeamMember(teamMember TeamMemberDB) TeamMember {
	var startDate, endDate *time.Time
	if teamMember.StartDate.Valid {
		startDate = &teamMember.StartDate.Time
	}
	if teamMember.EndDate.Valid {
		endDate = &teamMember.EndDate.Time
	}
	return TeamMember{
		UserName:           teamMember.UserName,
		Avatar:             teamMember.Avatar,
		EmailAlias:         teamMember.EmailAlias,
		OfficerName:        teamMember.OfficerName,
		OfficerDescription: teamMember.OfficerDescription,
		HistoryWikiURL:     teamMember.HistoryWikiURL,
		StartDate:          startDate,
		EndDate:            endDate,
	}
}

// ListTeams returns a list of the ystv teams and their current members.
func (s *Store) ListTeams(ctx context.Context) ([]Team, error) {
	var t []Team
	//nolint:musttag
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

	teamMembersDB, err := s.ListTeamMembers(ctx, t.TeamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by email: %w", err)
	}

	teamMembers := make([]TeamMember, 0)
	for _, m := range teamMembersDB {
		var startDate, endDate *time.Time
		if m.StartDate.Valid {
			startDate = &m.StartDate.Time
		}
		if m.EndDate.Valid {
			endDate = &m.EndDate.Time
		}

		teamMembers = append(teamMembers, TeamMember{
			UserName:           m.UserName,
			Avatar:             m.Avatar,
			EmailAlias:         m.EmailAlias,
			OfficerName:        m.OfficerName,
			OfficerDescription: m.OfficerDescription,
			HistoryWikiURL:     m.HistoryWikiURL,
			StartDate:          startDate,
			EndDate:            endDate,
		})
	}

	t.Members = teamMembers

	return t, nil
}

// GetTeamByID returns a single team including its members
func (s *Store) GetTeamByID(ctx context.Context, teamID int) (Team, error) {
	t, err := s.getTeamByID(ctx, teamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team by id: %w", err)
	}

	teamMembersDB, err := s.ListTeamMembers(ctx, t.TeamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by id: %w", err)
	}

	teamMembers := make([]TeamMember, 0)
	for _, m := range teamMembersDB {
		var startDate, endDate *time.Time
		if m.StartDate.Valid {
			startDate = &m.StartDate.Time
		}
		if m.EndDate.Valid {
			endDate = &m.EndDate.Time
		}

		teamMembers = append(teamMembers, TeamMember{
			UserName:           m.UserName,
			Avatar:             m.Avatar,
			EmailAlias:         m.EmailAlias,
			OfficerName:        m.OfficerName,
			OfficerDescription: m.OfficerDescription,
			HistoryWikiURL:     m.HistoryWikiURL,
			StartDate:          startDate,
			EndDate:            endDate,
		})
	}

	t.Members = teamMembers

	return t, nil
}

// GetTeamByYearByEmail returns a team by a calendar year
func (s *Store) GetTeamByYearByEmail(ctx context.Context, emailAlias string, year int) (Team, error) {
	t, err := s.getTeamByEmail(ctx, emailAlias)
	if err != nil {
		return t, fmt.Errorf("failed to get team by year by email: %w", err)
	}

	teamMembers := make([]TeamMemberDB, 0)

	err = s.db.SelectContext(ctx, &teamMembers, `
		SELECT CONCAT(users.first_name, ' ', users.last_name) AS user_name, users.email AS user_email, COALESCE(users.avatar, '') AS avatar,
		users.use_gravatar AS use_gravatar, officer.email_alias, officer.name AS officer_name, officer.description AS officer_description,
		officer.historywiki_url, officerTeamMembers.start_date, officerTeamMembers.end_date
		FROM people.officership_teams teams
		INNER JOIN people.officership_team_members teamMembers ON teams.team_id = teamMembers.team_id
		INNER JOIN people.officerships officer ON teamMembers.officer_id = officer.officer_id
		INNER JOIN people.officership_members officerTeamMembers ON officerTeamMembers.officer_id = teamMembers.officer_id
		INNER JOIN people.users users ON officerTeamMembers.user_id = users.user_id
		WHERE EXTRACT(year FROM officerTeamMembers.start_date) <= $1 AND (EXTRACT(year FROM officerTeamMembers.end_date) >= $1 OR officerTeamMembers.end_date IS NULL) AND
		teams.email_alias = $2
		ORDER BY officerTeamMembers.start_date, CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' AND officer.name NOT LIKE '%Assistant%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END;`, year, emailAlias)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by year by email: %w", err)
	}

	for _, teamMember := range teamMembers {
		switch avatar := teamMember.Avatar; {
		case teamMember.UseGravatar:
			//nolint:gosec
			hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(teamMember.UserEmail))))
			teamMember.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
		case avatar == "":
			teamMember.Avatar = "https://auth.ystv.co.uk/public/ystv-colour-white.png"
		case strings.Contains(avatar, s.cdnEndpoint):
		case strings.Contains(avatar, fmt.Sprintf("%d.", teamMember.UserID)):
			teamMember.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
		default:
			log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", teamMember.UserID, len(teamMember.Avatar), teamMember.Avatar)
			teamMember.Avatar = ""
		}
		t.Members = append(t.Members, s.TeamMemberDBToTeamMember(teamMember))
	}

	return t, nil
}

// GetTeamByYearByID returns a team by a calendar year
func (s *Store) GetTeamByYearByID(ctx context.Context, teamID, year int) (Team, error) {
	t, err := s.getTeamByID(ctx, teamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team by year by id: %w", err)
	}

	teamMembers := make([]TeamMemberDB, 0)

	err = s.db.SelectContext(ctx, &teamMembers, `
		SELECT CONCAT(users.first_name, ' ', users.last_name) AS user_name, COALESCE(users.avatar, '') AS avatar,
		officer.email_alias, officer.name AS officer_name, officer.description AS officer_description,
		officer.historywiki_url, officerTeamMembers.start_date, officerTeamMembers.end_date
		FROM people.officership_teams teams
		INNER JOIN people.officership_team_members teamMembers ON teams.team_id = teamMembers.team_id
		INNER JOIN people.officerships officer ON teamMembers.officer_id = officer.officer_id
		INNER JOIN people.officership_members officerTeamMembers ON officerTeamMembers.officer_id = teamMembers.officer_id
		INNER JOIN people.users users ON officerTeamMembers.user_id = users.user_id
		WHERE EXTRACT(year FROM officerTeamMembers.start_date) <= $1 AND (EXTRACT(year FROM officerTeamMembers.end_date) >= $1 OR officerTeamMembers.end_date IS NULL) AND
		teams.team_id = $2
		ORDER BY officerTeamMembers.start_date, CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' AND officer.name NOT LIKE '%Assistant%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END;`, year, teamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by year by id: %w", err)
	}

	for _, teamMember := range teamMembers {
		switch avatar := teamMember.Avatar; {
		case teamMember.UseGravatar:
			//nolint:gosec
			hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(teamMember.UserEmail))))
			teamMember.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
		case avatar == "":
			teamMember.Avatar = "https://auth.ystv.co.uk/public/ystv-colour-white.png"
		case strings.Contains(avatar, s.cdnEndpoint):
		case strings.Contains(avatar, fmt.Sprintf("%d.", teamMember.UserID)):
			teamMember.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
		default:
			log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", teamMember.UserID, len(teamMember.Avatar), teamMember.Avatar)
			teamMember.Avatar = ""
		}
		t.Members = append(t.Members, s.TeamMemberDBToTeamMember(teamMember))
	}

	return t, nil
}

// GetTeamByStartEndYearByEmail returns a team by an academic year
func (s *Store) GetTeamByStartEndYearByEmail(ctx context.Context, emailAlias string, startYear, endYear int) (Team, error) {
	t, err := s.getTeamByEmail(ctx, emailAlias)
	if err != nil {
		return t, fmt.Errorf("failed to get team by start end year by email: %w", err)
	}

	teamMembers := make([]TeamMemberDB, 0)

	err = s.db.SelectContext(ctx, &teamMembers, `
		SELECT CONCAT(users.first_name, ' ', users.last_name) AS user_name, COALESCE(users.avatar, '') AS avatar,
		officer.email_alias, officer.name AS officer_name, officer.description AS officer_description,
		officer.historywiki_url, officerTeamMembers.start_date, officerTeamMembers.end_date
		FROM people.officership_teams teams
		INNER JOIN people.officership_team_members teamMembers ON teams.team_id = teamMembers.team_id
		INNER JOIN people.officerships officer ON teamMembers.officer_id = officer.officer_id
		INNER JOIN people.officership_members officerTeamMembers ON officerTeamMembers.officer_id = teamMembers.officer_id
		INNER JOIN people.users users ON officerTeamMembers.user_id = users.user_id
		WHERE ((EXTRACT(year FROM officerTeamMembers.start_date) <= $1 AND EXTRACT(year FROM officerTeamMembers.end_date) >= $1) OR
		       (EXTRACT(year FROM officerTeamMembers.start_date) >= $1 AND EXTRACT(year FROM officerTeamMembers.end_date) <= $2) OR
		       (EXTRACT(year FROM officerTeamMembers.start_date) <= $2 AND (EXTRACT(year FROM officerTeamMembers.end_date) >= $2 OR officerTeamMembers.end_date IS NULL))) AND
		teams.email_alias = $3
		ORDER BY officerTeamMembers.start_date, CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' AND officer.name NOT LIKE '%Assistant%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END;`, startYear, endYear, emailAlias)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by start end year by email: %w", err)
	}

	for _, teamMember := range teamMembers {
		switch avatar := teamMember.Avatar; {
		case teamMember.UseGravatar:
			//nolint:gosec
			hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(teamMember.UserEmail))))
			teamMember.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
		case avatar == "":
			teamMember.Avatar = "https://auth.ystv.co.uk/public/ystv-colour-white.png"
		case strings.Contains(avatar, s.cdnEndpoint):
		case strings.Contains(avatar, fmt.Sprintf("%d.", teamMember.UserID)):
			teamMember.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
		default:
			log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", teamMember.UserID, len(teamMember.Avatar), teamMember.Avatar)
			teamMember.Avatar = ""
		}
		t.Members = append(t.Members, s.TeamMemberDBToTeamMember(teamMember))
	}

	return t, nil
}

// GetTeamByStartEndYearByID returns a team by an academic year
func (s *Store) GetTeamByStartEndYearByID(ctx context.Context, teamID, startYear, endYear int) (Team, error) {
	t, err := s.getTeamByID(ctx, teamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team by start end year by id: %w", err)
	}

	teamMembers := make([]TeamMemberDB, 0)

	err = s.db.SelectContext(ctx, &teamMembers, `
		SELECT CONCAT(users.first_name, ' ', users.last_name) AS user_name, COALESCE(users.avatar, '') AS avatar,
		officer.email_alias, officer.name AS officer_name, officer.description AS officer_description,
		officer.historywiki_url, officerTeamMembers.start_date, officerTeamMembers.end_date
		FROM people.officership_teams teams
		INNER JOIN people.officership_team_members teamMembers ON teams.team_id = teamMembers.team_id
		INNER JOIN people.officerships officer ON teamMembers.officer_id = officer.officer_id
		INNER JOIN people.officership_members officerTeamMembers ON officerTeamMembers.officer_id = teamMembers.officer_id
		INNER JOIN people.users users ON officerTeamMembers.user_id = users.user_id
		WHERE ((EXTRACT(year FROM officerTeamMembers.start_date) <= $1 AND EXTRACT(year FROM officerTeamMembers.end_date) >= $1) OR
		       (EXTRACT(year FROM officerTeamMembers.start_date) >= $1 AND EXTRACT(year FROM officerTeamMembers.end_date) <= $2) OR
		       (EXTRACT(year FROM officerTeamMembers.start_date) <= $2 AND (EXTRACT(year FROM officerTeamMembers.end_date) >= $2 OR officerTeamMembers.end_date IS NULL))) AND
		teams.team_id = $3
		ORDER BY officerTeamMembers.start_date, CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' AND officer.name NOT LIKE '%Assistant%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END;`, startYear, endYear, teamID)
	if err != nil {
		return t, fmt.Errorf("failed to get team members by start end year by id: %w", err)
	}

	for _, teamMember := range teamMembers {
		switch avatar := teamMember.Avatar; {
		case teamMember.UseGravatar:
			//nolint:gosec
			hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(teamMember.UserEmail))))
			teamMember.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
		case avatar == "":
			teamMember.Avatar = "https://auth.ystv.co.uk/public/ystv-colour-white.png"
		case strings.Contains(avatar, s.cdnEndpoint):
		case strings.Contains(avatar, fmt.Sprintf("%d.", teamMember.UserID)):
			teamMember.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
		default:
			log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", teamMember.UserID, len(teamMember.Avatar), teamMember.Avatar)
			teamMember.Avatar = ""
		}
		t.Members = append(t.Members, s.TeamMemberDBToTeamMember(teamMember))
	}

	return t, nil
}

func (s *Store) getTeamByEmail(ctx context.Context, emailAlias string) (Team, error) {
	var team Team
	//nolint:musttag
	err := s.db.GetContext(ctx, &team, `
		SELECT team_id, name, email_alias, short_description, full_description
		FROM people.officership_teams
		WHERE email_alias = $1;`, emailAlias)
	if err != nil {
		return team, err
	}

	return team, nil
}

func (s *Store) getTeamByID(ctx context.Context, teamID int) (Team, error) {
	var team Team
	//nolint:musttag
	err := s.db.GetContext(ctx, &team, `
		SELECT team_id, name, email_alias, short_description, full_description
		FROM people.officership_teams
		WHERE team_id = $1;`, teamID)
	if err != nil {
		return team, err
	}

	return team, nil
}

// ListTeamMembers returns a list of TeamMembers who are part of a team
func (s *Store) ListTeamMembers(ctx context.Context, teamID int) ([]TeamMemberDB, error) {
	m := make([]TeamMemberDB, 0)
	var temp []TeamMemberDB

	err := s.db.SelectContext(ctx, &temp, `
		SELECT CONCAT(first_name, ' ', last_name) AS user_name, COALESCE(avatar, '') AS avatar,
		email_alias, officer.name AS officer_name, officer.description AS officer_description,
		historywiki_url, off_mem.start_date, off_mem.end_date
		FROM people.officerships officer
		INNER JOIN people.officership_members off_mem ON officer.officer_id = off_mem.officer_id
		INNER JOIN people.users u ON off_mem.user_id = u.user_id
		INNER JOIN people.officership_team_members tm ON officer.officer_id = tm.officer_id
		WHERE off_mem.start_date < NOW() AND (off_mem.end_date > NOW() OR off_mem.end_date IS NULL) AND
		tm.team_id = $1
		ORDER BY CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' AND officer.name NOT LIKE '%Assistant%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END, off_mem.start_date;`, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}

	for _, teamMember := range temp {
		switch avatar := teamMember.Avatar; {
		case teamMember.UseGravatar:
			//nolint:gosec
			hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(teamMember.UserEmail))))
			teamMember.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
		case avatar == "":
			teamMember.Avatar = "https://auth.ystv.co.uk/public/ystv-colour-white.png"
		case strings.Contains(avatar, s.cdnEndpoint):
		case strings.Contains(avatar, fmt.Sprintf("%d.", teamMember.UserID)):
			teamMember.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
		default:
			log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", teamMember.UserID, len(teamMember.Avatar), teamMember.Avatar)
			teamMember.Avatar = ""
		}
		m = append(m, teamMember)
	}

	return m, nil
}

// ListOfficers returns the list of current officers
func (s *Store) ListOfficers(ctx context.Context) ([]TeamMemberDB, error) {
	m := make([]TeamMemberDB, 0)
	var temp []TeamMemberDB

	err := s.db.SelectContext(ctx, &temp, `
		SELECT CONCAT(first_name, ' ', last_name) AS user_name, COALESCE(avatar, '') AS avatar,
		email_alias, officer.name AS officer_name, officer.description AS officer_description,
		historywiki_url, off_mem.start_date, off_mem.end_date
		FROM people.officerships officer
		INNER JOIN people.officership_members off_mem ON officer.officer_id = off_mem.officer_id
		INNER JOIN people.users u ON off_mem.user_id = u.user_id
		WHERE off_mem.start_date < NOW() AND (off_mem.end_date > NOW() OR off_mem.end_date IS NULL)
		ORDER BY CASE
		    WHEN officer.name = 'Station Director' THEN 0
		    WHEN officer.name LIKE '%Director%' AND officer.name NOT LIKE '%Deputy%' AND officer.name NOT LIKE '%Assistant%' THEN 1
		    WHEN officer.name LIKE '%Deputy%' THEN 2
		    WHEN officer.name LIKE '%Assistant%' THEN 3
		    WHEN officer.name = 'Head of Welfare and Training' THEN 4
		    WHEN officer.name LIKE '%Head of%' THEN 5
		    ELSE 6 END, off_mem.start_date;`)
	if err != nil {
		return nil, fmt.Errorf("failed to list all officers: %w", err)
	}

	for _, teamMember := range temp {
		switch avatar := teamMember.Avatar; {
		case teamMember.UseGravatar:
			//nolint:gosec
			hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(teamMember.UserEmail))))
			teamMember.Avatar = "https://www.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
		case avatar == "":
			teamMember.Avatar = "https://auth.ystv.co.uk/public/ystv-colour-white.png"
		case strings.Contains(avatar, s.cdnEndpoint):
		case strings.Contains(avatar, fmt.Sprintf("%d.", teamMember.UserID)):
			teamMember.Avatar = "https://ystv.co.uk/static/images/members/thumb/" + avatar
		default:
			log.Printf("unknown avatar, user id: %d, length: %d, db string: %s, continuing", teamMember.UserID, len(teamMember.Avatar), teamMember.Avatar)
			teamMember.Avatar = ""
		}
		m = append(m, teamMember)
	}

	return m, nil
}
