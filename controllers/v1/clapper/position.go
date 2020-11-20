package clapper

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/clapper"
)

// ListPosition handles listing all possible positions
// @Summary List positions
// @Description Lists all positions.
// @ID get-positions
// @Tags positions
// @Produce json
// @Success 200 {array} clapper.Position
// @Router /v1/internal/clapper/positions [get]
func (r *Repos) ListPosition(c echo.Context) error {
	p, err := r.position.List(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListPosition: failed to list: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// NewPosition handles creating a new position
// @Summary New position
// @Description creates a new position.
// @ID new-position
// @Tags positions
// @Accept json
// @Param event body clapper.Position true "Position object"
// @Success 201 body int "Position ID"
// @Router /v1/internal/clapper/positions [post]
func (r *Repos) NewPosition(c echo.Context) error {
	p := clapper.Position{}
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("NewPosition: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	positionID, err := r.position.New(c.Request().Context(), &p)
	if err != nil {
		err = fmt.Errorf("NewPosition: failed to insert position: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, positionID)
}

// UpdatePosition updates an existing position
// @Summary Update position
// @Description updates a position.
// @ID update-position
// @Tags positions
// @Accept json
// @Param quote body clapper.Position true "Position object"
// @Success 200
// @Router /v1/internal/clapper/positions [put]
func (r *Repos) UpdatePosition(c echo.Context) error {
	p := clapper.Position{}
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("UpdatePosition: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = r.position.Update(c.Request().Context(), &p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		err = fmt.Errorf("UpdatePosition failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// TODO add delete
