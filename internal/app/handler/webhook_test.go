package handler

import (
	"fmt"
	"testing"

	"strvucks-go/internal/app/model"
	"strvucks-go/pkg/swagger"

	"github.com/stretchr/testify/assert"
)

type WebhookClientMock struct {
	A *swagger.DetailedActivity
	E error
}

func (w *WebhookClientMock) GetActivity(activityID int64, permission *model.Permission) (*swagger.DetailedActivity, error) {

	return w.A, w.E
}

func TestUpdateSummary(t *testing.T) {
	db := model.DB()

	user := model.User{AthleteID: 1}
	permission := model.Permission{AthleteID: user.AthleteID}
	if err := db.Save(&user).Save(&permission).Error; err != nil {
		t.Error("cannot prepare", err)
	}
	defer db.Delete(&user).Delete(&permission)

	webhookAct := Webhook{&WebhookClientMock{&swagger.DetailedActivity{}, nil}}
	summaryAct := webhookAct.updateSummary(100, 1)
	assert.Equal(t, int64(1), summaryAct.WeeklyCount, "get updated summary if activity exists")

	webhookNoAct := Webhook{&WebhookClientMock{nil, fmt.Errorf("error")}}
	summaryNoAct := webhookNoAct.updateSummary(100, 1)
	assert.Nil(t, summaryNoAct, "get nil if no activity")
}
