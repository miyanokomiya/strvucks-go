package model

import (
	"time"

	"strvucks-go/pkg/swagger"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Summary model
type Summary struct {
	ID                        int64
	AthleteID                 int64
	LatestDistance            float64
	LatestMovingTime          int64
	LatestTotalElevationGain  float64
	LatestCalories            float64
	MonthBaseDate             time.Time
	MonthlyCount              int64
	MonthlyDistance           float64
	MonthlyMovingTime         int64
	MonthlyTotalElevationGain float64
	MonthlyCalories           float64
	WeekBaseDate              time.Time
	WeeklyCount               int64
	WeeklyDistance            float64
	WeeklyMovingTime          int64
	WeeklyTotalElevationGain  float64
	WeeklyCalories            float64
}

// FindOrNew Summary by AthleteID
func (s *Summary) FindOrNew(db *gorm.DB, athleteID int64) *gorm.DB {
	orm := db.Where("athlete_id = ?", athleteID).First(s)
	if orm.RecordNotFound() {
		s.AthleteID = athleteID
		return db
	} else if orm.Error != nil {
		return orm
	}
	return db
}

// Migrate Summary
func (s Summary) Migrate(activity *swagger.DetailedActivity) Summary {
	n := now.New(activity.StartDate)
	monthBaseDate := n.BeginningOfMonth()
	weekBaseDate := n.BeginningOfWeek()
	distance := float64(activity.Distance)
	movingTime := int64(activity.MovingTime)
	totalElevationGain := float64(activity.TotalElevationGain)
	calories := float64(activity.Calories)

	if monthBaseDate.Equal(s.MonthBaseDate) {
		s.MonthlyCount++
		s.MonthlyDistance += distance
		s.MonthlyMovingTime += movingTime
		s.MonthlyTotalElevationGain += totalElevationGain
		s.MonthlyCalories += calories
	} else {
		s.MonthlyCount = 1
		s.MonthlyDistance = distance
		s.MonthlyMovingTime = movingTime
		s.MonthlyTotalElevationGain = totalElevationGain
		s.MonthlyCalories = calories
	}

	if weekBaseDate.Equal(s.WeekBaseDate) {
		s.WeeklyCount++
		s.WeeklyDistance += distance
		s.WeeklyMovingTime += movingTime
		s.WeeklyTotalElevationGain += totalElevationGain
		s.WeeklyCalories += calories
	} else {
		s.WeeklyCount = 1
		s.WeeklyDistance = distance
		s.WeeklyMovingTime = movingTime
		s.WeeklyTotalElevationGain = totalElevationGain
		s.WeeklyCalories = calories
	}

	s.MonthBaseDate = monthBaseDate
	s.WeekBaseDate = weekBaseDate
	s.LatestDistance = distance
	s.LatestMovingTime = movingTime
	s.LatestTotalElevationGain = totalElevationGain
	s.LatestCalories = calories

	return s
}
