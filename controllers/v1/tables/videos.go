package tables

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/services"
)

// VideoCreate Video create API
// @Summary Video create API
// @Description create new video, there is more to videos!
// @Accept json
// @Produce json
// @Param	body	body	models.Video	true	"video create parameter"
// @Success 200 {object} models.Video
// @Router /v1/tables/videos [post]
func VideoCreate(c echo.Context) error {
	return c.JSON(http.StatusOK, "update ok")
}

// VideoList Video list API
// @Summary Video list API
// @Description list videos
// @Accept json
// @Produce json
// @Success 200 {object} models.VideoSlice
// @Router /v1/tables/videos [get]
func VideoList(c echo.Context) error {
	res, err := services.VideoList()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// VideoFind Video find API
// @Summary Video find API
// @Description find video
// @Accept json
// @Produce json
// @Success 200 {object} models.Video
// @Router /v1/tables/videos/{video_id} [get]
func VideoFind(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(400, "Number pls")
	}
	res, err := services.VideoFind(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// VideoUpdate Video update API
// @Summary Video update API
// @Description update videos
// @Accept  json
// @Produce  json
// @Param   video_id     path    int     true        "video id parameter"
// @Success 200 string string	""
// @Router /v1/tables/videos/{video_id} [put]
func VideoUpdate(c echo.Context) error {
	return c.JSON(http.StatusOK, "update ok")
}

// VideoDelete Video delete API
// @Summary Video delete API
// @Description delete videos
// @Accept  json
// @Produce  json
// @Param   video_id     path    int     true        "video id parameter"
// @Success 200 string string	""
// @Router /v1/tables/videos/{video_id} [delete]
func VideoDelete(c echo.Context) error {
	return c.JSON(http.StatusOK, "delete ok")
}
