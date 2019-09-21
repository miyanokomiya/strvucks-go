package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionSave(t *testing.T) {
	tx := DB().Begin()
	defer tx.Rollback()

	tx = tx.Save(&User{AthleteID: 1})
	tx = tx.Save(&User{AthleteID: 2})

	initial := Permission{AthleteID: 1}
	if err := initial.Save(tx).Error; err != nil {
		t.Fatal("cannot create permission", err)
	}

	assert.Equal(t, int64(1), initial.AthleteID, "create permission")

	exist := Permission{AthleteID: 2, AccessToken: "10"}
	tx = tx.Create(&exist)

	first := Permission{AthleteID: 2, AccessToken: "20"}
	if err := first.Save(tx).Error; err != nil {
		t.Fatal("cannot update permission", err)
	}

	assert.Equal(t, "20", first.AccessToken, "update permission")
}
