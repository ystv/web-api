package types

type (
	ErrorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	SuccessResponse struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}

	UsersListPaginationOptions struct {
		// Size indicates the size of the pagination to return.
		//
		// If you want all the results then leave Size nil.
		// If Size is not nil, then Page will need to be set as well
		Size *UserPaginationSize `json:"size,omitempty"`
		// Page is used in combination with Size, and both must be either
		// set or unset
		Page *int `json:"page,omitempty"`
		// Search is a search string that you want to query against
		Search *string `json:"search,omitempty"`
		// Column is used to determine the order of the output.
		// If Column is set then Direction must be set as well
		Column *UserPaginationColumn `json:"column,omitempty"`
		// Direction is either asc or desc and is only used when
		// Column is used
		Direction *UserPaginationDirection `json:"direction,omitempty"`
		// Enabled checks if the users you are listing are enabled
		// or disabled.
		//
		// If you want any user regardless, then leave Enabled nil
		Enabled *UserPaginationEnabled `json:"enabled,omitempty"`
		// Deleted checks if the users you are listing are deleted
		// or not.
		//
		// If you want any user regardless, then leave Deleted nil
		Deleted *UserPaginationDeleted `json:"deleted,omitempty"`
	}

	UserPaginationSize      int
	UserPaginationDirection string
	UserPaginationColumn    string
	UserPaginationEnabled   string
	UserPaginationDeleted   string
)

const (
	FIVE         UserPaginationSize = 5
	TEN          UserPaginationSize = 10
	TWENTY_FIVE  UserPaginationSize = 25
	FIFTY        UserPaginationSize = 50
	SEVENTY_FIVE UserPaginationSize = 75
	ONE_HUNDRED  UserPaginationSize = 100
)

const (
	ASCSENDING UserPaginationDirection = "asc"
	DESCENDING UserPaginationDirection = "desc"
)

const (
	USER_ID    UserPaginationColumn = "userId"
	NAME       UserPaginationColumn = "name"
	USERNAME   UserPaginationColumn = "username"
	EMAIL      UserPaginationColumn = "email"
	LAST_LOGIN UserPaginationColumn = "lastLogin"
)

const (
	ENABLED  UserPaginationEnabled = "enabled"
	DISABLED UserPaginationEnabled = "disabled"
)

const (
	DELETED     UserPaginationDeleted = "deleted"
	NOT_DELETED UserPaginationDeleted = "not_deleted"
)
