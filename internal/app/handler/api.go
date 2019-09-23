package handler

import (
	"net/url"

	"strvucks-go/internal/app/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// API handles API
type API struct {
	AuthClient AuthClient
}

// NewAPI returns implemented API
func NewAPI() *API {
	return &API{&AuthClientImpl{}}
}

// AuthClient handles auth
type AuthClient interface {
	getAuthUserID(c *gin.Context) (int64, error)
}

// AuthClientImpl handles auth
type AuthClientImpl struct{}

func (w *AuthClientImpl) getAuthUserID(c *gin.Context) (int64, error) {
	return GetAuthUserID(c.Request)
}

func (w *API) internalError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}

// StravaAuthURL returns a Strava oauth2 URL
func (w *API) StravaAuthURL(c *gin.Context) {
	config := Config()
	authURL, _ := url.QueryUnescape(config.AuthCodeURL("strvucks", AuthCodeOption()...))

	c.JSON(200, map[string]interface{}{
		"url": authURL,
	})
}

func (w *API) getUser(c *gin.Context) (*model.User, bool) {
	id, err := w.AuthClient.getAuthUserID(c)
	if err != nil {
		log.Error("Not auth", err)
		c.JSON(401, nil)
		return nil, false
	}

	user := &model.User{ID: id}
	if err := model.DB().First(user).Error; err != nil {
		log.Error("Not found user:", id, err)
		c.JSON(404, nil)
		return nil, false
	}

	return user, true
}

// CurrentUserHandler handles current user
func (w *API) CurrentUserHandler(c *gin.Context) {
	user, ok := w.getUser(c)
	if !ok {
		return
	}

	c.JSON(200, user)
}

// UpdateCurrentUserHandler handles current user
func (w *API) UpdateCurrentUserHandler(c *gin.Context) {
	user, ok := w.getUser(c)
	if !ok {
		return
	}

	userParams := model.User{}
	if err := c.BindJSON(&userParams); err != nil {
		log.Error("Invalid params", err)
		c.JSON(400, nil)
		return
	}

	user.IftttKey = userParams.IftttKey
	user.IftttMessage = userParams.IftttMessage
	if err := user.Save(model.DB()).Error; err != nil {
		log.Error("Failure save user", err)
		c.JSON(500, w.internalError(err))
		return
	}

	c.JSON(200, user)
}

// MySummaryHandler handles current user
func (w *API) MySummaryHandler(c *gin.Context) {
	user, ok := w.getUser(c)
	if !ok {
		return
	}

	summary := &model.Summary{}
	if err := summary.FirstOrInit(model.DB(), user.AthleteID).Error; err != nil {
		log.Error("Failure get summary", err)
		c.JSON(500, w.internalError(err))
		return
	}

	c.JSON(200, summary)
}
