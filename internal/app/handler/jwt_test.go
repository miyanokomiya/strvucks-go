package handler

import (
	"net/http"
	"testing"
	"time"

	"strvucks-go/internal/app/model"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	godotenv.Load("../../../.env")

	user := model.User{ID: 10}
	token, err := CreateToken(&user, 100)

	assert.Nil(t, err, "seccess create token")
	assert.NotEqual(t, "", token, "seccess create token")
}

func TestGetAuthUserID(t *testing.T) {
	godotenv.Load("../../../.env")

	getReq := func() *http.Request {
		req, err := http.NewRequest("GET", "/hoge", nil)
		if err != nil {
			t.Fatal("NewRequest URI error")
		}
		return req
	}

	{
		token, err := CreateToken(&model.User{ID: 10}, time.Now().Unix()+1000)
		if err != nil {
			t.Fatal("cannot create token")
		}

		req := getReq()
		req.Header.Set("Authorization", token)
		if id, err := GetAuthUserID(req); err != nil {
			t.Fatal("cannot create token", err)
		} else {
			assert.Equal(t, int64(10), id, "get auth user id from JWT token")
		}
	}

	{
		token, err := CreateToken(&model.User{ID: 10}, time.Now().Unix()-1000)
		if err != nil {
			t.Fatal("cannot create token")
		}

		req := getReq()
		req.Header.Add("Authorization", token)
		_, err = GetAuthUserID(req)
		assert.NotNil(t, err, "cannot get auth user from expired JWT token")
	}
}
