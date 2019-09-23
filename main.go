package main

import (
	"net/http"
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

	apiRoute := r.Group("/api")
	{
		api := handler.NewAPI()
		apiRoute.GET("/strava_auth", api.StravaAuthURL)
		apiRoute.GET("/current_user", api.CurrentUserHandler)
		apiRoute.POST("/current_user", api.UpdateCurrentUserHandler)
		apiRoute.GET("/current_user/summary", api.MySummaryHandler)
	}

	r.StaticFS("/assets", http.Dir("web/assets"))
	r.StaticFS("/web", http.Dir("web/dist"))
	r.StaticFS("/favicon.ico", http.Dir("web/assets/favicon.ico"))

	r.GET("/", indexHandler)

	r.Run(":" + os.Getenv("PORT"))
}

func indexHandler(c *gin.Context) {
	c.Redirect(303, "/web/index.html")
}
