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
	"github.com/google/uuid"
)

type Newsfeed struct {
	service services.INewsfeedService
}

func NewNewsFeedAPI(engine *gin.Engine, service services.INewsfeedService) {
	handler := &Newsfeed{
		service: service,
	}
	Group := engine.Group("api/v1/newsfeed")
	{
		//newsfeed
		Group.GET("", middleware.AuthMdw.RequestAuthorization(), handler.GetNewsfeed)
		Group.POST("post", middleware.AuthMdw.RequestAuthorization(), handler.CreatePost)
		Group.PATCH("post/:id", middleware.AuthMdw.RequestAuthorization())
		Group.DELETE("post/:id", middleware.AuthMdw.RequestAuthorization())

		//Group.GET("user/:id/posts", middleware.AuthMdw.RequestAuthorization())

		//interact newsfeed
		Group.POST("post/:postId/like", middleware.AuthMdw.RequestAuthorization(), handler.ToggleLikePost)
		Group.GET("post/:postId/like", middleware.AuthMdw.RequestNoRequiredAuthorization(), handler.GetLikers)

		Group.POST("post/:postId/comment", middleware.AuthMdw.RequestAuthorization(), handler.PostComment)
		Group.PUT("post/:postId/comment", middleware.AuthMdw.RequestAuthorization(), handler.PutComment)
		Group.GET("post/:postId/comments", middleware.AuthMdw.RequestAuthorization(), handler.RetrieveComments)

	}
}

func (h *Newsfeed) CreatePost(c *gin.Context) {
	newpost := new(model.NewsfeedPost)
	if !validate.ValidateRequest(c, newpost) {
		return
	}
	userId, existed := c.Get("userId")
	if !existed || userId == "" {
		response.ErrorResponse[string](c, http.StatusBadRequest, "user id not found")
		return
	}
	mypost, err := h.service.CreatePost(c, userId.(string), newpost)
	if err != nil {
		response.ErrorResponse[string](c, http.StatusInternalServerError, "can not create a new post")
		return
	}

	response.SuccessResponse(c, "create post successfully", mypost)
}

func (h *Newsfeed) GetNewsfeed(c *gin.Context) {
	userId, existed := c.Get("userId")
	if !existed || userId == "" {
		response.ErrorResponse[string](c, http.StatusBadRequest, "user id not found")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		response.ErrorResponse[string](c, http.StatusBadRequest, "limit is a number")
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		response.ErrorResponse[string](c, http.StatusBadRequest, "offset is a number")
		return
	}
	newsfeed, err := h.service.GetNewsfeed(c, limit, offset, userId.(string))
	if err != nil {
		response.ErrorResponse[string](c, http.StatusInternalServerError, "can not get newsfeed")
		return
	}
	response.SuccessResponseWithPagination(c, limit, offset, userId.(string), newsfeed)

}

func (h *Newsfeed) ToggleLikePost(c *gin.Context) {
	userId, existed := c.Get("userId")
	if !existed || userId == "" {
		response.ErrorResponse[string](c, http.StatusBadRequest, "user id not found")
		return
	}
	postId := c.Param("postId")
	if postId == "" {
		response.ErrorResponse[string](c, http.StatusBadRequest, "post id can not be empty")
		return
	}
	if _, err := uuid.Parse(postId); err != nil {
		response.ErrorResponse[string](c, http.StatusBadRequest, "post id is not a valid UUID")
		return
	}
	err := h.service.ToggleLikePost(c, userId.(string), postId)
	if err != nil {
		response.ErrorResponse[string](c, http.StatusInternalServerError, "can not toggle like or unlike")
		return
	}
	response.SuccessResponse(c, "toggle like post successfully", "")
}

func (h *Newsfeed) GetLikers(c *gin.Context) {
	userId, existed := c.Get("userId")
	if !existed || userId == "" {
		userId = "guest"
	}
	postId := c.Param("postId")
	if postId == "" {
		response.ErrorResponse[string](c, http.StatusBadRequest, "post id can not be empty")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		response.ErrorResponse[string](c, http.StatusBadRequest, "limit is a number")
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		response.ErrorResponse[string](c, http.StatusBadRequest, "offset is a number")
		return
	}

	if userId == "guest" {
		likers, err := h.service.GetLikers(c, limit, offset, userId.(string), postId, true)
		if err != nil {
			response.ErrorResponse[string](c, http.StatusInternalServerError, err.Error())
			return
		}
		response.SuccessResponseWithPagination(c, limit, offset, userId.(string), likers)
	} else {
		likers, err := h.service.GetLikers(c, limit, offset, userId.(string), postId, false)
		if err != nil {
			response.ErrorResponse[string](c, http.StatusInternalServerError, err.Error())
			return
		}
		response.SuccessResponseWithPagination(c, limit, offset, userId.(string), likers)
	}
}

func (h *Newsfeed) PostComment(c *gin.Context) {
	newcomment := new(model.CommentPost)
	if !validate.ValidateRequest(c, newcomment) {
		return
	}
	user_id, existed := c.Get("userId")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user id not found")))
		return
	}
	id := c.Param("postId")
	if len(id) <= 0 {
		c.JSON(response.BadRequest(errors.New("id is empty")))
		return
	}
	mycomment, err := h.service.PostComment(c, user_id.(string), id, newcomment)
	if err != nil {
		response.ErrorResponse[string](c, http.StatusInternalServerError, "can not create a new comment")
		return
	}
	response.SuccessResponse(c, "create post successfully", mycomment)
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
