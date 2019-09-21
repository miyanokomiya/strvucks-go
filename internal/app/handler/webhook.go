package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"strvucks-go/internal/app/model"
	"strvucks-go/pkg/swagger"

	log "github.com/sirupsen/logrus"

	"github.com/antihax/optional"
	"github.com/gin-gonic/gin"
)

// WebhookVarifyHandler varifies webhook from Strava
func WebhookVarifyHandler(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode != "" && token != "" {
		if mode == "subscribe" && token == os.Getenv("STRAVA_VERIFY_TOKEN") {
			log.Info("Webhook verified")
			c.JSON(200, gin.H{
				"hub.challenge": challenge,
			})
		} else {
			c.JSON(403, nil)
		}
	}
}

// WebhookHandler handles webhook from Strava
func WebhookHandler(c *gin.Context) {
	event := model.WebhookEvent{}
	if err := c.BindJSON(&event); err != nil {
		log.Error("Invalid Webhook Body", err)
		c.JSON(400, nil)
		return
	}

	log.Println("activityID:", event.ObjectID)
	log.Println("athleteID:", event.OwnerID)

	if event.ObjectType != "activity" {
		log.Info("Not an activity event and ignore")
		c.JSON(200, nil)
		return
	}

	if event.AspectType != "create" {
		log.Info("Not an create event and ignore")
		c.JSON(200, nil)
		return
	}

	db := model.DB()
	if err := db.Create(&event).Error; err != nil {
		log.Error("Failure:", err)
		c.JSON(500, nil)
		return
	}

	c.JSON(200, nil)

	summary := updateSummary(event.ObjectID, event.OwnerID)
	if summary == nil {
		log.Error("Failure get summary")
		return
	}

	postIfttt(summary, event.ObjectID)
}

func getActivity(activityID int64, athleteID int64) *swagger.DetailedActivity {
	l := log.WithFields(log.Fields{"activityID": activityID, "athleteID": athleteID})
	l.Info("Start get activity from Strava")

	permission := model.Permission{}
	if err := model.DB().Where("athlete_id = ?", athleteID).First(&permission).Error; err != nil {
		l.Error("Failure get permission:", err)
		return nil
	}

	client := Client(&permission)
	sconfig := swagger.NewConfiguration()
	sconfig.HTTPClient = client
	sclient := swagger.NewAPIClient(sconfig)
	activity, _, err := sclient.ActivitiesApi.GetActivityById(context.Background(), activityID, &swagger.GetActivityByIdOpts{IncludeAllEfforts: optional.EmptyBool()})

	if err != nil {
		l.Error("Failure get activity from Strava:", err)
		return nil
	}

	l.Info("Success get activity from Strava")
	return &activity
}

func updateSummary(activityID int64, athleteID int64) *model.Summary {
	activity := getActivity(activityID, athleteID)
	db := model.DB()
	summary := model.Summary{}

	if err := summary.FirstOrInit(db, athleteID).Error; err != nil {
		log.Error("Failure get summary:", err)
		return nil
	}

	summary = summary.Migrate(activity)

	if err := summary.Save(db).Error; err != nil {
		log.Error("Failure save summary:", err)
		return nil
	}

	log.Info("Success save summary")

	return &summary
}

func getIftttURL(athleteID int64) string {
	db := model.DB()
	user := model.User{}
	if err := db.Where("athlete_id = ?", athleteID).First(&user).Error; err != nil {
		log.Error("Failure get user:", err)
		return ""
	}

	return fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", user.IftttMessage, user.IftttKey)
}

func postIfttt(summary *model.Summary, activityID int64) {
	l := log.WithFields(log.Fields{"activityID": activityID, "summaryID": summary.ID})
	l.Info("Start post summary to IFTTT")

	text := summary.GenerateText(activityID)
	body := model.IftttBody{
		Value1: text,
	}

	buff := new(bytes.Buffer)
	json.NewEncoder(buff).Encode(body)

	iftttURL := getIftttURL(summary.AthleteID)
	if iftttURL == "" {
		return
	}

	_, err := http.Post(iftttURL, "application/json; charset=utf-8", buff)
	if err != nil {
		l.Error("Failure post summary to IFTTT:", err)
		return
	}

	l.Info("Success post summary to IFTTT")
}
