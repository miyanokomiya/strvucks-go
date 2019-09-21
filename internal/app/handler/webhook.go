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

// Webhook handles webhook of Strava
type Webhook struct {
	WebhookClient WebhookClient
}

// NewWebhook returns implemented Webhook
func NewWebhook() *Webhook {
	return &Webhook{&WebhookClientImpl{}}
}

// WebhookClient is an external module
type WebhookClient interface {
	getActivity(activityID int64, permission *model.Permission) (*swagger.DetailedActivity, error)
	postTextToIfttt(text string, user *model.User) error
}

// WebhookClientImpl implements WebhookClient
type WebhookClientImpl struct{}

func (w *WebhookClientImpl) getActivity(activityID int64, permission *model.Permission) (*swagger.DetailedActivity, error) {
	config := swagger.NewConfiguration()
	config.HTTPClient = Client(permission)
	client := swagger.NewAPIClient(config)
	activity, _, err := client.ActivitiesApi.GetActivityById(context.Background(), activityID, &swagger.GetActivityByIdOpts{IncludeAllEfforts: optional.EmptyBool()})

	return &activity, err
}

func (w *WebhookClientImpl) postTextToIfttt(text string, user *model.User) error {
	iftttURL := user.IftttURL()
	body := model.IftttBody{
		Value1: text,
	}
	buff := new(bytes.Buffer)
	json.NewEncoder(buff).Encode(body)

	_, err := http.Post(iftttURL, "application/json; charset=utf-8", buff)
	return err
}

// WebhookVarifyHandler varifies webhook from Strava
func (w *Webhook) WebhookVarifyHandler(c *gin.Context) {
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
func (w *Webhook) WebhookHandler(c *gin.Context) {
	event := model.WebhookEvent{}
	if err := c.BindJSON(&event); err != nil {
		log.Error("Invalid Webhook Body:", err)
		c.JSON(400, nil)
		return
	}

	l := log.WithFields(log.Fields{"activityID": event.ObjectID})

	if event.ObjectType != "activity" {
		l.Info("Not activity event and ignore")
		c.JSON(200, nil)
		return
	}

	if event.AspectType != "create" {
		l.Info("Not create event and ignore")
		c.JSON(200, nil)
		return
	}

	db := model.DB()
	if err := db.Create(&event).Error; err != nil {
		l.Error("Failure create event:", err)
		c.JSON(500, nil)
		return
	}
	l.Info("Success create event")

	c.JSON(200, nil)

	summary := w.updateSummary(event.ObjectID, event.OwnerID)
	if summary == nil {
		l.Error("Failure update summary")
		return
	}
	l.Info("Success update summary")

	if err := summary.Save(db).Error; err != nil {
		l.Error("Failure save summary:", err)
		return
	}
	l.Info("Success save summary")

	if err := w.postSummaryToIfttt(summary, event.ObjectID); err != nil {
		l.Error("Failure post IFTTT:", err)
		return
	}
	l.Info("Success post IFTTT")
}

func (w *Webhook) updateSummary(activityID int64, athleteID int64) *model.Summary {
	permission := model.Permission{}
	if err := model.DB().First(&permission, model.Permission{AthleteID: athleteID}).Error; err != nil {
		log.Error("Failure get permission of Strava:", err)
		return nil
	}

	activity, err := w.WebhookClient.getActivity(activityID, &permission)
	if err != nil {
		log.Error("Failure get activity from Strava:", err)
		return nil
	}

	db := model.DB()
	summary := model.Summary{}

	if err := summary.FirstOrInit(db, athleteID).Error; err != nil {
		log.Error("Failure get summary:", err)
		return nil
	}

	summary = summary.Migrate(activity)

	return &summary
}

func (w *Webhook) postSummaryToIfttt(summary *model.Summary, activityID int64) error {
	l := log.WithFields(log.Fields{"activityID": activityID, "summaryID": summary.ID})
	l.Info("Start post summary to IFTTT")

	text := summary.GenerateText(activityID)
	user := model.User{}
	if err := model.DB().First(&user, model.User{AthleteID: summary.AthleteID}).Error; err != nil {
		return fmt.Errorf("Failure get user\n%s", err)
	}

	if err := w.WebhookClient.postTextToIfttt(text, &user); err != nil {
		return fmt.Errorf("Failure get user\n%s", err)
	}

	l.Info("Success post summary to IFTTT")
	return nil
}
