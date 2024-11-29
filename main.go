// Ualá Challenge - Patricio Yegros
//
// # This is the swagger documentation
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
	swaggerFiles "github.com/swaggo/files"
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
	url := ginSwagger.URL("http://localhost:3000/swagger/doc.json")

	r.POST("/user/:userID/tweet", tweet)
	r.POST("/user/:userID/follower/:followerID", follow)
	r.GET("/user/:userID/timeline", timeline)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	log.Fatalln(r.Run())
}

// swagger:operation POST /user/{userID}/follower/{followerID}
// Follow user
// @Summary Given 2 ids, the followerID begins to follow the userID
// @Tags Follow
// @Produce text/plain
//
// parameters:
//   - name: userID
//     type: uint
//     in: path
//     description: ID of user to be followed
//     required: true
//   - name: followerID
//     type: uint
//     in: path
//     description: ID of the follower
//     required: true
//
// responses:
// '200':
//
//	description: 'Followed correctly'
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

// swagger:operation POST /user/{userID}/tweet
// Tweet
// @Summary Create a tweet
// @Tags Tweet
// @Produce text/plain
//
// parameters:
//   - name: userID
//     type: uint
//     in: path
//     description: ID of owner of the tweet
//     required: true
//   - name: body
//     type: string
//     in: body
//     description: Body of the tweet
//     required: true
//
// responses:
// '201':
//
//	description: 'Tweet created correctly'
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

// swagger:operation GET /user/{userID}/timeline Timeline
// Get TimeLine
// @Summary Return the userID's timeline
// @Tags Timeline
// @Produce application/json
//
// parameters:
//   - name: userID
//     type: uint
//     in: path
//     description: ID of owner of timeline
//     required: true
//
// responses:
// '200':
//
//	description: 'Timeline obtained correctly'.
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
