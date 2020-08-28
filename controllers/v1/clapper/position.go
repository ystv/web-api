package clapper

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/clapper"
)

// PositionList handles listing all possible positions
func (r *Repos) PositionList(c echo.Context) error {
	p, err := r.position.List(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("PositionList: failed to list: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// PositionNew handles creating a new position
func (r *Repos) PositionNew(c echo.Context) error {
	p := clapper.Position{}
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("PositionNew: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	positionID, err := r.position.New(c.Request().Context(), &p)
	if err != nil {
		err = fmt.Errorf("PositionNew: failed to insert position: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, positionID)
}

// PositionUpdate updates an existing position
func (r *Repos) PositionUpdate(c echo.Context) error {
	p := clapper.Position{}
	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("PositionUpdate: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = r.position.Update(c.Request().Context(), &p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		err = fmt.Errorf("PositionUpdate failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// TODO add delete
