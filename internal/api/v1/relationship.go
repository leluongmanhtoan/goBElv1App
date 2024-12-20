package api

import (
	"errors"
	"net/http"
	"program/internal/middleware"
	"program/internal/response"
	"program/internal/services"
	"program/internal/validate"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Relationships struct {
	service services.IRelationshipsService
}

func NewRelationshipsAPI(engine *gin.Engine, relationshipsService services.IRelationshipsService) {
	handler := &Relationships{
		service: relationshipsService,
	}
	Group := engine.Group("api/v1")
	{
		Group.POST("follow", middleware.AuthMdw.RequestAuthorization(), handler.ToggleFollow)
		Group.GET(":id/followers", middleware.AuthMdw.RequestAuthorization(), handler.RetrieveFollowers)
		Group.GET(":id/following", middleware.AuthMdw.RequestAuthorization(), handler.RetrieveFollowing)
		Group.GET(":id/count-follow", middleware.AuthMdw.RequestAuthorization(), handler.RetrieveNumberOfFollowRelationship)
	}
}
func (h *Relationships) RetrieveNumberOfFollowRelationship(c *gin.Context) {
	//user_id, existed := c.Get("user_id")
	_, existed := c.Get("user_id")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user_id not found")))
		return
	}
	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(response.BadRequest(errors.New("id is empty")))
		return
	}

	//Note: Yeu cau thuc hien logic kiem tra xem user hien tai co duoc nguoi khac cho phep lay danh sach khong
	/*if user_id != id{

	}*/
	getResponse, err := h.service.GetFollowRelationshipCount(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, getResponse)
}

func (h *Relationships) RetrieveFollowers(c *gin.Context) {
	//user_id, existed := c.Get("user_id")
	_, existed := c.Get("user_id")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user_id not found")))
		return
	}
	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(response.BadRequest(errors.New("id is empty")))
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(response.BadRequest(errors.New("limit is a number")))
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(response.BadRequest(errors.New("offset is a number")))
		return
	}

	//Note: Yeu cau thuc hien logic kiem tra xem user hien tai co duoc nguoi khac cho phep lay danh sach khong
	/*if user_id != id{

	}*/
	getResponse, err := h.service.GetFollowers(c, limit, offset, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, getResponse)
}

func (h *Relationships) RetrieveFollowing(c *gin.Context) {
	//user_id, existed := c.Get("user_id")
	_, existed := c.Get("user_id")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user_id not found")))
		return
	}
	id := c.Param("id")
	if len(id) <= 0 {
		c.JSON(response.BadRequest(errors.New("id is empty")))
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(response.BadRequest(errors.New("limit is a number")))
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(response.BadRequest(errors.New("offset is a number")))
		return
	}

	//Note: Yeu cau thuc hien logic kiem tra xem user hien tai co duoc nguoi khac cho phep lay danh sach khong
	/*if user_id != id{

	}*/
	getResponse, err := h.service.GetFollowing(c, limit, offset, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, getResponse)
}

func (h *Relationships) ToggleFollow(c *gin.Context) {
	var request struct {
		FollowingUserID string `json:"followingId" validate:"required"`
	}
	if !validate.ValidateRequest(c, &request) {
		return
	}
	user_id, existed := c.Get("user_id")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user_id not found")))
		return
	}
	insertResponse, err := h.service.ToggleFollow(c, user_id.(string), request.FollowingUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, insertResponse)
}
