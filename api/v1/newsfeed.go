package api

import (
	"errors"
	"net/http"
	"program/internal/response"
	"program/internal/validate"
	"program/middleware"
	"program/model"
	"program/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Newsfeed struct {
	service service.INewsfeed
}

func NewNewsFeedAPI(engine *gin.Engine, service service.INewsfeed) {
	handler := &Newsfeed{
		service: service,
	}
	Group := engine.Group("api/v1")
	{
		Group.POST("posts", middleware.AuthMdw.RequestAuthorization(), handler.PostNewsFeed)
		Group.GET("user/:id/posts", middleware.AuthMdw.RequestAuthorization())
		Group.GET("user/:id/following/posts", middleware.AuthMdw.RequestAuthorization(), handler.RetrieveNewsfeed)
	}
}

func (h *Newsfeed) PostNewsFeed(c *gin.Context) {
	newpost := new(model.NewsfeedPost)
	if !validate.ValidateRequest(c, newpost) {
		return
	}
	user_id, existed := c.Get("user_id")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user_id not found")))
		return
	}
	insertResponse, err := h.service.PostNewsfeed(c, user_id.(string), newpost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, insertResponse)
}

func (h *Newsfeed) RetrieveNewsfeed(c *gin.Context) {
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
	getResponse, err := h.service.GetNewsfeed(c, limit, offset, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, getResponse)

}
