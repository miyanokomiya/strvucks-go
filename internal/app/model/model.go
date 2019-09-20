package model

import (
	"os"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dbInstance *gorm.DB

// DB Get DB
func DB() *gorm.DB {
	if dbInstance != nil {
		return dbInstance
	}

	DBMS := "postgres"
	USER := os.Getenv("POSTGRES_USER")
	PASS := os.Getenv("POSTGRES_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("POSTGRES_DB")

	CONNECT := "user=" + USER + " password=" + PASS + " dbname=" + DBNAME + " host=" + HOST + " sslmode=disable"
	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}

	dbInstance = db
	return dbInstance
}

// Athlete model
type Athlete struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// User model
type User struct {
	ID           int64
	AthleteID    int64
	Username     string
	IftttKey     string
	IftttMessage string
}

// Permission model
type Permission struct {
	ID           int64
	AthleteID    int64
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       int64
}

// WebhookEvent model
type WebhookEvent struct {
	ID             int64
	AspectType     string `json:"aspect_type"`
	EventTime      int64  `json:"event_time"`
	ObjectID       int64  `json:"object_id"`
	ObjectType     string `json:"object_type"`
	OwnerID        int64  `json:"owner_id"`
	SubscriptionID int64  `json:"subscription_id"`
}

type IftttBody struct {
	Value1 string `json:"value1"`
}

func (user *User) Save(db *gorm.DB) *gorm.DB {
	old := User{}
	if orm := db.Where("athlete_id = ?", user.AthleteID).First(&old); orm.RecordNotFound() {
		return db.Create(user)
	} else if orm.Error == nil {
		user.ID = old.ID
		return db.Save(user)
	} else {
		return orm
	}
}

func (permission *Permission) Save(db *gorm.DB) *gorm.DB {
	old := Permission{}
	if orm := db.Where("athlete_id = ?", permission.AthleteID).First(&old); orm.RecordNotFound() {
		return db.Create(permission)
	} else if orm.Error == nil {
		permission.ID = old.ID
		return db.Save(permission)
	} else {
		return orm
	}
}