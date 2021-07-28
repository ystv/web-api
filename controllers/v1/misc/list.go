package misc

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ystv/web-api/utils"
)

// GetLists handles listing mailing lists
//
// @Summary Get Mailing lists
// @Description Lists all mailing lists.
// @ID get-mailing-lists
// @Tags misc-list
// @Produce json
// @Success 200 {array} misc.List
// @Router /v1/internal/misc/lists [get]
func (r *Repos) GetLists(c echo.Context) error {
	l, err := r.misc.GetLists(c.Request().Context())
	if err != nil {
		err = fmt.Errorf("GetLists failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, l)
}

// GetListsByToken handles listing mailing lists
// enables IsSubscribed check
//
// @Summary Get Mailing lists by token
// @Description Lists all mailing lists. Provides extra context for what the user has subscribed too
// @ID get-mailing-lists-token
// @Tags misc-list
// @Produce json
// @Success 200 {array} misc.List
// @Router /v1/internal/misc/lists/my [get]
func (r *Repos) GetListsByToken(c echo.Context) error {
	p, err := utils.GetToken(c.Request().Response.Request)
	if err != nil {
		err = fmt.Errorf("SetCrew: failed to get token: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	l, err := r.misc.GetListsByUserID(c.Request().Context(), p.UserID)
	if err != nil {
		err = fmt.Errorf("GetListsByToken failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, l)
}

// GetList handles listing a mailing list including the subscribers
//
// @Summary Get Mailing list
// @Description Get a mailing list. Provides list subscribers also
// @ID get-mailing-list-id
// @Tags misc-list
// @Produce json
// @Param listid path int true "List ID"
// @Success 200 {object} misc.List
// @Router /v1/internal/misc/list/{listid} [get]
func (r *Repos) GetList(c echo.Context) error {
	listID, err := strconv.Atoi(c.Param("listid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad listID")
	}
	l, err := r.misc.GetList(c.Request().Context(), listID)
	if err != nil {
		err = fmt.Errorf("GetList failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, l)
}

// GetSubscribers handles listing a mailing list's subscribers
//
// @Summary Get subscribers
// @Description Get a mailing list's subscribers
// @ID get-mailing-list-subscribers-id
// @Tags misc-list
// @Produce json
// @Param listid path int true "List ID"
// @Success 200 {object} misc.List
// @Router /v1/internal/misc/list/{listid}/subscribers [get]
func (r *Repos) GetSubscribers(c echo.Context) error {
	listID, err := strconv.Atoi(c.Param("listid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad listID")
	}
	s, err := r.misc.GetSubscribers(c.Request().Context(), listID)
	if err != nil {
		err = fmt.Errorf("GetSubscribers failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, s)
}

// SubscribeByToken handles subscribing a user to a mailing list
//
// @Summary Subscribe to mailing list by token
// @Description Subscribe to a mailing list by a JWT
// @ID new-mailing-list-subscriber-token
// @Tags misc-list
// @Accept json
// @Param listid path int true "List ID"
// @Success 201
// @Router /v1/internal/misc/list/{listid}/subscribe [post]
func (r *Repos) SubscribeByToken(c echo.Context) error {
	listID, err := strconv.Atoi(c.Param("listid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad listID")
	}
	claims, err := utils.GetToken(c.Request().Response.Request)
	if err != nil {
		err = fmt.Errorf("SubscribeByToken failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = r.misc.Subscribe(c.Request().Context(), claims.UserID, listID)
	if err != nil {
		err = fmt.Errorf("SubscribeByToken failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

// SubscribeByID handles subscribing a user to a mailing list
//
// @Summary Subscribe to mailing list by user ID
// @Description Subscribe to a mailing list by a user ID
// @ID new-mailing-list-subscriber-id
// @Tags misc-list
// @Accept json
// @Param listid path int true "List ID"
// @Param userid path int true "User ID"
// @Success 201
// @Router /v1/internal/misc/list/{listid}/subscribe/{userid} [post]
func (r *Repos) SubscribeByID(c echo.Context) error {
	listID, err := strconv.Atoi(c.Param("listid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad listID")
	}
	userID, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad userID")
	}
	err = r.misc.Subscribe(c.Request().Context(), userID, listID)
	if err != nil {
		err = fmt.Errorf("SubscribeByToken failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

// UnsubscribeByToken handles unsubscribing a user to a mailing list
//
// @Summary Unsubscribe to mailing list by token
// @Description Unsubscribe to a mailing list by a JWT
// @ID delete-mailing-list-subscriber-token
// @Tags misc-list
// @Accept json
// @Param listid path int true "List ID"
// @Success 200
// @Router /v1/internal/misc/list/{listid}/unsubscribe [delete]
func (r *Repos) UnsubscribeByToken(c echo.Context) error {
	listID, err := strconv.Atoi(c.Param("listid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad listID")
	}
	claims, err := utils.GetToken(c.Request().Response.Request)
	if err != nil {
		err = fmt.Errorf("UnsubscribeByToken failed to get user ID: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = r.misc.UnsubscribeByID(c.Request().Context(), claims.UserID, listID)
	if err != nil {
		err = fmt.Errorf("UnsubscribeByToken failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// UnsubscribeByID handles subscribing a user to a mailing list
//
// @Summary Unsubscribe to mailing list by user ID
// @Description Unsubscribe to a mailing list by a user ID
// @ID delete-mailing-list-subscriber-id
// @Tags misc-list
// @Accept json
// @Param listid path int true "List ID"
// @Param userid path int true "User ID"
// @Success 200
// @Router /v1/internal/misc/list/{listid}/unsubscribe/{userid} [delete]
func (r *Repos) UnsubscribeByID(c echo.Context) error {
	listID, err := strconv.Atoi(c.Param("listid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad listID")
	}
	userID, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad userID")
	}
	err = r.misc.UnsubscribeByID(c.Request().Context(), userID, listID)
	if err != nil {
		err = fmt.Errorf("UnsubscribeByToken failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// UnsubscribeByUUID handles subscribing a user to a mailing list
//
// @Summary Unsubscribe to mailing list by subscriber UUID
// @Description Unsubscribe to a mailing list by a subscriber UUID
// @ID delete-mailing-list-subscriber-uuid
// @Tags misc-list
// @Accept json
// @Param uuid path int true "Subscriber UUID"
// @Success 200
// @Router /v1/list_unsubscribe/{uuid} [get]
func (r *Repos) UnsubscribeByUUID(c echo.Context) error {
	err := r.misc.UnsubscribeByUUID(c.Request().Context(), c.Param("uuid"))
	if err != nil {
		err = fmt.Errorf("UnsubscribeByUUID failed: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
