package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"strvucks-go/internal/app/model"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

// JwtName is cookie name of JWT token
var JwtName = "jwt_token"

func getRsaPrivate() ([]byte, error) {
	pri64 := os.Getenv("JWT_RSA_PRI")
	return base64.StdEncoding.DecodeString(pri64)
}

func getRsaPublic() ([]byte, error) {
	pri64 := os.Getenv("JWT_RSA_PUB")
	return base64.StdEncoding.DecodeString(pri64)
}

// BindAuthToken initializes JWT token
func BindAuthToken(c *gin.Context, user *model.User) error {
	tokenString, err := CreateToken(user)
	if err != nil {
		return fmt.Errorf("Failure create jwt token\n%s", err)
	}

	c.SetCookie(JwtName, tokenString, 3600, "", "", false, false)
	return nil
}

// CreateToken returns JWT token
func CreateToken(user *model.User) (string, error) {
	signBytes, err := getRsaPrivate()
	if err != nil {
		return "", fmt.Errorf("Failure get rsa\n%s", err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", fmt.Errorf("Failure parse rsa\n%s", err)
	}

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", fmt.Errorf("Failure sign jwt\n%s", err)
	}

	return tokenString, nil
}

// GetAuthUserID returns user id from JWT token
func GetAuthUserID(r *http.Request) (int64, error) {
	verifyBytes, err := getRsaPublic()
	if err != nil {
		return 0, fmt.Errorf("Failure get rsa pub\n%s", err)
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return 0, fmt.Errorf("Failure parse rsa pub\n%s", err)
	}

	token, err := request.ParseFromRequest(r, &CookieExtractor{}, func(token *jwt.Token) (interface{}, error) {
		_, err := token.Method.(*jwt.SigningMethodRSA)
		if !err {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("Unauthorized")
	}

	claims := token.Claims.(jwt.MapClaims)
	idFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("Not found ID from token")
	}

	return int64(idFloat), nil
}

// CookieExtractor extracts token from cookie
type CookieExtractor []string

// ExtractToken extracts token from cookie
func (e CookieExtractor) ExtractToken(req *http.Request) (string, error) {
	token, err := req.Cookie(JwtName)
	if err != nil {
		return "", request.ErrNoTokenInRequest
	}
	return token.Value, nil
}
