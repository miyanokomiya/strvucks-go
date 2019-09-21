package handler

import (
	"testing"

	"strvucks-go/internal/app/model"

	"github.com/stretchr/testify/assert"
)

func TestSaveUserAndPermission(t *testing.T) {
	db := model.DB()

	userPre := &model.User{AthleteID: 1}
	permissionPre := &model.Permission{AthleteID: 1}

	if err := db.Save(userPre).Save(permissionPre).Error; err != nil {
		t.Fatal("cannot prepare", err)
	}
	defer db.Delete(&model.User{}, model.User{AthleteID: 1})
	defer db.Delete(&model.Permission{}, model.Permission{AthleteID: 1})

	user := &model.User{AthleteID: 1, Username: "a"}
	permission := &model.Permission{AthleteID: 1, AccessToken: "b"}

	err := saveUserAndPermission(user, permission)
	assert.Nil(t, err, "success save user and permission")
}
