package misc

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/misc"
)

// Repo stores our dependencies
type Repo struct {
	misc misc.Repo
}

// NewRepo creates our data store
func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{misc.NewStore(db)}
}

// ListQuotes handles listing quotes by pagination
func (r *Repo) ListQuotes(c echo.Context) error {
	amount, err := strconv.Atoi(c.Param("amount"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad offset")
	}
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad page")
	}
	q, err := r.misc.ListQuotes(c.Request().Context(), amount, page)
	if err != nil {
		log.Printf("PresetList failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, q)
}

// NewQuote handles creating a quote
func (r *Repo) NewQuote(c echo.Context) error {
	q := misc.Quote{}
	err := c.Bind(&q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = r.misc.NewQuote(c.Request().Context(), q)
	if err != nil {
		log.Printf("NewQuote failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

// UpdateQuote handles updating a quote
func (r *Repo) UpdateQuote(c echo.Context) error {
	q := misc.Quote{}
	err := c.Bind(&q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = r.misc.NewQuote(c.Request().Context(), q)
	if err != nil {
		log.Printf("UpdateQuote failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

// DeleteQuote handles deleting quotes
func (r *Repo) DeleteQuote(c echo.Context) error {
	quoteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad ID")
	}
	err = r.misc.DeleteQuote(c.Request().Context(), quoteID)
	if err != nil {
		log.Printf("PresetList failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
