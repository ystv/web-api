package clapper

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/clapper"
	"github.com/ystv/web-api/utils"
)

// ListPositions handles listing all possible positions
// @Summary List positions
// @Description Lists all positions.
// @ID get-positions
// @Tags clapper-positions
// @Produce json
// @Success 200 {array} clapper.Position
// @Router /v1/internal/clapper/positions [get]
func (s *Store) ListPositions(c echo.Context) error {
	p, err := s.position.List(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListPositions: failed to list: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, utils.NonNil(p))
}

// NewPosition handles creating a new position
// @Summary New position
// @ID new-position
// @Tags clapper-positions
// @Accept json
// @Param event body clapper.Position true "Position object"
// @Success 201 body int "Position ID"
// @Router /v1/internal/clapper/positions [post]
func (s *Store) NewPosition(c echo.Context) error {
	var p clapper.Position

	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("NewPosition: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	positionID, err := s.position.New(c.Request().Context(), &p)
	if err != nil {
		err = fmt.Errorf("NewPosition: failed to insert position: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, positionID)
}

// UpdatePosition updates an existing position
// @Summary Update position
// @ID update-position
// @Tags clapper-positions
// @Accept json
// @Param quote body clapper.Position true "Position object"
// @Success 200
// @Router /v1/internal/clapper/positions [put]
func (s *Store) UpdatePosition(c echo.Context) error {
	var p clapper.Position

	err := c.Bind(&p)
	if err != nil {
		err = fmt.Errorf("UpdatePosition: failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = s.position.Update(c.Request().Context(), &p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		err = fmt.Errorf("UpdatePosition failed: %w", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// DeletePosition removes a position
//
// @Summary Delete position
// @ID delete-position
// @Tags clapper-positions
// @Accept json
// @Param positionid path int true "Position ID"
// @Success 200
// @Router /v1/internal/clapper/positions/{positionid} [delete]
func (s *Store) DeletePosition(c echo.Context) error {
	positionID, err := strconv.Atoi(c.Param("positionid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid position ID")
	}

	err = s.position.Delete(c.Request().Context(), positionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to delete event: %w", err))
	}

	return c.NoContent(http.StatusOK)
}
