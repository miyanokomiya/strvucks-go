package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/now"
	"github.com/joho/godotenv"
	"github.com/strava/go.strava"
	"strvucks-go/src"
)

var authenticator *strava.OAuthAuthenticator

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Not found .env file")
	}

	strava.ClientId, _ = strconv.Atoi(os.Getenv("STRAVA_CLIENTID"))
	strava.ClientSecret = os.Getenv("STRAVA_CLIENTSECRET")

	authenticator = &strava.OAuthAuthenticator{
		CallbackURL:            "https://" + os.Getenv("CALLBACK_HOST") + "/exchange_token",
		RequestClientGenerator: nil,
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		indexHandler(c.Writer, c.Request)
	})

  r.StaticFS("/assets", http.Dir("assets"))

	r.GET("/exchange_token", func(c *gin.Context) {
		authenticator.HandlerFunc(oAuthSuccess, oAuthFailure)(c.Writer, c.Request)
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
	// you should make this a template in your real application
	fmt.Fprintf(w, `<a href="%s">`, authenticator.AuthorizationURL("state1", strava.Permissions.ViewPrivate, true))
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

	postIfttt(summary, event.ObjectID)
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

	client := strava.NewClient(permission.StravaToken)
	service := strava.NewActivitiesService(client)
	call := service.Get(activityID)

	activity, err := call.Do()
	if err != nil {
		log.Println("Failure get activity: ", err)
		return nil
	}

	distance := activity.Distance
	movingTime := activity.MovingTime
	totalElevationGain := activity.TotalElevationGain
	calories := activity.Calories

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

func oAuthSuccess(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {

	user := src.User{
		AthleteID: auth.Athlete.Id,
		Username:  auth.Athlete.FirstName + auth.Athlete.LastName,
	}
	permission := src.Permission{
		AthleteID:   auth.Athlete.Id,
		StravaToken: auth.AccessToken,
	}

	tx := src.DB().Begin()
	tx = user.Save(tx)
	tx = permission.Save(tx)

	if err := tx.Error; err != nil {
		tx.Rollback()
		fmt.Fprintf(w, "FAILURE: %s", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Fprintf(w, "FAILURE: %s", err)
		return
	}

	fmt.Fprintf(w, "SUCCESS:\nAt this point you can use this information to create a new user or link the account to one of your existing users\n")
	fmt.Fprintf(w, "State: %s\n\n", auth.State)
	fmt.Fprintf(w, "Access Token: %s\n\n", auth.AccessToken)

	fmt.Fprintf(w, "The Authenticated Athlete (you):\n")
	content, _ := json.MarshalIndent(auth.Athlete, "", " ")
	fmt.Fprint(w, string(content))
}

func oAuthFailure(err error, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authorization Failure:\n")

	// some standard error checking
	if err == strava.OAuthAuthorizationDeniedErr {
		fmt.Fprint(w, "The user clicked the 'Do not Authorize' button on the previous page.\n")
		fmt.Fprint(w, "This is the main error your application should handle.")
	} else if err == strava.OAuthInvalidCredentialsErr {
		fmt.Fprint(w, "You provided an incorrect client_id or client_secret.\nDid you remember to set them at the begininng of this file?")
	} else if err == strava.OAuthInvalidCodeErr {
		fmt.Fprint(w, "The temporary token was not recognized, this shouldn't happen normally")
	} else if err == strava.OAuthServerErr {
		fmt.Fprint(w, "There was some sort of server error, try again to see if the problem continues")
	} else {
		fmt.Fprint(w, err)
	}
}
