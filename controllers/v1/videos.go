package v1

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/ystv/web-api/models"
)

// VideoCreate Video create API
// @Summary Video create API
// @Description create new video, there is more to videos!
// @Accept json
// @Produce json
// @Param	body	body	v1.videoCreateReq	true	"video create parameter"
// @Success 200 {object} models.VideoCreate
// @Router /v1/videos [post]
func VideoCreate(c echo.Context) error {
	v := &videoCreateReq{}
	if err := c.Bind(v); err != nil {
		return err
	}
	return c.JSON(http..StatusOK, map[string]uint{"video_id": videoID})
}

// VideoList Video list API
// @Summary User list API
// @Description list videos
// @Accept json
// @Produce json
// @Success 200 {object} models.VideoList
// @Router /v1/videos [get]
func VideoList(c echo.Context) error {
	res, err := services.VideoList()
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
// @Router /v1/videos/{video_id} [put]
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
// @Router /v1/videos/{video_id} [delete]
func videoDelete(c echo.Context) error {
	return c.JSON(http.StatusOK, "delete ok")
}