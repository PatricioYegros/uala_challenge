// Ualá Challenge Patricio Yegros
//
// This is the swagger documentation
//
//	Schemes: http
//	Host: localhost:8080
//	BasePath: /
//	Version: 1.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- text/plain
//	- application/json
//
// swagger:meta
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PatricioYegros/uala_challenge/app"
	"github.com/PatricioYegros/uala_challenge/app/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:generate swagger generate spec -o ./swagger.json
var twitterService *service.TwitterService

func init() {
	var err error

	twitterService, _, err = app.NewService()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.POST("/user/:userID/tweet", tweet)
	r.POST("/user/:userID/follower/:followerID", follow)
	r.GET("/user/:userID/timeline", timeline)

	log.Fatalln(r.Run())
}

// @Summary Follow User
// @Description FollowerID start to follow UserID
// @Tags Twitter
// @Param followerID path uint true "followerID"
// @Param userID path uint true "userID"
// @Produce text/plain
// @Success 200
// @Router /user/{userID}/follower/{followerID} [post]
func follow(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		returnError(c, err)
		return
	}

	followerID, err := strconv.Atoi(c.Param("followerID"))
	if err != nil {
		returnError(c, err)
		return
	}

	err = twitterService.Follow(uint(followerID), uint(userID))
	if err != nil {
		returnError(c, err)
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%d has followed %d", followerID, userID))
}

type TweetRequestBody struct {
	Body string `json:"content"`
}

// @Summary Tweet
// @Description User makes a Tweet
// @Tags Twitter
// @Param userID path uint true "userID"
// @Param body body string true "body"
// @Produce text/plain
// @Success 201
// @Router /user/{userID}/tweet [post]
func tweet(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		returnError(c, err)
		return
	}

	var requestBody TweetRequestBody

	if err = c.BindJSON(&requestBody); err != nil {
		returnError(c, err)
		return
	}

	tweetID, err := twitterService.Tweet(uint(userID), requestBody.Body)
	if err != nil {
		returnError(c, err)
		return
	}

	c.String(http.StatusCreated, fmt.Sprintf("%d tweet %s created", userID, tweetID))
}

// @Summary Timeline
// @Description Get the timeline of certain user
// @Tags Twitter
// @Param userID path uint true "userID"
// @Produce application/json
// @Success 200
// @Router /user/{userID}/timeline [get]
func timeline(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		returnError(c, err)
		return
	}

	timeline, err := twitterService.GetTimeLine(uint(userID))
	if err != nil {
		returnError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, timeline)
}

func returnError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, err.Error())
}
