package handler

import (
	"errors"
	"testing"
	"time"

	"strvucks-go/internal/app/model"
	"strvucks-go/pkg/swagger"

	"github.com/stretchr/testify/assert"
)

type WebhookClientMock struct {
	A  *swagger.DetailedActivity
	AL []swagger.SummaryActivity
	E  error
}

func (w *WebhookClientMock) getActivity(activityID int64, permission *model.Permission) (*swagger.DetailedActivity, error) {

	return w.A, w.E
}

func (w *WebhookClientMock) getActivitiesInMonth(date time.Time, permission *model.Permission) ([]swagger.SummaryActivity, error) {
	return w.AL, w.E
}

func (w *WebhookClientMock) postTextToIfttt(text string, user *model.User) error {
	return w.E
}

func TestUpdateSummary(t *testing.T) {
	db := model.DB()

	user := model.User{AthleteID: 1}
	permission := model.Permission{AthleteID: user.AthleteID}

	if err := db.Save(&user).Error; err != nil {
		t.Fatal("cannot prepare", err)
	}
	defer db.Delete(&user)
	if err := db.Save(&permission).Error; err != nil {
		t.Fatal("cannot prepare", err)
	}
	defer db.Delete(&permission)

	webhookAct := Webhook{&WebhookClientMock{&swagger.DetailedActivity{}, nil, nil}}
	summaryAct := webhookAct.updateSummary(100, 1)
	assert.Equal(t, int64(1), summaryAct.WeeklyCount, "get updated summary if activity exists")

	webhookNoAct := Webhook{&WebhookClientMock{nil, nil, errors.New("error")}}
	summaryNoAct := webhookNoAct.updateSummary(100, 1)
	assert.Nil(t, summaryNoAct, "get nil if no activity")
}

func TestPostSummaryToIfttt(t *testing.T) {
	db := model.DB()

	user := model.User{AthleteID: 1}
	if err := db.Save(&user).Error; err != nil {
		t.Fatal("cannot prepare", err)
	}
	defer db.Delete(&user)

	validSummary := model.Summary{AthleteID: 1}
	invalidSummary := model.Summary{AthleteID: 2}

	successHook := Webhook{&WebhookClientMock{nil, nil, nil}}
	successRet := successHook.postSummaryToIfttt(&validSummary, 1)
	assert.Nil(t, successRet, "return nil if success post IFTTT")

	invalidAthleteRet := successHook.postSummaryToIfttt(&invalidSummary, 1)
	assert.NotNil(t, invalidAthleteRet, "return error if invalid athlete")

	failureHook := Webhook{&WebhookClientMock{nil, nil, errors.New("error")}}
	failureRet := failureHook.postSummaryToIfttt(&validSummary, 1)
	assert.NotNil(t, failureRet, "return error if success post IFTTT")
}
