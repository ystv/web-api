package stream

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/ystv/web-api/services/stream"
)

// Repos encapsulates the dependency
type Repos struct {
	stream *stream.Store
}

// NewRepos creates our data store
func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{stream.NewStore(db)}
}

type SRSPublish struct {
	Action      string `json:"action"`
	IP          string `json:"ip"`
	VHost       string `json:"vhost"`
	Application string `json:"app"`
	Url         string `json:"tcUrl"`
	Stream      string `json:"stream"`
	Param       string `json:"param"`
}

func handleSRSPublish(c echo.Context) (string, string, string, string, error) {
	var publish SRSPublish
	var application, name, pwd, action string
	var err error

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(c.Request().Body)
	dec := json.NewDecoder(c.Request().Body)
	err = dec.Decode(&publish)
	if err != nil {
		return "", "", "", "", err
	}

	// skip question mark
	if len(publish.Param) > 0 {
		publish.Param = publish.Param[1:]
	}

	val, err := url.ParseQuery(publish.Param)
	if err != nil {
		return "", "", "", "", err
	}
	application = publish.Application
	name = publish.Stream
	pwd = val.Get("pwd")
	action = publish.Action

	return application, pwd, name, action, nil
}

func handleNginxPublish(c echo.Context) (string, string, string, string) {
	var application, name, pwd, action string
	application = c.FormValue("app")
	name = c.FormValue("name")
	pwd = c.FormValue("pwd")
	action = c.FormValue("call")

	return application, pwd, name, action
}

func (r *Repos) PublishHandler(c echo.Context) error {
	var application, name, pwd, action string
	var err error

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(c.Request().Body)

	if c.Request().Header.Get("Content-Type") == "application/json" {
		// SRS publish handler
		application, name, pwd, action, err = handleSRSPublish(c)
		if action != "on_publish" {
			err = errors.New("publish stream invalid action " + action)
		}
	} else {
		// Form DATA from nginx-rtmp/srtrelay
		application, name, pwd, action = handleNginxPublish(c)

		// only apply auth for publish
		if action != "publish" {
			return nil
		}
	}

	if err != nil {
		log.Printf("failed to parse publish data: %+v", err)
		return c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	endpoint, err := r.stream.GetEndpointByApplicationNamePwd(c.Request().Context(), application, name, pwd)
	if err != nil {
		log.Printf("failed to get endpoint: %+v", err)
		return c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	if endpoint.Active || endpoint.Blocked {
		log.Printf("endpoint active (%t) or blocked (%t)", endpoint.Active, endpoint.Blocked)
		return c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	err = r.stream.SetEndpointActiveByID(c.Request().Context(), endpoint.EndpointID)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("publish stream %s %s/%s ok", endpoint.EndpointID, application, name)
	}

	// SRS needs zero response
	return c.String(http.StatusOK, "0")
}

func (r *Repos) UnpublishHandler(c echo.Context) error {
	var application, name, pwd, action string
	var err error

	if c.Request().Header.Get("Content-Type") == "application/json" {
		// SRS publish handler
		application, name, _, action, err = handleSRSPublish(c)
		if action != "on_unpublish" {
			err = fmt.Errorf("unpublish stream invalid action %s", action)
		}
	} else {
		// Form DATA from nginx-rtmp/srtrelay
		application, name, pwd, action = handleNginxPublish(c)
		log.Println("unpublish action", action)
		// ignore actions except unpublish
		if action != "publish_done" {
			return nil
		}
	}

	if err != nil {
		log.Printf("failed to parse unpublish data: %+v", err)
		return c.String(http.StatusUnauthorized, "401 Unauthorized")
	}
	err = r.stream.SetEndpointInactiveByApplicationNamePwd(c.Request().Context(), application, name, pwd)
	if err != nil {
		log.Printf("failed to unpublish stream, continuing, %s/%s: %+v", application, name, err)
	} else {
		log.Printf("unpublish stream %s/%s ok", application, name)
	}

	// SRS needs zero response
	return c.String(http.StatusOK, "0")
}
