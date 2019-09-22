package handler

import (
	"net/url"

	"strvucks-go/internal/app/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// API handles API
type API struct {
	WebhookClient WebhookClient
}

// StravaAuthURL returns a Strava oauth2 URL
func (w *API) StravaAuthURL(c *gin.Context) {
	config := Config()
	authURL, _ := url.QueryUnescape(config.AuthCodeURL("strvucks", AuthCodeOption()...))

	c.JSON(200, map[string]interface{}{
		"url": authURL,
	})
}

// CurrentUserHandler handles current user
func (w *API) CurrentUserHandler(c *gin.Context) {
	id, err := GetAuthUserID(c.Request)
	if err != nil {
		c.JSON(401, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	user := &model.User{ID: id}
	if err := model.DB().First(user).Error; err != nil {
		c.JSON(404, nil)
		return
	}

	c.JSON(200, user)
}

// UpdateCurrentUserHandler handles current user
func (w *API) UpdateCurrentUserHandler(c *gin.Context) {
	id, err := GetAuthUserID(c.Request)
	if err != nil {
		c.JSON(401, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	user := &model.User{ID: id}
	if err := model.DB().First(user).Error; err != nil {
		c.JSON(404, nil)
		return
	}

	userParams := model.User{}
	if err := c.BindJSON(&userParams); err != nil {
		c.JSON(400, nil)
		return
	}

	user.IftttKey = userParams.IftttKey
	user.IftttMessage = userParams.IftttMessage
	if err := user.Save(model.DB()).Error; err != nil {
		log.Error("Failure save user", err)
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, user)
}
