package misc

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/utils"
)

// ListWebcams handles listing all webcams a user can access
// @Summary List webcams
// @Description List webcams available to user by using the permission ID
// @ID list-webcams
// @Tags misc-webcams
// @Success 200 {array} misc.Webcam
// @Router /v1/internal/misc/webcams [get]
func (s *Store) ListWebcams(c echo.Context) error {
	claims, status, err := s.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("ListWebcams failed to get user token: %w", err)
		return echo.NewHTTPError(status, err)
	}

	var perms []string
	perms = append(perms, claims.Permissions...)

	w, err := s.misc.ListWebcams(c.Request().Context(), perms)
	if err != nil {
		err = fmt.Errorf("failed to list webcams: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, utils.NonNil(w))
}

// GetWebcam handles reverse proxying a webcam
// @Summary Get webcam
// @Description Reverse proxies the selected webcam returns the jpeg feed as a result.
// @ID get-webcam
// @Tags misc-webcams
// @Param cameraID path int true "Camera ID"
// @Router /v1/internal/misc/webcams/{cameraID} [get]
func (s *Store) GetWebcam(c echo.Context) error {
	cameraID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid camera ID")
	}

	claims, status, err := s.access.GetToken(c.Request())
	if err != nil {
		err = fmt.Errorf("GetWebcam failed to get user token: %w", err)
		return echo.NewHTTPError(status, err)
	}

	var perms []string
	perms = append(perms, claims.Permissions...)

	w, err := s.misc.GetWebcam(c.Request().Context(), cameraID, perms)
	if err != nil {
		err = fmt.Errorf("failed to get camera: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	target, err := url.Parse(w.URL)
	if err != nil {
		err = fmt.Errorf("failed to parse webcam URL: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	targetQuery := target.RawQuery
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		req.URL.Path = path.Base(req.URL.Path)
		req.URL.RawPath = path.Base(req.URL.Path)

		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

func singleJoiningSlash(a, b string) string {
	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")
	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(proxy, req *url.URL) (path, rawPath string) {
	req.Path = strings.TrimSuffix(req.Path, proxy.Path)
	if proxy.RawPath == "" && req.RawPath == "" {
		return singleJoiningSlash(proxy.Path, req.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	aPath := proxy.EscapedPath()
	bPath := req.EscapedPath()

	aSlash := strings.HasSuffix(aPath, "/")
	bSlash := strings.HasPrefix(bPath, "/")

	switch {
	case aSlash && bSlash:
		return proxy.Path + req.Path[1:], aPath + bPath[1:]
	case !aSlash && !bSlash:
		return proxy.Path + "/" + req.Path, aPath + "/" + bPath
	}
	return proxy.Path + req.Path, aPath + bPath
}
