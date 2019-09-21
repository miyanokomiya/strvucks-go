package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"strvucks-go/internal/app/handler"
	"strvucks-go/internal/app/model"

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
		exchangeToken(c)
	})

	r.GET("/webhooks", func(c *gin.Context) {
		handler.WebhookVarifyHandler(c)
	})
	r.POST("/webhooks", func(c *gin.Context) {
		handler.WebhookHandler(c)
	})

	r.Run(":" + os.Getenv("PORT"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	config := handler.Config()
	authURL, _ := url.QueryUnescape(config.AuthCodeURL("strvucks", handler.AuthCodeOption()...))

	// you should make this a template in your real application
	fmt.Fprintf(w, `<a href="%s">`, authURL)
	fmt.Fprint(w, `<p>Login by Strava</p>`)
	fmt.Fprint(w, `<img src="/assets/strava.jpg" style="width: 120px; height: auto;" />`)
	fmt.Fprint(w, `</a>`)
}

func exchangeToken(c *gin.Context) {
	code := c.Query("code")

	config := handler.Config()

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Error("Failure exchange token.")
		c.String(400, "Failure exchange token.")
		return
	}

	athlete, ok := token.Extra("athlete").(map[string]interface{})
	if !ok {
		log.Error("Failure get athlete from Strava response.")
		c.String(400, "Failure get athlete from Strava response.")
		return
	}

	idFloat, ok := athlete["id"].(float64)
	if !ok {
		log.Error("Failure get athlete from Strava response.")
		c.String(400, "Failure get athlete from Strava response.")
		return
	}
	id := int64(idFloat)

	username, ok := athlete["username"].(string)
	if !ok {
		log.Error("Failure get athlete from Strava response.")
		c.String(400, "Failure get athlete from Strava response.")
		return
	}

	log.Info("Success get user from Strava response.", id, username)

	permission := model.Permission{
		AthleteID:    id,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.Unix(),
	}

	db := model.DB()

	user := model.User{}
	if orm := db.Where("athlete_id = ?", id).First(&user); orm.Error == nil || orm.RecordNotFound() {
		user.AthleteID = id
		user.Username = username
	} else {
		log.Error("Failure get user.")
		c.String(500, "Failure get user.")
		return
	}

	tx := db.Begin()
	tx = permission.Save(tx)
	tx = user.Save(tx)

	if err := tx.Error; err != nil {
		tx.Rollback()
		log.Error("Failure save token & user.")
		c.String(500, "Failure save token & user.")
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("Failure save token & user.")
		c.String(500, "Failure save token & user.")
		return
	}

	c.Redirect(200, "/?auth=success")
}
