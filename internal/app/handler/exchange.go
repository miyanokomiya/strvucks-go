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
		log.Error("Failure exchange token.", err)
		c.String(400, "Failure exchange token")
		return
	}
	log.Info("Success exchange token")

	user := parseAthlete(token)
	if user == nil {
		log.Error("Failure get athlete from Strava response.")
		c.String(400, "Failure get athlete from Strava response.")
		return
	}
	log.WithFields(log.Fields{
		"AthleteID": user.AthleteID,
		"Username":  user.Username,
	}).Info("Success get athlete")

	permission := &model.Permission{
		AthleteID:    user.AthleteID,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.Unix(),
	}

	if err := saveUserAndPermission(user, permission); err != nil {
		log.Error("Failure save token & user", err)
		c.String(500, "Failure save token & user")
		return
	}
	log.Info("Success save token & user")

	c.Redirect(303, "/?auth=success")
}

func saveUserAndPermission(user *model.User, permission *model.Permission) error {
	tx := model.DB().Begin()
	tx = user.Save(tx)
	tx = permission.Save(tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
