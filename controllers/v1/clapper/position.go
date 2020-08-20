package clapper

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services/clapper/position"
)

// PositionList handles listing all possible positions
func PositionList(c echo.Context) error {
	p, err := position.List()
	if err != nil {
		log.Printf("PositionList failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// PositionNew handles creating a new position
func PositionNew(c echo.Context) error {
	p := position.Position{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = position.New(&p)
	if err != nil {
		log.Printf("PositionNew failed: %+v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

// PositionUpdate updates an existing position
func PositionUpdate(c echo.Context) error {
	p := position.Position{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = position.Update(&p)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, err)
		}
		log.Printf("PositionUpdate failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
