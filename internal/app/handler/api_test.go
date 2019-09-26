package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"strvucks-go/internal/app/model"

  "github.com/miyanokomiya/strava-client-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// AuthClientMock handles auth
type AuthClientMock struct{}

func (w *AuthClientMock) getAuthUserID(c *gin.Context) (int64, error) {
	return 10, nil
}

// AuthClientMock handles auth
type AuthClientMockError struct{}

func (w *AuthClientMockError) getAuthUserID(c *gin.Context) (int64, error) {
	return 0, errors.New("error")
}

func TestInternalError(t *testing.T) {
	api := API{}
	assert.Equal(t, map[string]interface{}{
		"error": "hoge",
	}, api.internalError(errors.New("hoge")), "returns error hash")
}

func TestGetUser(t *testing.T) {
	user := model.User{ID: 10}
	db := model.DB()
	if err := db.Save(&user).Error; err != nil {
		t.Fatal("cannot create user", err)
	}
	defer db.Delete(&user)

	{
		router := gin.New()
		router.GET("/hoge", func(c *gin.Context) {
			api := API{nil, &AuthClientMock{}}
			user, ok := api.getUser(c)
			assert.True(t, ok, "success get user")
			assert.Equal(t, int64(10), user.ID, "success get user")
			c.JSON(200, user)
		})

		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code, "success get user")
	}

	{
		router := gin.New()
		router.GET("/hoge", func(c *gin.Context) {
			api := API{nil, &AuthClientMockError{}}
			user, ok := api.getUser(c)
			assert.False(t, ok, "not auth")
			assert.Nil(t, user, "not auth")
		})

		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code, "not auth")
	}

	{
		db.Delete(&user)
		router := gin.New()
		router.GET("/hoge", func(c *gin.Context) {
			api := API{nil, &AuthClientMock{}}
			user, ok := api.getUser(c)
			assert.False(t, ok, "not found")
			assert.Nil(t, user, "not found")
		})

		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 404, w.Code, "not found")
	}
}

func TestGetPermission(t *testing.T) {
	user := model.User{ID: 10, AthleteID: 100}
	db := model.DB()
	if err := db.Save(&user).Error; err != nil {
		t.Fatal("cannot create user", err)
	}
	defer db.Delete(&user)

	permission := model.Permission{AthleteID: 100}
	if err := db.Save(&permission).Error; err != nil {
		t.Fatal("cannot create permission", err)
	}
	defer db.Delete(&permission)

	{
		router := gin.New()
		router.GET("/hoge", func(c *gin.Context) {
			api := API{}
			p, ok := api.getPermission(c, &user)
			assert.True(t, ok, "success get permission")
			assert.Equal(t, int64(100), p.AthleteID, "success get permission")
			c.JSON(200, p)
		})

		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code, "success get permission")
	}

	{
		db.Delete(&user)
		router := gin.New()
		router.GET("/hoge", func(c *gin.Context) {
			api := API{}
			p, ok := api.getPermission(c, &user)
			assert.False(t, ok, "not auth")
			assert.Nil(t, p, "not auth")
		})

		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 401, w.Code, "not auth")
	}
}

func TestStravaAuthURL(t *testing.T) {
	api := API{}
	router := gin.New()
	router.GET("/hoge", api.StravaAuthURL)

	{
		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatal("cannot get URL", w.Code, w.Body)
		} else {
			decoder := json.NewDecoder(w.Body)
			data := map[string]interface{}{}
			if err := decoder.Decode(&data); err != nil {
				t.Fatal("cannot get URL", err)
			}
			url, ok := data["url"].(string)
			if !ok {
				t.Fatal("cannot get URL")
			}
			assert.NotEqual(t, "", url, "success get URL")
		}
	}
}

func TestCurrentUserHandler(t *testing.T) {
	godotenv.Load("../../../.env")

	api := API{nil, &AuthClientMock{}}
	router := gin.New()
	router.GET("/hoge", api.CurrentUserHandler)

	user := model.User{ID: 10}
	db := model.DB()
	if err := db.Save(&user).Error; err != nil {
		t.Fatal("cannot create user", err)
	}
	defer db.Delete(&user)

	{
		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatal("cannot get user", w.Code, w.Body)
		} else {
			decoder := json.NewDecoder(w.Body)
			user := model.User{}
			if err := decoder.Decode(&user); err != nil {
				t.Fatal("cannot get user", err)
			}
			assert.Equal(t, int64(10), user.ID, "success get user")
		}
	}
}

func TestUpdateCurrentUserHandler(t *testing.T) {
	godotenv.Load("../../../.env")

	api := API{nil, &AuthClientMock{}}
	router := gin.New()
	router.POST("/hoge", api.UpdateCurrentUserHandler)

	user := model.User{ID: 10}
	db := model.DB()
	if err := db.Save(&user).Error; err != nil {
		t.Fatal("cannot create user", err)
	}
	defer db.Delete(&user)

	{
		var jsonStr = []byte(`{"iftttKey":"key", "iftttMessage":"message"}`)
		req, err := http.NewRequest("POST", "/hoge", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatal("cannot update user", w.Code, w.Body)
		} else {
			user := model.User{}
			if err := db.First(&user, model.User{ID: 10}).Error; err != nil {
				t.Fatal("cannot find user", err)
			}
			assert.Equal(t, "key", user.IftttKey, "update IftttKey")
			assert.Equal(t, "message", user.IftttMessage, "update IftttMessage")
		}
	}
}

func TestMySummaryHandler(t *testing.T) {
	godotenv.Load("../../../.env")

	api := API{nil, &AuthClientMock{}}
	router := gin.New()
	router.GET("/hoge", api.MySummaryHandler)

	user := model.User{ID: 10, AthleteID: 20}
	db := model.DB()
	if err := db.Save(&user).Error; err != nil {
		t.Fatal("cannot create user", err)
	}
	defer db.Delete(&user)

	summary := model.Summary{AthleteID: 20, WeeklyCount: 3}
	if err := db.Save(&summary).Error; err != nil {
		t.Fatal("cannot create summary", err)
	}
	defer db.Delete(&summary)

	{
		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatal("cannot get summary", w.Code, w.Body)
		} else {
			decoder := json.NewDecoder(w.Body)
			summary := model.Summary{}
			if err := decoder.Decode(&summary); err != nil {
				t.Fatal("cannot get summary", err)
			}
			assert.Equal(t, int64(20), summary.AthleteID, "success get summary")
			assert.Equal(t, int64(3), summary.WeeklyCount, "success get summary")
		}
	}
}

func TestRecalcMySummaryHandler(t *testing.T) {
	user := model.User{ID: 10, AthleteID: 100}
	db := model.DB()
	if err := db.Save(&user).Error; err != nil {
		t.Fatal("cannot create user", err)
	}
	defer db.Delete(&user)

	permission := model.Permission{AthleteID: 100}
	if err := db.Save(&permission).Error; err != nil {
		t.Fatal("cannot create permission", err)
	}
	defer db.Delete(&permission)

	summary := model.Summary{AthleteID: 100}
	if err := db.Save(&summary).Error; err != nil {
		t.Fatal("cannot create summary", err)
	}
	defer db.Delete(&summary)

	{
		api := API{&WebhookClientMock{
			AL: []strava.SummaryActivity{
				strava.SummaryActivity{Distance: 1},
				strava.SummaryActivity{Distance: 2},
			},
		}, &AuthClientMock{}}
		router := gin.New()
		router.GET("/hoge", api.RecalcMySummaryHandler)

		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code, "success recalc summary")
		sum := model.Summary{ID: summary.ID}
		if err := db.First(&sum).Error; err != nil {
			t.Fatal("cannot get summary", err)
		}
		assert.Equal(t, float64(3), sum.WeeklyDistance, "success recalc summary")
		db.Delete(&sum)
	}

	{
		api := API{&WebhookClientMock{
			E: errors.New("error"),
		}, &AuthClientMock{}}
		router := gin.New()
		router.GET("/hoge", api.RecalcMySummaryHandler)

		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, 500, w.Code, "failure recalc summary")
	}
}
