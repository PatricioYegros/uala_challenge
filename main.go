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
)

var ErrPermission = "User dont have permission to perform action"

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

	r.POST("/user/login/:userID", login)
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

	log, err := checkUserLog(uint(followerID))
	if err != nil {
		returnError(c, err)
		return
	}

	if log {

		err = twitterService.Follow(uint(followerID), uint(userID))
		if err != nil {
			returnError(c, err)
			return
		}

		c.String(http.StatusNoContent, "")

	} else {

		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrPermission})
	}
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

	log, err := checkUserLog(uint(userID))
	if err != nil {
		returnError(c, err)
		return
	}

	if log {

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

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrPermission})
	}
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

	log, err := checkUserLog(uint(userID))
	if err != nil {
		returnError(c, err)
		return
	}

	if log {

		timeline, err := twitterService.GetTimeLine(uint(userID))
		if err != nil {
			returnError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, timeline)

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrPermission})
	}
}

func returnError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, err.Error())
}

// @Summary Login
// @Description Logs
// @Tags Twtter
// @Param userID path uint true "userID"
// @Produce application/json
// @Success 200
// @Router /user/login/:userID [post]
func login(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		returnError(c, err)
		return
	}

	err = twitterService.Login(uint(userID))
	if err != nil {
		returnError(c, err)
		return
	}

	c.String(http.StatusNoContent, "")
}

func checkUserLog(userID uint) (bool, error) {
	value, err := twitterService.CheckUserLog(userID)
	if err != nil {
		return false, err
	}

	return value, nil
}
