package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserFirstOrInit(t *testing.T) {
	tx := DB().Begin()
	defer tx.Rollback()

	initial := User{}
	if err := initial.FirstOrInit(tx, 1).Error; err != nil {
		t.Fatal("cannot init user", err)
	}

	assert.Equal(t, int64(1), initial.AthleteID, "init user")

	tx = tx.Create(&User{AthleteID: 2, Username: "10"})

	first := User{}
	if err := first.FirstOrInit(tx, 2).Error; err != nil {
		t.Fatal("cannot find user", err)
	}

	assert.Equal(t, "10", first.Username, "find user")
}

func TestUserSave(t *testing.T) {
	tx := DB().Begin()
	defer tx.Rollback()

	initial := User{AthleteID: 1}
	if err := initial.Save(tx).Error; err != nil {
		t.Fatal("cannot create user", err)
	}

	assert.Equal(t, int64(1), initial.AthleteID, "create user")

	exist := User{AthleteID: 2, Username: "10"}
	tx = tx.Create(&exist)

	first := User{ID: exist.ID, AthleteID: 2, Username: "20"}
	if err := first.Save(tx).Error; err != nil {
		t.Fatal("cannot update user", err)
	}

	assert.Equal(t, "20", first.Username, "update user")
}

func TestUserIftttURL(t *testing.T) {
	user := User{
		IftttMessage: "ifttt_message",
		IftttKey:     "ifttt_key",
	}

	assert.Equal(t, "https://maker.ifttt.com/trigger/ifttt_message/with/key/ifttt_key", user.IftttURL())
}
