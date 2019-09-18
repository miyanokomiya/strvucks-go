package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"strvucks-go/src"
	st "strvucks-go/src/strava"
	"strvucks-go/src/swagger"

	"github.com/antihax/optional"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Not found .env file")
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
		webhookVarifyHandler(c)
	})
	r.POST("/webhooks", func(c *gin.Context) {
		webhookHandler(c)
	})

	r.Run(":" + os.Getenv("PORT"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	config := st.Config()
	authURL, _ := url.QueryUnescape(config.AuthCodeURL("strvucks", st.AuthCodeOption()...))

	// you should make this a template in your real application
	fmt.Fprintf(w, `<a href="%s">`, authURL)
	fmt.Fprint(w, `<p>Login by Strava</p>`)
	fmt.Fprint(w, `<img src="/assets/strava.jpg" style="width: 120px; height: auto;" />`)
	fmt.Fprint(w, `</a>`)
}

func webhookVarifyHandler(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode != "" && token != "" {
		if mode == "subscribe" && token == os.Getenv("STRAVA_VERIFY_TOKEN") {
			log.Println("WEBHOOK_VERIFIED")
			c.JSON(200, gin.H{
				"hub.challenge": challenge,
			})
		} else {
			c.JSON(403, nil)
		}
	}
}

func webhookHandler(c *gin.Context) {
	event := src.WebhookEvent{}
	if err := c.BindJSON(&event); err != nil {
		log.Println("Invalid Webhook Body")
		c.JSON(400, nil)
		return
	}

	log.Println("activityID: ", event.ObjectID)
	log.Println("athleteID: ", event.OwnerID)

	if event.ObjectType != "activity" {
		log.Println("Not an activity event and ignore")
		c.JSON(200, nil)
		return
	}

	if event.AspectType != "create" {
		log.Println("Not an create event and ignore")
		c.JSON(200, nil)
		return
	}

	db := src.DB()
	if err := db.Create(&event).Error; err != nil {
		log.Println("Failure: ", err)
		c.JSON(500, nil)
		return
	}

	summary := updateSummary(event.ObjectID, event.OwnerID)
	if summary == nil {
		log.Println("Failure get summary")
		return
	}

	// postIfttt(summary, event.ObjectID)
}

func postIfttt(summary *src.Summary, activityID int64) {
	db := src.DB()

	user := src.User{}
	if err := db.Where("athlete_id = ?", summary.AthleteID).First(&user).Error; err != nil {
		log.Println("Failure get user: ", err)
		return
	}

	lines := []string{
		"New Act:",
		fmt.Sprintf("%.2fkm", summary.LatestDistance/1000),
		fmt.Sprintf("%dmin", summary.LatestMovingTime/60),
		fmt.Sprintf("%.0fkcal", summary.LatestCalories),
		"\nWeekly:",
		fmt.Sprintf("%.2fkm", summary.WeeklyDistance/1000),
		fmt.Sprintf("%dmin", summary.WeeklyMovingTime/60),
		fmt.Sprintf("%.0fkcal", summary.WeeklyCalories),
		fmt.Sprintf("(%d)", summary.WeeklyCount),
		"\nMonthly:",
		fmt.Sprintf("%.2fkm", summary.MonthlyDistance/1000),
		fmt.Sprintf("%dmin", summary.MonthlyMovingTime/60),
		fmt.Sprintf("%.0fkcal", summary.MonthlyCalories),
		fmt.Sprintf("(%d)", summary.MonthlyCount),
		"\n",
		fmt.Sprintf("https://www.strava.com/activities/%d", activityID),
	}
	text := strings.Join(lines, " ")

	body := src.IftttBody{
		Value1: text,
	}

	buff := new(bytes.Buffer)
	json.NewEncoder(buff).Encode(body)

	iftttURL := fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", user.IftttMessage, user.IftttKey)

	response, err := http.Post(iftttURL, "application/json; charset=utf-8", buff)
	if err != nil {
		log.Println("Failure post ifttt: ", err)
		return
	}
	fmt.Println(response)
	log.Println("Success post ifttt")
}

func updateSummary(activityID int64, athleteID int64) *src.Summary {
	db := src.DB()

	permission := src.Permission{}
	if err := db.Where("athlete_id = ?", athleteID).First(&permission).Error; err != nil {
		log.Println("Failure get permission: ", err)
		return nil
	}

	client := st.Client(&permission)
	sconfig := swagger.NewConfiguration()
	sconfig.HTTPClient = client
	sclient := swagger.NewAPIClient(sconfig)
	activity, _, err := sclient.ActivitiesApi.GetActivityById(context.Background(), activityID, &swagger.GetActivityByIdOpts{IncludeAllEfforts: optional.EmptyBool()})

	if err != nil {
		log.Println("Failure get activity: ", err)
		return nil
	}

	distance := float64(activity.Distance)
	movingTime := int64(activity.MovingTime)
	totalElevationGain := float64(activity.TotalElevationGain)
	calories := float64(activity.Calories)

	monthBaseDate := now.BeginningOfMonth()
	weekBaseDate := now.BeginningOfWeek()

	summary := src.Summary{}
	if orm := db.Where("athlete_id = ?", athleteID).First(&summary); orm.RecordNotFound() {
		summary.AthleteID = athleteID

		summary.MonthlyCount = 1
		summary.MonthlyDistance = distance
		summary.MonthlyMovingTime = movingTime
		summary.MonthlyTotalElevationGain = totalElevationGain
		summary.MonthlyCalories = calories

		summary.WeeklyCount = 1
		summary.WeeklyDistance = distance
		summary.WeeklyMovingTime = movingTime
		summary.WeeklyTotalElevationGain = totalElevationGain
		summary.WeeklyCalories = calories
	} else if orm.Error != nil {
		log.Println("Failure get summary: ", err)
		return nil
	} else {
		if monthBaseDate.Equal(summary.MonthBaseDate) {
			summary.MonthlyCount++
			summary.MonthlyDistance += distance
			summary.MonthlyMovingTime += movingTime
			summary.MonthlyTotalElevationGain += totalElevationGain
			summary.MonthlyCalories += calories
		} else {
			summary.MonthlyCount = 1
			summary.MonthlyDistance = distance
			summary.MonthlyMovingTime = movingTime
			summary.MonthlyTotalElevationGain = totalElevationGain
			summary.MonthlyCalories = calories
		}

		if weekBaseDate.Equal(summary.WeekBaseDate) {
			summary.WeeklyCount++
			summary.WeeklyDistance += distance
			summary.WeeklyMovingTime += movingTime
			summary.WeeklyTotalElevationGain += totalElevationGain
			summary.WeeklyCalories += calories
		} else {
			summary.WeeklyCount = 1
			summary.WeeklyDistance = distance
			summary.WeeklyMovingTime = movingTime
			summary.WeeklyTotalElevationGain = totalElevationGain
			summary.WeeklyCalories = calories
		}
	}

	summary.MonthBaseDate = monthBaseDate
	summary.WeekBaseDate = weekBaseDate
	summary.LatestDistance = distance
	summary.LatestMovingTime = movingTime
	summary.LatestTotalElevationGain = totalElevationGain
	summary.LatestCalories = calories

	if err := summary.Save(db).Error; err != nil {
		log.Println("Failure save summary: ", err)
		return nil
	}

	log.Println("Success save summary: ", summary)

	return &summary
}

func exchangeToken(c *gin.Context) {
	code := c.Query("code")

	config := st.Config()

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Failure exchange token.")
		c.String(400, "Failure exchange token.")
		return
	}

	athlete, ok := token.Extra("athlete").(map[string]interface{})
	if !ok {
		log.Println("Failure get athlete from Strava response.")
		c.String(400, "Failure get athlete from Strava response.")
		return
	}

	idFloat, ok := athlete["id"].(float64)
	if !ok {
		log.Println("Failure get athlete from Strava response.")
		c.String(400, "Failure get athlete from Strava response.")
		return
	}
	id := int64(idFloat)

	username, ok := athlete["username"].(string)
	if !ok {
		log.Println("Failure get athlete from Strava response.")
		c.String(400, "Failure get athlete from Strava response.")
		return
	}
	log.Println("Success get user from Strava response.", id, username)

	permission := src.Permission{
		AthleteID:    id,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.Unix(),
	}

	db := src.DB()

	user := src.User{}
	if orm := db.Where("athlete_id = ?", id).First(&user); orm.Error == nil || orm.RecordNotFound() {
		user.AthleteID = id
		user.Username = username
	} else {
		log.Println("Failure get user.")
		c.String(500, "Failure get user.")
		return
	}

	tx := db.Begin()
	tx = permission.Save(tx)
	tx = user.Save(tx)

	if err := tx.Error; err != nil {
		tx.Rollback()
		log.Println("Failure save token & user.")
		c.String(500, "Failure save token & user.")
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failure save token & user.")
		c.String(500, "Failure save token & user.")
		return
	}

	c.Redirect(200, "/?auth=success")
}
