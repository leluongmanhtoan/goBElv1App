package api

import (
	"errors"
	"net/http"
	"program/internal/response"
	"program/internal/validate"
	"program/middleware"
	"program/service"

	"github.com/gin-gonic/gin"
)

type Relationships struct {
	service service.IRelationships
}

func NewRelationshipsAPI(engine *gin.Engine, relationshipsService service.IRelationships) {
	handler := &Relationships{
		service: relationshipsService,
	}
	Group := engine.Group("api/v1")
	{
		Group.POST("follow", middleware.AuthMdw.RequestAuthorization(), handler.FollowUser)
		Group.DELETE("follow", middleware.AuthMdw.RequestAuthorization(), handler.UnFollowUser)
		Group.GET(":id/followers", middleware.AuthMdw.RequestAuthorization(), handler.RetrieveFollowers)
		Group.GET(":id/following", middleware.AuthMdw.RequestAuthorization(), handler.RetrieveFollowing)
	}
}

func (h *Relationships) RetrieveFollowers(c *gin.Context) {
}

func (h *Relationships) RetrieveFollowing(c *gin.Context) {

}

func (h *Relationships) FollowUser(c *gin.Context) {
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
	insertResponse, err := h.service.FollowUser(c, user_id.(string), request.FollowingUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, insertResponse)
}

func (h *Relationships) UnFollowUser(c *gin.Context) {
	var request struct {
		FollowingUserID string `json:"unfollowingId" validate:"required"`
	}
	if !validate.ValidateRequest(c, &request) {
		return
	}
	user_id, existed := c.Get("user_id")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user_id not found")))
		return
	}
	deleteResponse, err := h.service.UnFollowUser(c, user_id.(string), request.FollowingUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, deleteResponse)
}
