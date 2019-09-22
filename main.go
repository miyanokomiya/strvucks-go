package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"strvucks-go/internal/app/handler"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Info("Not found .env file")
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		indexHandler(c.Writer, c.Request)
	})

	r.StaticFS("/assets", http.Dir("assets"))

	r.GET("/exchange_token", func(c *gin.Context) {
		handler.ExchangeToken(c)
	})

	webhook := handler.NewWebhook()
	r.GET("/webhooks", func(c *gin.Context) {
		webhook.WebhookVarifyHandler(c)
	})
	r.POST("/webhooks", func(c *gin.Context) {
		webhook.WebhookHandler(c)
	})

	r.Run(":" + os.Getenv("PORT"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	id, err := handler.GetAuthUserID(r)
	if err != nil {
		fmt.Fprint(w, `<p>Not Auth</p>`)
	} else {
		fmt.Fprintf(w, `<p>Your ID: %d</p>`, id)
	}

	config := handler.Config()
	authURL, _ := url.QueryUnescape(config.AuthCodeURL("strvucks", handler.AuthCodeOption()...))

	// you should make this a template in your real application
	fmt.Fprintf(w, `<a href="%s">`, authURL)
	fmt.Fprint(w, `<p>Login by Strava</p>`)
	fmt.Fprint(w, `<img src="/assets/strava.jpg" style="width: 120px; height: auto;" />`)
	fmt.Fprint(w, `</a>`)
}
