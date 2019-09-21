package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	ID           int64
	AthleteID    int64
	Username     string
	IftttKey     string
	IftttMessage string
}

// FirstOrInit User by AthleteID
func (u *User) FirstOrInit(db *gorm.DB, athleteID int64) *gorm.DB {
	return db.FirstOrInit(u, User{AthleteID: athleteID})
}

// Save User by treating AthleteID as primaly
func (u *User) Save(db *gorm.DB) *gorm.DB {
	old := User{}
	if orm := db.FirstOrInit(&old, User{AthleteID: u.AthleteID}); orm.Error != nil {
		return orm
	}

	u.ID = old.ID
	return db.Save(u)
}

// IftttURL returns URL of IFTTT webhook
func (u *User) IftttURL() string {
	return fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", u.IftttMessage, u.IftttKey)
}
