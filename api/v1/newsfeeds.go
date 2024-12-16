package api

import "github.com/gin-gonic/gin"

type Newsfeed struct {
}

func NewNewsFeed(engine *gin.Engine) {
	Group := engine.Group("api/v1")
	{
		Group.POST("")
	}
}
