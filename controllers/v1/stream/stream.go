package stream

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"

	"github.com/ystv/web-api/services/stream"
	"github.com/ystv/web-api/utils"
)

// Repos encapsulates the dependency
type Repos struct {
	stream stream.Repo
}

// NewRepos creates our data store
func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{stream.NewStore(db)}
}

// PublishStream handles a stream publish request
//
// @Summary Publish a stream
// @Description Checks existing stream endpoints and changes it to active; this is for Nginx RTMP module
// @Description containing the application, name and pwd
// @ID publish-stream
// @Tags stream-endpoints
// @Accept json
// @P aram event body
// @Success 200 body int "EndpointDB ID"
// @Error 401
// @Router /v1/internal/stream/publish [post]
func (r *Repos) PublishStream(c echo.Context) error {
	var application, name, pwd, action string
	var err error

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(c.Request().Body)

	if c.Request().Header.Get("Content-Type") == "application/json" {
		// SRS publish handler
		application, name, pwd, action, err = _handleSRSPublish(c)
		if action != "on_publish" {
			err = errors.New("invalid action " + action)
		}
	} else {
		// Form DATA from nginx-rtmp/srtrelay
		application, name, pwd, action = _handleNginxPublish(c)

		// only apply auth for a publish request
		if action != "publish" {
			return nil
		}
	}

	if err != nil {
		c.Logger().Warnf("PublishStream: failed to parse publish data: %+v", err)
		return c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	endpoint, err := r.stream.GetEndpointByApplicationNamePwd(c.Request().Context(), application, name, pwd)
	if err != nil {
		c.Logger().Errorf("PublishStream: failed to get endpoint: %+v", err)
		return c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	if endpoint.Active || endpoint.Blocked {
		c.Logger().Warnf("PublishStream: endpoint active (%t) or blocked (%t)", endpoint.Active, endpoint.Blocked)
		return c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	err = r.stream.SetEndpointActiveByID(c.Request().Context(), endpoint.EndpointID)
	if err != nil {
		c.Logger().Errorf("PublishStream: failed to set endpoint active: %+v", err)
		return c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	c.Logger().Infof("PublishStream: published %d %s/%s", endpoint.EndpointID, application, name)

	// SRS needs zero response
	return c.String(http.StatusOK, "0")
}

// UnpublishStream handles a stream unpublish request
//
// @Summary Unpublish a stream
// @Description Checks existing stream endpoints and changes it to inactive; this is for Nginx RTMP module
// @Description containing the application, name, authentication and start and end times
// @ID unpublish-stream
// @Tags stream-endpoints
// @Accept json
// @P aram event body
// @Success 200 body int
// @Error 401
// @Router /v1/internal/stream/unpublish [post]
func (r *Repos) UnpublishStream(c echo.Context) error {
	var application, name, pwd, action string
	var err error

	if c.Request().Header.Get("Content-Type") == "application/json" {
		// SRS publish handler
		application, name, _, action, err = _handleSRSPublish(c)
		if action != "on_unpublish" {
			err = fmt.Errorf("invalid action %s", action)
		}
	} else {
		// Form DATA from nginx-rtmp/srtrelay
		application, name, pwd, action = _handleNginxPublish(c)
		// ignore actions except unpublish
		if action != "publish_done" {
			return nil
		}
	}

	if err != nil {
		c.Logger().Warnf("UnpublishStream: failed to parse unpublish data: %+v", err)
		return c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	err = r.stream.SetEndpointInactiveByApplicationNamePwd(c.Request().Context(), application, name, pwd)
	if err != nil {
		c.Logger().Errorf("UnpublishStream: failed to unpublish stream, continuing, %s/%s: %+v", application, name, err)
	}

	c.Logger().Infof("UnpublishStream: unpublished %s/%s", application, name)

	// SRS needs zero response
	return c.String(http.StatusOK, "0")
}

type _srsPublish struct {
	Action      string `json:"action"`
	IP          string `json:"ip"`
	VHost       string `json:"vhost"`
	Application string `json:"app"`
	URL         string `json:"tcUrl"`
	Stream      string `json:"stream"`
	Param       string `json:"param"`
}

func _handleSRSPublish(c echo.Context) (application, name, pwd, action string, err error) {
	var publish _srsPublish

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(c.Request().Body)
	dec := json.NewDecoder(c.Request().Body)
	err = dec.Decode(&publish)
	if err != nil {
		return
	}

	// skip question mark
	if len(publish.Param) > 0 {
		publish.Param = publish.Param[1:]
	}

	val, err := url.ParseQuery(publish.Param)
	if err != nil {
		return
	}
	application = publish.Application
	name = publish.Stream
	pwd = val.Get("pwd")
	action = publish.Action

	return
}

func _handleNginxPublish(c echo.Context) (application, name, pwd, action string) {
	application = c.FormValue("app")
	name = c.FormValue("name")
	pwd = c.FormValue("pwd")
	action = c.FormValue("call")

	return
}

// ListStreams handles a listing stream endpoints
//
// @Summary ListStreams stream endpoints
// @Description Lists all stream endpoints; this is for Nginx RTMP module
// @Description containing the application, name, authentication and start and end times
// @ID get-stream
// @Tags stream-endpoints
// @Accept json
// @Success 200 {array} stream.Endpoint
// @Router /v1/internal/streams [get]
func (r *Repos) ListStreams(c echo.Context) error {
	e, err := r.stream.ListEndpoints(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("ListStreams: failed to get: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	endpoints := make([]stream.Endpoint, 0)

	for _, endpoint := range e {
		var pwd, startValid, endValid, notes string

		if endpoint.Pwd.Valid {
			pwd = endpoint.Pwd.String
		}
		if endpoint.StartValid.Valid {
			startValid = endpoint.StartValid.Time.Format(time.RFC3339)
		}
		if endpoint.EndValid.Valid {
			endValid = endpoint.EndValid.Time.Format(time.RFC3339)
		}
		if endpoint.Notes.Valid {
			notes = endpoint.Notes.String
		}

		endpoints = append(endpoints, stream.Endpoint{
			EndpointID:  endpoint.EndpointID,
			Application: endpoint.Application,
			Name:        endpoint.Name,
			Pwd:         pwd,
			StartValid:  startValid,
			EndValid:    endValid,
			Notes:       notes,
			Active:      endpoint.Active,
			Blocked:     endpoint.Blocked,
		})
	}

	return c.JSON(http.StatusOK, utils.NonNil(endpoints))
}

// FindStream handles finding a stream
// @Summary Finds stream
// @Description finds existing stream
// @ID find-stream
// @Tags stream-endpoints
// @Accept json
// @Param endpoint body stream.FindEndpoint true "Find Endpoint object"
// @Success 200 {object} stream.EndpointDB
// @Error 400
// @Error 404
// @Router /v1/internal/streams/find [get]
func (r *Repos) FindStream(c echo.Context) error {
	var findEndpoint stream.FindEndpoint

	err := c.Bind(&findEndpoint)
	if err != nil {
		err = fmt.Errorf("failed to bind to request json: %w", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if findEndpoint.EndpointID == 0 && (len(findEndpoint.Application) == 0 || len(findEndpoint.Name) == 0) {
		err = errors.New("failed to bind to request json: missing application, name or endpoint id")
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	foundStream, err := r.stream.FindEndpoint(c.Request().Context(), &findEndpoint)
	if err != nil {
		err = fmt.Errorf("failed to find endpoint: %w", err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, foundStream)
}

// NewStream handles a creating a stream endpoint
//
// @Summary NewStream stream endpoint
// @Description Creates a new stream endpoint; this is for Nginx RTMP module
// @Description containing the application, name, authentication and start and end times
// @ID new-stream
// @Tags stream-endpoints
// @Accept json
// @Param endpoint body stream.NewEditEndpoint true "Stream endpoint object"
// @Success 201 body int "EndpointDB ID"
// @Error 400
// @Router /v1/internal/streams [post]
func (r *Repos) NewStream(c echo.Context) error {
	var newEndpoint stream.NewEditEndpoint

	err := c.Bind(&newEndpoint)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("NewStream: failed to bind to request json: %w", err))
	}

	if len(newEndpoint.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("NewStream: endpoint name must be set"))
	}

	if len(newEndpoint.Application) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("NewStream: endpoint application must be set"))
	}

	startTime := null.NewTime(time.Time{}, false)

	if newEndpoint.StartValid != "" {
		var parseStart time.Time

		parseStart, err = time.Parse(time.RFC3339, newEndpoint.StartValid)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("NewStream: failed to parse start date: %w", err))
		}

		startTime = null.TimeFrom(parseStart)
	}

	endTime := null.NewTime(time.Time{}, false)

	if newEndpoint.EndValid != "" {
		var parseEnd time.Time

		parseEnd, err = time.Parse(time.RFC3339, newEndpoint.EndValid)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("NewStream: failed to parse end date: %w", err))
		}

		diffEnd := time.Now().Compare(parseEnd)
		if diffEnd != -1 {
			return echo.NewHTTPError(http.StatusBadRequest, "NewStream: end date must be after now")
		}

		if startTime.Valid {
			diffStartEnd := startTime.Time.Compare(parseEnd)
			if diffStartEnd != -1 {
				return echo.NewHTTPError(http.StatusBadRequest, "NewStream: end date must be after start date")
			}
		}

		endTime = null.TimeFrom(parseEnd)
	}

	pwd := null.NewString("", false)
	notes := null.NewString("", false)

	if newEndpoint.Pwd != "" {
		pwd = null.StringFrom(newEndpoint.Pwd)
	}

	if newEndpoint.Notes != "" {
		notes = null.StringFrom(newEndpoint.Notes)
	}

	endpoint, err := r.stream.NewEndpoint(c.Request().Context(), stream.EndpointDB{
		Application: newEndpoint.Application,
		Name:        newEndpoint.Name,
		Pwd:         pwd,
		StartValid:  startTime,
		EndValid:    endTime,
		Notes:       notes,
		Blocked:     newEndpoint.Blocked,
		AutoRemove:  newEndpoint.AutoRemove,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("NewStream: failed to insert stream endpoint: %w", err))
	}

	endpointCreated := struct {
		EndpointID int `db:"endpointId"`
	}{
		EndpointID: endpoint.EndpointID,
	}

	return c.JSON(http.StatusCreated, endpointCreated)
}

// UpdateStream updates an existing position
// @Summary UpdateStream stream endpoint
// @ID update-stream
// @Tags stream-endpoints
// @Accept json
// @Param endpoint body stream.NewEditEndpoint true "Endpoint object"
// @Success 200
// @Router /v1/internal/streams/{endpointid} [put]
func (r *Repos) UpdateStream(c echo.Context) error {
	endpointID, err := strconv.Atoi(c.Param("endpointid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "UpdateStream: invalid endpoint ID")
	}

	_, err = r.stream.GetEndpointByID(c.Request().Context(), endpointID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "UpdateStream: endpoint not found")
	}

	var editEndpoint stream.NewEditEndpoint

	err = c.Bind(&editEndpoint)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("UpdateStream: failed to bind to request json: %w", err))
	}

	if len(editEndpoint.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("UpdateStream: endpoint name must be set"))
	}

	if len(editEndpoint.Application) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("UpdateStream: endpoint application must be set"))
	}

	startTime := null.NewTime(time.Time{}, false)

	if editEndpoint.StartValid != "" {
		var parseStart time.Time

		parseStart, err = time.Parse(time.RFC3339, editEndpoint.StartValid)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("UpdateStream: failed to parse start date: %w", err))
		}

		startTime = null.TimeFrom(parseStart)
	}

	endTime := null.NewTime(time.Time{}, false)

	if editEndpoint.EndValid != "" {
		var parseEnd time.Time

		parseEnd, err = time.Parse(time.RFC3339, editEndpoint.EndValid)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("UpdateStream: failed to parse end date: %w", err))
		}

		if startTime.Valid {
			diffStartEnd := startTime.Time.Compare(parseEnd)
			if diffStartEnd != -1 {
				return echo.NewHTTPError(http.StatusBadRequest, "UpdateStream: end date must be after start date")
			}
		}

		endTime = null.TimeFrom(parseEnd)
	}

	pwd := null.NewString("", false)
	notes := null.NewString("", false)

	if editEndpoint.Pwd != "" {
		pwd = null.StringFrom(editEndpoint.Pwd)
	}

	if editEndpoint.Notes != "" {
		notes = null.StringFrom(editEndpoint.Notes)
	}

	err = r.stream.EditEndpoint(c.Request().Context(), stream.EndpointDB{
		EndpointID:  endpointID,
		Application: editEndpoint.Application,
		Name:        editEndpoint.Name,
		Pwd:         pwd,
		StartValid:  startTime,
		EndValid:    endTime,
		Notes:       notes,
		Blocked:     editEndpoint.Blocked,
		AutoRemove:  editEndpoint.AutoRemove,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("UpdateStream: failed to update stream endpoint: %w", err))
	}

	return c.NoContent(http.StatusOK)
}

// DeleteStream deletes a stream endpoint
//
// @Summary Delete stream endpoint
// @ID delete-stream
// @Tags stream-endpoints
// @Accept json
// @Param endpointid path int true "Endpoint ID"
// @Success 200
// @Router /v1/internal/streams/{endpointid} [delete]
func (r *Repos) DeleteStream(c echo.Context) error {
	endpointID, err := strconv.Atoi(c.Param("endpointid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "DeleteStream: invalid endpoint ID")
	}

	err = r.stream.DeleteEndpoint(c.Request().Context(), endpointID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("DeleteStream: failed to delete stream endpoint: %w", err))
	}

	return c.NoContent(http.StatusOK)
}
