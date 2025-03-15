package people

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/ystv/web-api/utils"
)

func (s *Store) OfficershipDBToOfficership(officershipDB OfficershipDB) Officership {
	var roleID, teamID *int64
	var ifUnfilled, isTeamLeader, isTeamDeputy *bool
	var teamName *string

	if officershipDB.RoleID.Valid {
		roleID = &officershipDB.RoleID.Int64
	}

	if officershipDB.TeamID.Valid {
		teamID = &officershipDB.TeamID.Int64
	}

	if officershipDB.IfUnfilled.Valid {
		ifUnfilled = &officershipDB.IfUnfilled.Bool
	}

	if officershipDB.IsTeamLeader.Valid {
		isTeamLeader = &officershipDB.IsTeamLeader.Bool
	}

	if officershipDB.IsTeamDeputy.Valid {
		isTeamDeputy = &officershipDB.IsTeamDeputy.Bool
	}

	if officershipDB.TeamName.Valid {
		teamName = &officershipDB.TeamName.String
	}

	return Officership{
		OfficershipID:    officershipDB.OfficershipID,
		Name:             officershipDB.Name,
		EmailAlias:       officershipDB.EmailAlias,
		Description:      officershipDB.Description,
		HistoryWikiURL:   officershipDB.HistoryWikiURL,
		RoleID:           roleID,
		IsCurrent:        officershipDB.IsCurrent,
		IfUnfilled:       ifUnfilled,
		CurrentOfficers:  officershipDB.CurrentOfficers,
		PreviousOfficers: officershipDB.PreviousOfficers,
		TeamID:           teamID,
		TeamName:         teamName,
		IsTeamLeader:     isTeamLeader,
		IsTeamDeputy:     isTeamDeputy,
	}
}

func (s *Store) OfficershipMemberDBToOfficershipMember(officershipDB OfficershipMemberDB) OfficershipMember {
	var startDate, endDate *time.Time
	var teamID *int
	var teamName *string

	if officershipDB.StartDate.Valid {
		startDate = &officershipDB.StartDate.Time
	}

	if officershipDB.EndDate.Valid {
		endDate = &officershipDB.EndDate.Time
	}

	if officershipDB.TeamID.Valid {
		temp := int(officershipDB.TeamID.Int64)
		teamID = &temp
	}

	if officershipDB.TeamName.Valid {
		teamName = &officershipDB.TeamName.String
	}

	return OfficershipMember{
		OfficershipMemberID: officershipDB.OfficershipMemberID,
		UserID:              officershipDB.UserID,
		OfficerID:           officershipDB.OfficerID,
		StartDate:           startDate,
		EndDate:             endDate,
		OfficershipName:     officershipDB.OfficershipName,
		UserName:            officershipDB.UserName,
		TeamID:              teamID,
		TeamName:            teamName,
	}
}

func (s *Store) CountOfficerships(ctx context.Context) (CountOfficerships, error) {
	var countOfficerships CountOfficerships

	err := s.db.GetContext(ctx, &countOfficerships,
		`SELECT
		(SELECT COUNT(*) FROM people.officerships) as total_officerships,
		(SELECT COUNT(*) FROM people.officerships WHERE is_current = true) as current_officerships,
		(SELECT COUNT(*) FROM people.officership_members) as total_officers,
		(SELECT COUNT(*) FROM people.officership_members WHERE end_date IS NULL) as current_officers;`)
	if err != nil {
		return countOfficerships, fmt.Errorf("failed to count officerships all from db: %w", err)
	}

	return countOfficerships, nil
}

func (s *Store) GetOfficerships(ctx context.Context, officershipStatus OfficershipsStatus) ([]OfficershipDB, error) {
	var o []OfficershipDB

	builder := utils.PSQL().Select("o.*", "COUNT(DISTINCT omc.officership_member_id) AS current_officers",
		"COUNT(DISTINCT omp.officership_member_id) AS previous_officers", "otm.team_id AS team_id",
		"ot.name AS team_name").
		From("people.officerships o").
		LeftJoin("people.officership_members omc ON o.officer_id = omc.officer_id AND omc.end_date IS NULL").
		LeftJoin("people.officership_members omp ON o.officer_id = omp.officer_id AND omp.end_date IS NOT NULL").
		LeftJoin("people.officership_team_members otm ON o.officer_id = otm.officer_id").
		LeftJoin("people.officership_teams ot ON ot.team_id = otm.team_id").
		GroupBy("o", "o.officer_id", "o.name", "o.email_alias", "description", "historywiki_url", "role_id",
			"is_current", "if_unfilled", "otm.team_id", "ot.name")

	switch officershipStatus {
	case Any:
	case Current:
		builder = builder.Where("o.is_current = true")
	case Retired:
		builder = builder.Where("o.is_current = false")
	}

	builder = builder.GroupBy("o", "o.officer_id", "o.name", "o.email_alias", "description",
		"historywiki_url", "role_id", "is_current", "if_unfilled").
		OrderBy(`CASE WHEN o.name = 'Station Director' THEN 0
	WHEN o.name LIKE '%Director%' AND o.name NOT LIKE '%Deputy%' AND o.name NOT LIKE '%Assistant%' THEN 1
	WHEN o.name LIKE '%Deputy%' THEN 2
	WHEN o.name LIKE '%Assistant%' THEN 3
	WHEN o.name = 'Head of Welfare and Training' THEN 4
	WHEN o.name LIKE '%Head of%' THEN 5
	ELSE 6 END`, "o.name")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getOfficerships: %w", err))
	}

	err = s.db.SelectContext(ctx, &o, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get officerships: %w", err)
	}

	return o, nil
}

func (s *Store) GetOfficership(ctx context.Context, officershipGetDTO OfficershipGetDTO) (OfficershipDB, error) {
	var o OfficershipDB

	builder := utils.PSQL().Select("o.*", "COUNT(DISTINCT omc.officership_member_id) AS current_officers",
		"COUNT(DISTINCT omp.officership_member_id) AS previous_officers", "otm.team_id AS team_id",
		"ot.name AS team_name", "otm.is_leader AS is_team_leader", "otm.is_deputy AS is_team_deputy").
		From("people.officerships o").
		LeftJoin("people.officership_members omc ON o.officer_id = omc.officer_id AND omc.end_date IS NULL").
		LeftJoin("people.officership_members omp ON o.officer_id = omp.officer_id AND omp.end_date IS NOT NULL").
		LeftJoin("people.officership_team_members otm ON o.officer_id = otm.officer_id").
		LeftJoin("people.officership_teams ot ON ot.team_id = otm.team_id").
		Where(sq.Or{
			sq.Eq{"o.officer_id": officershipGetDTO.OfficershipID},
			sq.And{
				sq.Eq{"o.name": officershipGetDTO.Name},
				sq.NotEq{"o.name": ""},
			},
		}).
		GroupBy("o", "o.officer_id", "o.name", "o.email_alias", "description", "historywiki_url", "role_id",
			"is_current", "if_unfilled", "otm.team_id", "ot.name", "otm.is_leader", "otm.is_deputy").
		OrderBy(`CASE WHEN o.name = 'Station Director' THEN 0
	WHEN o.name LIKE '%Director%' AND o.name NOT LIKE '%Deputy%' AND o.name NOT LIKE '%Assistant%' THEN 1
	WHEN o.name LIKE '%Deputy%' THEN 2
	WHEN o.name LIKE '%Assistant%' THEN 3
	WHEN o.name = 'Head of Welfare and Training' THEN 4
	WHEN o.name LIKE '%Head of%' THEN 5
	ELSE 6 END`, "o.name").
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getOfficership: %w", err))
	}

	err = s.db.GetContext(ctx, &o, sql, args...)
	if err != nil {
		return OfficershipDB{}, fmt.Errorf("failed to get officership: %w", err)
	}

	return o, nil
}

func (s *Store) AddOfficership(ctx context.Context, officershipAdd OfficershipAddEditDTO) (OfficershipDB, error) {
	builder := utils.PSQL().Insert("people.officerships").
		Columns("name", "email_alias", "description", "historywiki_url", "role_id", "is_current",
			"if_unfilled").
		Values(officershipAdd.Name, officershipAdd.EmailAlias, officershipAdd.Description, officershipAdd.HistoryWikiURL, nil, officershipAdd.IsCurrent, nil).
		Suffix("RETURNING officer_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for addOfficership: %w", err))
	}

	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return OfficershipDB{}, fmt.Errorf("failed to add officership: %w", err)
	}

	defer stmt.Close()

	var officershipID int

	err = stmt.QueryRow(args...).Scan(&officershipID)
	if err != nil {
		return OfficershipDB{}, fmt.Errorf("failed to add officership: %w", err)
	}

	return s.GetOfficership(ctx, OfficershipGetDTO{OfficershipID: officershipID})
}

func (s *Store) EditOfficership(ctx context.Context, officershipID int, officershipEdit OfficershipAddEditDTO) (OfficershipDB, error) {
	builder := utils.PSQL().Update("people.officerships").
		SetMap(map[string]interface{}{
			"name":            officershipEdit.Name,
			"email_alias":     officershipEdit.EmailAlias,
			"description":     officershipEdit.Description,
			"historywiki_url": officershipEdit.HistoryWikiURL,
			"role_id":         nil,
			"is_current":      officershipEdit.IsCurrent,
			"if_unfilled":     nil,
		}).
		Where(sq.Eq{"officer_id": officershipID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for editOfficership: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return OfficershipDB{}, fmt.Errorf("failed to edit officership: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return OfficershipDB{}, fmt.Errorf("failed to edit officership: %w", err)
	}

	if rows < 1 {
		return OfficershipDB{}, fmt.Errorf("failed to edit officerhip: invalid rows affected: %d", rows)
	}

	return s.GetOfficership(ctx, OfficershipGetDTO{OfficershipID: officershipID})
}

func (s *Store) DeleteOfficership(ctx context.Context, officershipID int) error {
	builder := utils.PSQL().Delete("people.officerships").
		Where(sq.Eq{"officer_id": officershipID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for deleteOfficership: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete officership: %w", err)
	}

	return nil
}

func (s *Store) GetOfficershipTeams(ctx context.Context) ([]OfficershipTeam, error) {
	var t []OfficershipTeam

	builder := utils.PSQL().Select("ot.*", "COUNT(DISTINCT otm.officer_id) AS current_officerships",
		"COUNT(DISTINCT om.officership_member_id) AS current_officers").
		From("people.officership_teams ot").
		LeftJoin("people.officership_team_members otm ON ot.team_id = otm.team_id").
		LeftJoin("people.officerships o ON otm.officer_id = o.officer_id").
		LeftJoin("people.officership_members om ON o.officer_id = om.officer_id AND om.end_date IS NULL AND o.is_current = true").
		GroupBy("ot", "ot.team_id", "ot.name", "ot.email_alias", "short_description", "full_description")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getOfficershipTeams: %w", err))
	}

	err = s.db.SelectContext(ctx, &t, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get officership teams: %w", err)
	}

	return t, nil
}

func (s *Store) GetOfficershipTeam(ctx context.Context, officershipTeamGet OfficershipTeamGetDTO) (OfficershipTeam, error) {
	var t OfficershipTeam

	builder := utils.PSQL().Select("ot.*", "COUNT(DISTINCT otm.officer_id) AS current_officerships",
		"COUNT(DISTINCT om.officership_member_id) AS current_officers").
		From("people.officership_teams ot").
		LeftJoin("people.officership_team_members otm ON ot.team_id = otm.team_id").
		LeftJoin("people.officerships o ON otm.officer_id = o.officer_id").
		LeftJoin("people.officership_members om ON o.officer_id = om.officer_id AND om.end_date IS NULL AND o.is_current = true").
		Where(sq.Or{sq.Eq{"ot.team_id": officershipTeamGet.TeamID}, sq.And{sq.Eq{"ot.name": officershipTeamGet.Name}, sq.NotEq{"ot.name": ""}}}).
		GroupBy("ot", "ot.team_id", "ot.name", "ot.email_alias", "short_description", "full_description").
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getOfficershipTeam: %w", err))
	}

	err = s.db.GetContext(ctx, &t, sql, args...)
	if err != nil {
		return OfficershipTeam{}, fmt.Errorf("failed to get officership team: %w", err)
	}

	return t, nil
}

func (s *Store) AddOfficershipTeam(ctx context.Context, officershipTeamAdd OfficershipTeamAddEditDTO) (OfficershipTeam, error) {
	builder := utils.PSQL().Insert("people.officership_teams").
		Columns("name", "email_alias", "short_description", "full_description").
		Values(officershipTeamAdd.Name, officershipTeamAdd.EmailAlias, officershipTeamAdd.ShortDescription, officershipTeamAdd.FullDescription).
		Suffix("RETURNING team_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for addOfficershipTeam: %w", err))
	}

	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return OfficershipTeam{}, fmt.Errorf("failed to add officership team: %w", err)
	}

	defer stmt.Close()

	var teamID int

	err = stmt.QueryRow(args...).Scan(&teamID)
	if err != nil {
		return OfficershipTeam{}, fmt.Errorf("failed to add offciership team: %w", err)
	}

	return s.GetOfficershipTeam(ctx, OfficershipTeamGetDTO{TeamID: teamID})
}

func (s *Store) EditOfficershipTeam(ctx context.Context, officershipTeamID int, officershipTeamAdd OfficershipTeamAddEditDTO) (OfficershipTeam, error) {
	builder := utils.PSQL().Update("people.officership_teams").
		SetMap(map[string]interface{}{
			"name":              officershipTeamAdd.Name,
			"email_alias":       officershipTeamAdd.EmailAlias,
			"short_description": officershipTeamAdd.ShortDescription,
			"full_description":  officershipTeamAdd.FullDescription,
		}).
		Where(sq.Eq{"team_id": officershipTeamID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for editOfficershipTeam: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return OfficershipTeam{}, fmt.Errorf("failed to edit officership team: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return OfficershipTeam{}, fmt.Errorf("failed to edit officership team: %w", err)
	}

	if rows < 1 {
		return OfficershipTeam{}, fmt.Errorf("failed to edit officerhip team: invalid rows affected: %d", rows)
	}

	return s.GetOfficershipTeam(ctx, OfficershipTeamGetDTO{TeamID: officershipTeamID})
}

func (s *Store) DeleteOfficershipTeam(ctx context.Context, officershipTeamID int) error {
	builder := utils.PSQL().Delete("people.officership_teams").
		Where(sq.Eq{"team_id": officershipTeamID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for deleteOfficershipTeam: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete officership team: %w", err)
	}

	return nil
}

func (s *Store) GetOfficershipTeamMembers(ctx context.Context, officershipTeamID *int,
	officershipStatus OfficershipsStatus) ([]OfficershipTeamMember, error) {
	var m []OfficershipTeamMember

	builder := utils.PSQL().Select("otm.*", "o.name AS officer_name",
		"COUNT(DISTINCT omc.officership_member_id) AS current_officers",
		"COUNT(DISTINCT omp.officership_member_id) AS previous_officers", "o.is_current AS is_current").
		From("people.officership_team_members otm").
		LeftJoin("people.officerships o on o.officer_id = otm.officer_id").
		LeftJoin("people.officership_members omc ON o.officer_id = omc.officer_id AND omc.end_date IS NULL").
		LeftJoin("people.officership_members omp ON o.officer_id = omp.officer_id AND omp.end_date IS NOT NULL")

	if officershipTeamID != nil {
		builder = builder.Where(sq.Eq{"otm.team_id": officershipTeamID})
	}

	switch officershipStatus {
	case Any:
	case Current:
		builder = builder.Where("o.is_current = true")
	case Retired:
		builder = builder.Where("o.is_current = false")
	}

	builder = builder.OrderBy(`CASE WHEN o.name = 'Station Director' THEN 0
	WHEN o.name LIKE '%Director%' AND o.name NOT LIKE '%Deputy%' AND o.name NOT LIKE '%Assistant%' THEN 1
	WHEN o.name LIKE '%Deputy%' THEN 2
	WHEN o.name LIKE '%Assistant%' THEN 3
	WHEN o.name = 'Head of Welfare and Training' THEN 4
	WHEN o.name LIKE '%Head of%' THEN 5
	ELSE 6 END`,
		"o.name").
		GroupBy("otm", "otm.officer_id", "otm.team_id", "o.officer_id", "name", "email_alias", "description",
			"historywiki_url", "role_id", "is_current", "if_unfilled")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getOfficershipTeamMembers: %w", err))
	}

	err = s.db.SelectContext(ctx, &m, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get officership team members: %w", err)
	}

	return m, nil
}

func (s *Store) GetOfficershipsNotInTeam(ctx context.Context, officershipTeamID int) ([]OfficershipDB, error) {
	var o []OfficershipDB

	subQuery := utils.PSQL().Select("o.officer_id").
		From("people.officerships o").
		LeftJoin("people.officership_team_members otm on o.officer_id = otm.officer_id").
		Where(sq.Eq{"otm.team_id": officershipTeamID})

	builder := utils.PSQL().Select("o.*").
		From("people.officerships o").
		Where(sq.And{
			utils.NotIn("o.officer_id", subQuery),
			utils.StringSQL("o.is_current = true"),
		}).
		OrderBy(`CASE WHEN o.name = 'Station Director' THEN 0
	WHEN o.name LIKE '%Director%' AND o.name NOT LIKE '%Deputy%' AND o.name NOT LIKE '%Assistant%' THEN 1
	WHEN o.name LIKE '%Deputy%' THEN 2
	WHEN o.name LIKE '%Assistant%' THEN 3
	WHEN o.name = 'Head of Welfare and Training' THEN 4
	WHEN o.name LIKE '%Head of%' THEN 5
	ELSE 6 END`, "o.name")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getOfficershipsNotInTeam: %w", err))
	}

	err = s.db.SelectContext(ctx, &o, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get officerships not in team: %w", err)
	}

	return o, nil
}

func (s *Store) GetOfficershipTeamMember(ctx context.Context, officershipTeamMemberGet OfficershipTeamMemberGetDeleteDTO) (OfficershipTeamMember, error) {
	var m OfficershipTeamMember

	builder := utils.PSQL().Select("otm.*", "o.name AS officer_name",
		"COUNT(DISTINCT omc.officership_member_id) AS current_officers",
		"COUNT(DISTINCT omp.officership_member_id) AS previous_officers").
		From("people.officership_team_members otm").
		LeftJoin("people.officerships o on o.officer_id = otm.officer_id").
		LeftJoin("people.officership_members omc ON o.officer_id = omc.officer_id AND omc.end_date IS NULL").
		LeftJoin("people.officership_members omp ON o.officer_id = omp.officer_id AND omp.end_date IS NOT NULL").
		Where(sq.And{
			sq.Eq{"otm.team_id": officershipTeamMemberGet.TeamID},
			sq.Eq{"otm.officer_id": officershipTeamMemberGet.OfficerID},
		}).
		GroupBy("otm", "otm.officer_id", "otm.team_id", "o.officer_id", "name", "email_alias", "description",
			"historywiki_url", "role_id", "is_current", "if_unfilled").
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getOfficershipTeamMember: %w", err))
	}

	err = s.db.GetContext(ctx, &m, sql, args...)
	if err != nil {
		return OfficershipTeamMember{}, fmt.Errorf("failed to get officership team member: %w", err)
	}

	return m, nil
}

func (s *Store) AddOfficershipTeamMember(ctx context.Context, officershipTeamMemberAdd OfficershipTeamMemberAddDTO) (OfficershipTeamMember, error) {
	builder := utils.PSQL().Insert("people.officership_team_members").
		Columns("team_id", "officer_id", "is_leader", "is_deputy").
		Values(officershipTeamMemberAdd.TeamID, officershipTeamMemberAdd.OfficerID, officershipTeamMemberAdd.IsLeader, officershipTeamMemberAdd.IsDeputy).
		Suffix("RETURNING team_id, officer_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for addOfficershipTeamMember: %w", err))
	}

	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return OfficershipTeamMember{}, fmt.Errorf("failed to add officership team member: %w", err)
	}

	defer stmt.Close()

	var teamID, officerID int

	err = stmt.QueryRow(args...).Scan(&teamID, &officerID)
	if err != nil {
		return OfficershipTeamMember{}, fmt.Errorf("failed to add offciership team member: %w", err)
	}

	return s.GetOfficershipTeamMember(ctx, OfficershipTeamMemberGetDeleteDTO{
		TeamID:    teamID,
		OfficerID: officerID,
	})
}

func (s *Store) DeleteOfficershipTeamMember(ctx context.Context, officershipTeamMemberDelete OfficershipTeamMemberGetDeleteDTO) error {
	builder := utils.PSQL().Delete("people.officership_team_members").
		Where(sq.And{sq.Eq{"team_id": officershipTeamMemberDelete.TeamID}, sq.Eq{"officer_id": officershipTeamMemberDelete.OfficerID}})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for deleteOfficershipTeam: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete officership team: %w", err)
	}

	return nil
}

func (s *Store) RemoveTeamForOfficershipTeamMembers(ctx context.Context, officershipTeamID int) error {
	builder := utils.PSQL().Delete("people.officership_team_members").
		Where(sq.Eq{"team_id": officershipTeamID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for removeTeamForOfficershipTeamMembers: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to remove team for officership team members: %w", err)
	}

	return nil
}

func (s *Store) GetOfficershipMembers(ctx context.Context, officershipGet *OfficershipGetDTO, userID *int,
	officershipStatus OfficershipsStatus, officershipMemberStatus OfficershipsStatus,
	orderByOfficerName bool) ([]OfficershipMemberDB, error) {
	var o []OfficershipMemberDB

	builder := utils.PSQL().Select("om.*", "o.name AS officership_name",
		"CONCAT(u.first_name, ' ', u.last_name) AS user_name", "otm.team_id AS team_id", "ot.name AS team_name").
		From("people.officership_members om").
		LeftJoin("people.officerships o ON o.officer_id = om.officer_id").
		LeftJoin("people.officership_team_members otm ON otm.officer_id = om.officer_id").
		LeftJoin("people.officership_teams ot ON ot.team_id = otm.team_id").
		LeftJoin("people.users u ON u.user_id = om.user_id")

	if officershipGet != nil {
		builder = builder.Where(sq.Or{
			sq.Eq{"o.officer_id": officershipGet.OfficershipID},
			sq.And{
				sq.Eq{"o.name": officershipGet.Name}, sq.NotEq{"o.name": ""},
			},
		})
	}

	if userID != nil {
		builder = builder.Where(sq.Eq{"u.user_id": userID})
	}

	switch officershipStatus {
	case Any:
	case Current:
		builder = builder.Where("o.is_current = true")
	case Retired:
		builder = builder.Where("o.is_current = false")
	}

	switch officershipMemberStatus {
	case Any:
	case Current:
		builder = builder.Where("om.end_date IS NULL")
	case Retired:
		builder = builder.Where("om.end_date IS NOT NULL")
	}

	if orderByOfficerName {
		builder = builder.OrderBy(`CASE WHEN o.name = 'Station Director' THEN 0
		WHEN o.name LIKE '%Director%' AND o.name NOT LIKE '%Deputy%' AND o.name NOT LIKE '%Assistant%' THEN 1
		WHEN o.name LIKE '%Deputy%' THEN 2
		WHEN o.name LIKE '%Assistant%' THEN 3
		WHEN o.name = 'Head of Welfare and Training' THEN 4
		WHEN o.name LIKE '%Head of%' THEN 5
		ELSE 6 END`, "o.name")
	}

	builder = builder.OrderBy("om.start_date DESC")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getOfficershipMembers: %w", err))
	}

	err = s.db.SelectContext(ctx, &o, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get officership members: %w", err)
	}

	return o, nil
}

func (s *Store) GetOfficershipMember(ctx context.Context, officershipMemberID int) (OfficershipMemberDB, error) {
	var m OfficershipMemberDB

	builder := utils.PSQL().Select("om.*", "o.name AS officership_name",
		"CONCAT(u.first_name, ' ', u.last_name) AS user_name", "otm.team_id AS team_id", "ot.name AS team_name").
		From("people.officership_members om").
		LeftJoin("people.officerships o ON o.officer_id = om.officer_id").
		LeftJoin("people.officership_team_members otm ON otm.officer_id = om.officer_id").
		LeftJoin("people.officership_teams ot ON ot.team_id = otm.team_id").
		LeftJoin("people.users u ON u.user_id = om.user_id").
		Where(sq.Eq{"om.officership_member_id": officershipMemberID}).
		Limit(1)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for getOfficershipMember: %w", err))
	}

	err = s.db.GetContext(ctx, &m, sql, args...)
	if err != nil {
		return OfficershipMemberDB{}, fmt.Errorf("failed to get officership member: %w", err)
	}

	return m, nil
}

func (s *Store) AddOfficershipMember(ctx context.Context, officershipMemberAdd OfficershipMemberAddEditDTO) (OfficershipMemberDB, error) {
	builder := utils.PSQL().Insert("people.officership_members").
		Columns("user_id", "officer_id", "start_date", "end_date").
		Values(officershipMemberAdd.UserID, officershipMemberAdd.OfficerID, officershipMemberAdd.StartDate, officershipMemberAdd.EndDate).
		Suffix("RETURNING officership_member_id")

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for addOfficershipMember: %w", err))
	}

	stmt, err := s.db.PrepareContext(ctx, sql)
	if err != nil {
		return OfficershipMemberDB{}, fmt.Errorf("failed to add officership member: %w", err)
	}

	defer stmt.Close()

	var officershipMemberID int

	err = stmt.QueryRow(args...).Scan(&officershipMemberID)
	if err != nil {
		return OfficershipMemberDB{}, fmt.Errorf("failed to add offciership member: %w", err)
	}

	return s.GetOfficershipMember(ctx, officershipMemberID)
}

func (s *Store) EditOfficershipMember(ctx context.Context, officershipMemberID int, officershipMemberAdd OfficershipMemberAddEditDTO) (OfficershipMemberDB, error) {
	builder := utils.PSQL().Update("people.officership_members").
		SetMap(map[string]interface{}{
			"user_id":    officershipMemberAdd.UserID,
			"officer_id": officershipMemberAdd.OfficerID,
			"start_date": officershipMemberAdd.StartDate,
			"end_date":   officershipMemberAdd.EndDate,
		}).
		Where(sq.Eq{"officership_member_id": officershipMemberID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for editOfficershipMember: %w", err))
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return OfficershipMemberDB{}, fmt.Errorf("failed to edit officership member: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return OfficershipMemberDB{}, fmt.Errorf("failed to edit officership member: %w", err)
	}

	if rows < 1 {
		return OfficershipMemberDB{},
			fmt.Errorf("failed to edit officerhip member: invalid rows affected: %d", rows)
	}

	return s.GetOfficershipMember(ctx, officershipMemberID)
}

func (s *Store) DeleteOfficershipMember(ctx context.Context, officershipMemberID int) error {
	builder := utils.PSQL().Delete("people.officership_members").
		Where(sq.Eq{"officership_member_id": officershipMemberID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for deleteOfficershipMember: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete officership member: %w", err)
	}

	return nil
}

func (s *Store) RemoveOfficershipForOfficershipMembers(ctx context.Context, officershipID int) error {
	builder := utils.PSQL().Delete("people.officership_members").
		Where(sq.Eq{"officer_id": officershipID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for removeOfficershipForOfficershipMembers: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to remove officership for officership members: %w", err)
	}

	return nil
}

func (s *Store) RemoveUserForOfficershipMembers(ctx context.Context, userID int) error {
	builder := utils.PSQL().Delete("people.officership_members").
		Where(sq.Eq{"user_id": userID})

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(fmt.Errorf("failed to build sql for removeUsersForOfficershipMembers: %w", err))
	}

	_, err = s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to remove users for officership members: %w", err)
	}

	return nil
}
