package api

import (
	"errors"
	"net/http"
	"program/internal/middleware"
	"program/internal/model"
	"program/internal/response"
	"program/internal/services"
	"program/internal/validate"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Newsfeed struct {
	service services.INewsfeedService
}

func NewNewsFeedAPI(engine *gin.Engine, service services.INewsfeedService) {
	handler := &Newsfeed{
		service: service,
	}
	Group := engine.Group("api/v1")
	{
		//newsfeed
		Group.POST("posts", middleware.AuthMdw.RequestAuthorization(), handler.PostNewsFeed)
		Group.GET("user/:id/posts", middleware.AuthMdw.RequestAuthorization())
		Group.GET("user/following/posts", middleware.AuthMdw.RequestAuthorization(), handler.RetrieveNewsfeed)

		//interact newsfeed
		Group.POST("posts/:postId/like", middleware.AuthMdw.RequestAuthorization(), handler.ToggleLikePost)
		Group.GET("posts/:postId/like", handler.RetrieveLikers)

		Group.POST("posts/:postId/comment", middleware.AuthMdw.RequestAuthorization(), handler.PostComment)
		Group.GET("posts/:postId/comments", middleware.AuthMdw.RequestAuthorization(), handler.RetrieveComments)
		Group.PUT("posts/:postId/comment", middleware.AuthMdw.RequestAuthorization(), handler.PutComment)
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
	user_id, existed := c.Get("user_id")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user_id not found")))
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
	getResponse, err := h.service.GetNewsfeed(c, limit, offset, user_id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, getResponse)

}

func (h *Newsfeed) ToggleLikePost(c *gin.Context) {
	user_id, existed := c.Get("user_id")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user_id not found")))
		return
	}
	id := c.Param("postId")
	if len(id) <= 0 {
		c.JSON(response.BadRequest(errors.New("id is empty")))
		return
	}
	toggleResponse, err := h.service.ToggleLikePost(c, user_id.(string), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, toggleResponse)
}

func (h *Newsfeed) RetrieveLikers(c *gin.Context) {
	//user_id, existed := c.Get("user_id")
	id := c.Param("postId")
	if len(id) <= 0 {
		c.JSON(response.BadRequest(errors.New("id is empty")))
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
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
	getResponse, err := h.service.GetLikers(c, limit, offset, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, getResponse)
}

func (h *Newsfeed) PostComment(c *gin.Context) {
	newcomment := new(model.CommentPost)
	if !validate.ValidateRequest(c, newcomment) {
		return
	}
	user_id, existed := c.Get("user_id")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user_id not found")))
		return
	}
	id := c.Param("postId")
	if len(id) <= 0 {
		c.JSON(response.BadRequest(errors.New("id is empty")))
		return
	}
	insertResponse, err := h.service.PostComment(c, user_id.(string), id, newcomment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, insertResponse)
}

func (h *Newsfeed) RetrieveComments(c *gin.Context) {
	id := c.Param("postId")
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
	getResponse, err := h.service.GetComments(c, limit, offset, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, getResponse)

}

func (h *Newsfeed) PutComment(c *gin.Context) {
	commentPut := new(model.CommentPut)
	if !validate.ValidateRequest(c, commentPut) {
		return
	}
	putResponse, err := h.service.PutComment(c, commentPut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, putResponse)
}
