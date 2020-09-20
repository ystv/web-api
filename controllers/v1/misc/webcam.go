package misc

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/controllers/v1/people"
)

// ListWebcams handles listing all webcams a user can access
// @Summary List webcams
// @Description List webcams available to user by using the permission ID
// @ID list-webcams
// @Tags misc, webcams
// @Success 200 {array} misc.Webcam
// @Router /v1/internal/misc/webcams [get]
func (r *Repos) ListWebcams(c echo.Context) error {
	// Get user token
	claims, err := people.GetToken(c)
	if err != nil {
		err = fmt.Errorf("ListWebcams failed to get user token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// Prepare slice of permission IDs
	perms := []int{}
	for _, permission := range claims.Permissions {
		perms = append(perms, permission.PermissionID)
	}
	w, err := r.misc.ListWebcams(c.Request().Context(), perms)
	if err != nil {
		err = fmt.Errorf("failed to list webcams: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, w)
}

// GetWebcam handles reverse proxying a webcam
// @Summary Get webcam
// @Description Reverse proxies the selected webcam returns the jpeg feed as a result.
// @ID get-webcam
// @Tags misc, webcams
// @Param cameraid path int true "Camera ID"
// @Router /v1/internal/misc/webcams{cameraid} [get]
func (r *Repos) GetWebcam(c echo.Context) error {
	cameraID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid camera ID")
	}
	// Get user token
	claims, err := people.GetToken(c)
	if err != nil {
		err = fmt.Errorf("ListWebcams failed to get user token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// Prepare slice of permission IDs
	perms := []int{}
	for _, permission := range claims.Permissions {
		perms = append(perms, permission.PermissionID)
	}
	// Get webcam URL and check user has permission for it
	w, err := r.misc.GetWebcam(c.Request().Context(), cameraID, perms)
	if err != nil {
		err = fmt.Errorf("failed to get camera: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	url, err := url.Parse(w.URL)
	if err != nil {
		err = fmt.Errorf("failed to parse webcam URL: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)

	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}
