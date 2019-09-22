package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"strvucks-go/internal/app/model"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestBindAuthToken(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Not found .env file", err)
	}

	router := gin.New()
	router.POST("/tokenAuth", func(c *gin.Context) {
		if err := BindAuthToken(c, &model.User{}); err != nil {
			c.String(500, "error")
			return
		}
		c.String(200, "success")
	})

	req, err := http.NewRequest("POST", "/tokenAuth", nil)
	if err != nil {
		t.Fatal("NewRequest URI error")
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code == 200 {
		cookies := w.Result().Cookies()
		assert.Equal(t, 1, len(cookies), "cookies exists")
		assert.Equal(t, JwtName, cookies[0].Name, "cookie name is jwt_token")
	} else {
		t.Fatal("Invalid Status", w.Code)
	}
}

func TestGetAuthUserID(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Not found .env file", err)
	}

	req, err := http.NewRequest("GET", "/hoge", nil)
	if err != nil {
		t.Error("NewRequest URI error")
	}

	token, err := CreateToken(&model.User{ID: 10})
	if err != nil {
		t.Fatal("cannot create token")
	}

	req.AddCookie(&http.Cookie{Name: JwtName, Value: token})
	id, err := GetAuthUserID(req)
	if err != nil {
		t.Fatal("cannot create token", err)
	}

	assert.Equal(t, int64(10), id, "get auth user id from JWT token")
}
