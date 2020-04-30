package tables

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/models"
	"github.com/ystv/web-api/services"
)

// VideoBoxCreate VideoBox create API
// @Summary VideoBox create API
// @Description Create new videobox
// @Accept json
// @Produce json
// @Param	body	body	models.VideoBox	true	"videobox create parameter"
// @Success 200 {object} models.VideoBox
// @Router /v1/tables/videoboxes [post]
func VideoBoxCreate(c echo.Context) error {
	q := new(models.VideoBox)
	err := c.Bind(q)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	id, err := services.VideoBoxCreate(q)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, id)
}

// VideoBoxList VideoBox list API
// @Summary VideoBox list API
// @Description list videoboxes
// @Accept json
// @Produce json
// @Success 200 {object} models.VideoBoxSlice
// @Router /v1/tables/videoboxes [get]
func VideoBoxList(c echo.Context) error {
	res, err := services.VideoBoxList()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// VideoBoxFind VideoBox find API
// @Summary VideoBox find API
// @Description find videobox
// @Accept json
// @Produce json
// @Success 200 {object} models.VideoBox
// @Router /v1/tables/videoboxes/{videobox_id} [get]
func VideoBoxFind(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	res, err := services.VideoBoxFind(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// VideoBoxUpdate VideoBox update API
// @Summary VideoBox update API
// @Description update videoboxes
// @Accept  json
// @Produce  json
// @Param   videobox_id     path    int     true        "videobox id parameter"
// @Success 200 string string	""
// @Router /v1/tables/videoboxes/{videobox_id} [put]
func VideoBoxUpdate(c echo.Context) error {
	// Check new videobox will bind
	newVideoBox := new(models.VideoBox)
	err := c.Bind(newVideoBox)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// Finding videobox to update
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	oldVideoBox, err := services.VideoBoxFind(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	// Update videobox
	err = services.VideoBoxUpdate(oldVideoBox, newVideoBox)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, oldVideoBox)
}

// VideoBoxDelete VideoBox delete API
// @Summary VideoBox delete API
// @Description delete videoboxes
// @Accept  json
// @Produce  json
// @Param   videobox_id     path    int     true        "videobox id parameter"
// @Success 200 string string	""
// @Router /v1/tables/videoboxes/{videobox_id} [delete]
func VideoBoxDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Number pls")
	}
	res, err := services.VideoBoxDelete(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
