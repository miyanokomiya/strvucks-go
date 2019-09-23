package model

import (
	"os"
	"time"

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

// Permission model
type Permission struct {
	ID           int64
	AthleteID    int64
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

// IftttBody model
type IftttBody struct {
	Value1 string `json:"value1"`
}

// Save Permission by treating AthleteID as primaly
func (p *Permission) Save(db *gorm.DB) *gorm.DB {
	old := Permission{}
	if orm := db.FirstOrInit(&old, Permission{AthleteID: p.AthleteID}); orm.Error != nil {
		return orm
	}

	p.ID = old.ID
	return db.Save(p)
}
