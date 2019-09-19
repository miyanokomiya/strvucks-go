package handler

import (
	"context"
	"net/http"
	"os"
	"time"

	"strvucks-go/internal/app/model"

	"golang.org/x/oauth2"
)

// Config returns oauth2 config
func Config() oauth2.Config {
	clientID := os.Getenv("STRAVA_CLIENTID")
	clientSecret := os.Getenv("STRAVA_CLIENTSECRET")
	redirectURL := "https://" + os.Getenv("CALLBACK_HOST") + "/exchange_token"
	authRUL := "http://www.strava.com/oauth/authorize"
	tokenURL := "https://www.strava.com/oauth/token"
	scopes := []string{"read,activity:read_all"}

	return oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   authRUL,
			TokenURL:  tokenURL,
			AuthStyle: 1,
		},
		RedirectURL: redirectURL,
		Scopes:      scopes,
	}
}

// AuthCodeOption returns AuthCodeOption
func AuthCodeOption() []oauth2.AuthCodeOption {
	return []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("approval_prompt", "force"),
		oauth2.SetAuthURLParam("response_type", "code"),
	}
}

// Client returns http client with oauth2
func Client(permission *model.Permission) *http.Client {
	token := oauth2.Token{
		AccessToken:  permission.AccessToken,
		RefreshToken: permission.RefreshToken,
		Expiry:       time.Unix(permission.Expiry, 0),
	}

	config := Config()
	client := config.Client(context.Background(), &token)

	return client
}
