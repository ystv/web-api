package misc

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/misc"
	"github.com/ystv/web-api/utils"
)

// Repos stores our dependencies
type (
	Repos interface {
		ListRepo
		QuoteRepo
		WebcamRepo
	}

	ListRepo interface {
		GetLists(c echo.Context) error
		GetListsByToken(c echo.Context) error
		GetList(c echo.Context) error
		GetSubscribers(c echo.Context) error
		SubscribeByToken(c echo.Context) error
		SubscribeByID(c echo.Context) error
		UnsubscribeByToken(c echo.Context) error
		UnsubscribeByID(c echo.Context) error
		UnsubscribeByUUID(c echo.Context) error
	}

	QuoteRepo interface {
		ListQuotes(c echo.Context) error
		NewQuote(c echo.Context) error
		UpdateQuote(c echo.Context) error
		DeleteQuote(c echo.Context) error
	}

	WebcamRepo interface {
		ListWebcams(c echo.Context) error
		GetWebcam(c echo.Context) error
	}

	Store struct {
		misc   misc.Repos
		access utils.Repo
	}
)

// NewRepos creates our data store
func NewRepos(db *sqlx.DB, access utils.Repo) Repos {
	return &Store{
		misc:   misc.NewStore(db),
		access: access,
	}
}
