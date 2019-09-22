package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"strvucks-go/internal/app/model"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
)

// JwtName is cookie name of JWT token
var JwtName = "jwt_token"

var rsaPri []byte
var rsaPub []byte

func getRsaPrivate() ([]byte, error) {
	if len(rsaPri) != 0 {
		return rsaPri, nil
	}

	pri64 := os.Getenv("JWT_RSA_PRI")
	buff, err := base64.StdEncoding.DecodeString(pri64)
	rsaPri = buff
	return buff, err
}

func getRsaPublic() ([]byte, error) {
	if len(rsaPub) != 0 {
		return rsaPub, nil
	}

	pub64 := os.Getenv("JWT_RSA_PUB")
	buff, err := base64.StdEncoding.DecodeString(pub64)
	rsaPub = buff
	return buff, err
}

// CreateToken returns JWT token
func CreateToken(user *model.User, expiry int64) (string, error) {
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
	claims["expiry"] = expiry

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

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		_, err := token.Method.(*jwt.SigningMethodRSA)
		if !err {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("Unauthorized\n%s", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	expiry, ok := claims["expiry"].(float64)
	if !ok {
		return 0, fmt.Errorf("Not found expiry in token")
	}
	if int64(expiry) <= time.Now().Unix() {
		return 0, fmt.Errorf("Token is expired at %d", int64(expiry))
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("Not found ID in token")
	}

	return int64(idFloat), nil
}
