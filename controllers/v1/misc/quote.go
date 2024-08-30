package misc

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/misc"
	"github.com/ystv/web-api/utils"
)

// ListQuotes handles listing quotes by pagination
// @Summary List quotes
// @Description Lists quotes by pagination.
// @ID get-quotes
// @Tags misc-quotes
// @Produce json
// @Param amount path int true "Amount"
// @Param page path int true "Page"
// @Success 200 {array} misc.QuotePage
// @Router /v1/internal/misc/quotes/{amount}/{page} [get]
func (r *Repos) ListQuotes(c echo.Context) error {
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
		err = fmt.Errorf("ListQuotes failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(q))
}

// NewQuote handles creating a quote
// @Summary New quote
// @Description creates a new quote.
// @Description web-api will overwrite created by User ID with the token's user ID.
// @ID new-quote
// @Tags misc-quotes
// @Accept json
// @Param quote body misc.Quote true "Quote object"
// @Success 201 {object} int "Quote ID"
// @Router /v1/internal/misc/quotes [post]
func (r *Repos) NewQuote(c echo.Context) error {
	q := misc.Quote{}
	err := c.Bind(&q)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	claims, err := r.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("NewQuote failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	q.CreatedBy = claims.UserID
	err = r.misc.NewQuote(c.Request().Context(), q)
	if err != nil {
		err = fmt.Errorf("NewQuote failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

// UpdateQuote handles updating a quote
// @Summary Update quote
// @Description updates a quote. Still need to provide the whole Quote object,
// @Description web-api will overwrite created by User ID to keep with existing record.
// @ID update-quote
// @Tags misc-quotes
// @Accept json
// @Param quote body misc.Quote true "Quote object"
// @Success 200
// @Router /v1/internal/misc/quotes [put]
func (r *Repos) UpdateQuote(c echo.Context) error {
	q := misc.Quote{}
	err := c.Bind(&q)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = r.misc.UpdateQuote(c.Request().Context(), q)
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
// @Tags misc-quotes
// @Param quoteid path int true "Quote ID"
// @Success 200
// @Router /v1/internal/misc/quotes/{quoteid} [delete]
func (r *Repos) DeleteQuote(c echo.Context) error {
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
