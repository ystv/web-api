package public

import "context"

type (
	//Team a group in ystv
	Team struct {
		Name        string       `json:"name"`
		Description string       `json:"description"`
		Members     []TeamMember `json:"members"`
	}
	//TeamMember a position within a group
	TeamMember struct {
		Name     string `json:"name"`
		Position string `json:"position"`
		Email    string `json:"email"`
	}
)

var _ TeamRepo = &Store{}

// ListTeams returns a list of the ystv teams and their current members.
func (s *Store) ListTeams(ctx context.Context) (*[]Team, error) {
	t := []Team{
		{
			"Admin Team",
			"The bossy people",
			[]TeamMember{
				{
					"Person A",
					"Station Director",
					"station.director@ystv.co.uk",
				},
				{
					"Person B",
					"Station Manager",
					"station.manager@ystv.co.uk",
				},
			},
		},
		{
			"Comp Team",
			"The cool kids",
			[]TeamMember{
				{
					"Person C",
					"Computing Director",
					"computing.director@ystv.co.uk",
				},
			},
		},
	}
	return &t, nil
}
