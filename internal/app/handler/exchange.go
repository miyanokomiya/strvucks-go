package handler

import (
	"context"

	"strvucks-go/internal/app/model"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func parseAthlete(token *oauth2.Token) *model.User {
	athlete, ok := token.Extra("athlete").(map[string]interface{})
	if !ok {
		return nil
	}

	idFloat, ok := athlete["id"].(float64)
	if !ok {
		return nil
	}
	athleteID := int64(idFloat)

	username, ok := athlete["username"].(string)
	if !ok {
		return nil
	}

	return &model.User{AthleteID: athleteID, Username: username}
}

// ExchangeToken exchange oauth2 token
func ExchangeToken(c *gin.Context) {
	code := c.Query("code")
	config := Config()
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Error("Failure exchange token.")
		c.String(400, "Failure exchange token.")
		return
	}

	user := parseAthlete(token)
	if user == nil {
		log.Error("Failure get athlete from Strava response.")
		c.String(400, "Failure get athlete from Strava response.")
		return
	}

	db := model.DB()
	if err := db.FirstOrInit(&user).Error; err != nil {
		log.Error("Failure get user.")
		c.String(500, "Failure get user.")
		return
	}

	permission := model.Permission{
		AthleteID:    user.AthleteID,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.Unix(),
	}

	tx := db.Begin().Save(&permission).Save(&user).Commit()
	if err := tx.Error; err != nil {
		tx.Rollback()
		log.Error("Failure save token & user.")
		c.String(500, "Failure save token & user.")
		return
	}

	c.Redirect(200, "/?auth=success")
}
