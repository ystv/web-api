package misc

import (
	"fmt"
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
// @Summary List quotes
// @Description Lists quotes by pagination.
// @ID get-quotes
// @Tags quotes
// @Produce json
// @Param amount path int true "Amount"
// @Param page path int true "Page"
// @Success 200 {array} misc.QuotePage
// @Router /v1/internal/misc/quote/{amount}/{page} [get]
func (r *Repo) ListQuotes(c echo.Context) error {
	amount, err := strconv.Atoi(c.Param("amount"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad offset")
	}
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad page")
	}
	q, err := r.misc.ListQuotes(c.Request().Context(), amount, page)
	if err != nil {
		err = fmt.Errorf("ListQuotes failed: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, q)
}

// NewQuote handles creating a quote
// @Summary New quote
// @Description creates a new quote.
// @ID new-quote
// @Tags quotes
// @Accept json
// @Param quote body misc.Quote true "Quote object"
// @Success 201 {object} int "Quote ID"
// @Router /v1/internal/misc/quote [post]
func (r *Repo) NewQuote(c echo.Context) error {
	q := misc.Quote{}
	err := c.Bind(&q)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = r.misc.NewQuote(c.Request().Context(), q)
	if err != nil {
		err = fmt.Errorf("NewQuote failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

// UpdateQuote handles updating a quote
// @Summary Update quote
// @Description updates a quote.
// @ID update-quote
// @Tags quotes
// @Accept json
// @Param quote body misc.Quote true "Quote object"
// @Success 200
// @Router /v1/internal/misc/quote [put]
func (r *Repo) UpdateQuote(c echo.Context) error {
	q := misc.Quote{}
	err := c.Bind(&q)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = r.misc.NewQuote(c.Request().Context(), q)
	if err != nil {
		err = fmt.Errorf("UpdateQuote failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// DeleteQuote handles deleting quotes
// @Summary Delete quote
// @Description deletes a quote by ID.
// @ID delete-quote
// @Tags quotes
// @Param quoteid path int true "Quote ID"
// @Success 200
// @Router /v1/internal/misc/quote/{quoteid} [delete]
func (r *Repo) DeleteQuote(c echo.Context) error {
	quoteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	err = r.misc.DeleteQuote(c.Request().Context(), quoteID)
	if err != nil {
		err = fmt.Errorf("DeleteQuote failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
