package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"strvucks-go/internal/app/model"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func StravaAuthURL(t *testing.T) {
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

	api := API{}
	router := gin.New()
	router.GET("/hoge", api.CurrentUserHandler)

	user := model.User{ID: 10}
	db := model.DB()
	if err := db.Save(&user).Error; err != nil {
		t.Fatal("cannot create user", err)
	}
	defer db.Delete(&user)

	{
		token, err := CreateToken(&model.User{ID: 10}, time.Now().Unix()+1000)
		if err != nil {
			t.Fatal("cannot create token")
		}

		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Error("NewRequest URI error")
		}
		req.Header.Set("Authorization", token)

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

	api := API{}
	router := gin.New()
	router.POST("/hoge", api.UpdateCurrentUserHandler)

	user := model.User{ID: 10}
	db := model.DB()
	if err := db.Save(&user).Error; err != nil {
		t.Fatal("cannot create user", err)
	}
	defer db.Delete(&user)

	{
		token, err := CreateToken(&model.User{ID: 10}, time.Now().Unix()+1000)
		if err != nil {
			t.Fatal("cannot create token")
		}

		var jsonStr = []byte(`{"iftttKey":"key", "iftttMessage":"message"}`)
		req, err := http.NewRequest("POST", "/hoge", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Error("NewRequest URI error")
		}
		req.Header.Set("Authorization", token)

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
