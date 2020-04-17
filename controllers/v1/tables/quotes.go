package tables

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/models"
	"github.com/ystv/web-api/services"
)

// QuoteCreate Quote create API
// @Summary Quote create API
// @Description Create new quote
// @Accept json
// @Produce json
// @Param	body	body	models.Quote	true	"quote create parameter"
// @Success 200 {object} models.Quote
// @Router /v1/tables/quotes [post]
func QuoteCreate(c echo.Context) error {
	q := new(models.Quote)
	err := c.Bind(q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	id, err := services.QuoteCreate(q)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, id)
}

// QuoteList Quote list API
// @Summary Quote list API
// @Description list quotes
// @Accept json
// @Produce json
// @Success 200 {object} models.QuoteSlice
// @Router /v1/tables/quotes [get]
func QuoteList(c echo.Context) error {
	res, err := services.QuoteList()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// QuoteFind Quote find API
// @Summary Quote find API
// @Description find quote
// @Accept json
// @Produce json
// @Success 200 {object} models.Quote
// @Router /v1/tables/quotes/{quote_id} [get]
func QuoteFind(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	res, err := services.QuoteFind(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// QuoteUpdate Quote update API
// @Summary Quote update API
// @Description update quotes
// @Accept  json
// @Produce  json
// @Param   quote_id     path    int     true        "quote id parameter"
// @Success 200 string string	""
// @Router /v1/tables/quotes/{quote_id} [put]
func QuoteUpdate(c echo.Context) error {
	// Check new quote will bind
	newQuote := new(models.Quote)
	err := c.Bind(newQuote)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// Finding quote to update
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	oldQuote, err := services.QuoteFind(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	// Update quote
	err = services.QuoteUpdate(oldQuote, newQuote)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, oldQuote)
}

// QuoteDelete Quote delete API
// @Summary Quote delete API
// @Description delete quotes
// @Accept  json
// @Produce  json
// @Param   quote_id     path    int     true        "quote id parameter"
// @Success 200 string string	""
// @Router /v1/tables/quotes/{quote_id} [delete]
func QuoteDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	res, err := services.QuoteDelete(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
