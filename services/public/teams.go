package public

type (
	Team struct {
		Name    string       `json:"name"`
		Members []TeamMember `json:"members"`
	}
	TeamMember struct {
		Name     string `json:"name"`
		Position string `json:"position"`
		Email    string `json:"email"`
	}
)
