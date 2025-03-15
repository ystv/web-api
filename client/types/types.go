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

	FindStreamEndpointOptions struct {
		// EndpointID is the unique database id of the stream
		EndpointID *int `json:"endpointId,omitempty"`
		// Application defines which RTMP application this is valid for
		Application *string `json:"application,omitempty"`
		// Name is the unique name given in an application
		Name *string `json:"name,omitempty"`
		// Pwd defines an extra layer of security for authentication
		Pwd *string `json:"pwd,omitempty"`
	}

	ListUsersPaginationOptions struct {
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
	FIVE UserPaginationSize = 5
	TEN  UserPaginationSize = 10
	//nolint:revive
	TWENTY_FIVE UserPaginationSize = 25
	FIFTY       UserPaginationSize = 50
	//nolint:revive
	SEVENTY_FIVE UserPaginationSize = 75
	//nolint:revive
	ONE_HUNDRED UserPaginationSize = 100
)

const (
	ASCSENDING UserPaginationDirection = "asc"
	DESCENDING UserPaginationDirection = "desc"
)

const (
	//nolint:revive
	USER_ID  UserPaginationColumn = "userId"
	NAME     UserPaginationColumn = "name"
	USERNAME UserPaginationColumn = "username"
	EMAIL    UserPaginationColumn = "email"
	//nolint:revive
	LAST_LOGIN UserPaginationColumn = "lastLogin"
)

const (
	ENABLED  UserPaginationEnabled = "enabled"
	DISABLED UserPaginationEnabled = "disabled"
)

const (
	DELETED UserPaginationDeleted = "deleted"
	//nolint:revive
	NOT_DELETED UserPaginationDeleted = "not_deleted"
)
