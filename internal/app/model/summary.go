package model

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/now"
	"github.com/miyanokomiya/strava-client-go"
)

// Summary model
type Summary struct {
	ID                        int64     `json:"id"`
	AthleteID                 int64     `json:"athleteId"`
	LatestDistance            float64   `json:"latestDistance"`
	LatestMovingTime          int64     `json:"latestMovingTime"`
	LatestTotalElevationGain  float64   `json:"latestTotalElevationGain"`
	LatestCalories            float64   `json:"latestCalories"`
	MonthBaseDate             time.Time `json:"monthBaseDate"`
	MonthlyCount              int64     `json:"monthlyCount"`
	MonthlyDistance           float64   `json:"monthlyDistance"`
	MonthlyMovingTime         int64     `json:"monthlyMovingTime"`
	MonthlyTotalElevationGain float64   `json:"monthlyTotalElevationGain"`
	MonthlyCalories           float64   `json:"monthlyCalories"`
	WeekBaseDate              time.Time `json:"weekBaseDate"`
	WeeklyCount               int64     `json:"weeklyCount"`
	WeeklyDistance            float64   `json:"weeklyDistance"`
	WeeklyMovingTime          int64     `json:"weeklyMovingTime"`
	WeeklyTotalElevationGain  float64   `json:"weeklyTotalElevationGain"`
	WeeklyCalories            float64   `json:"weeklyCalories"`
}

// FirstOrInit Summary by AthleteID
func (s *Summary) FirstOrInit(db *gorm.DB, athleteID int64) *gorm.DB {
	return db.FirstOrInit(s, Summary{AthleteID: athleteID})
}

// Save Summary by treating AthleteID as primaly
func (s *Summary) Save(db *gorm.DB) *gorm.DB {
	old := Summary{}
	if orm := db.FirstOrInit(&old, Summary{AthleteID: s.AthleteID}); orm.Error != nil {
		return orm
	}

	s.ID = old.ID
	return db.Save(s)
}

// Migrate Summary
func (s Summary) Migrate(activity *strava.DetailedActivity) Summary {
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

// MigrateBySummary migrate Summary from SummaryActivity
func (s Summary) MigrateBySummary(activity *strava.SummaryActivity) Summary {
	return s.Migrate(&strava.DetailedActivity{
		StartDate:          activity.StartDate,
		Distance:           activity.Distance,
		MovingTime:         activity.MovingTime,
		TotalElevationGain: activity.TotalElevationGain,
		Calories:           0,
	})
}

// GenerateText generates text from Summary
func (s *Summary) GenerateText(activityID int64) string {
	lines := []string{
		"New Act:",
		fmt.Sprintf("%.2fkm", s.LatestDistance/1000),
		formatTime(s.LatestMovingTime),
		formatLap(s.LatestDistance, s.LatestMovingTime),
		fmt.Sprintf("%.0fkcal", s.LatestCalories),
		"\nWeekly:",
		fmt.Sprintf("%.2fkm", s.WeeklyDistance/1000),
		formatTime(s.WeeklyMovingTime),
		formatLap(s.WeeklyDistance, s.WeeklyMovingTime),
		fmt.Sprintf("(%d)", s.WeeklyCount),
		"\nMonthly:",
		fmt.Sprintf("%.2fkm", s.MonthlyDistance/1000),
		formatTime(s.MonthlyMovingTime),
		formatLap(s.MonthlyDistance, s.MonthlyMovingTime),
		fmt.Sprintf("(%d)", s.MonthlyCount),
		fmt.Sprintf("\nhttps://www.strava.com/activities/%d", activityID),
	}
	return strings.Join(lines, " ")
}

func formatTime(sec int64) string {
	h := math.Floor(float64(sec) / 60 / 60)
	m := math.Floor(float64(sec)/60 - (60 * h))
	s := sec % 60
	if h == 0 {
		return fmt.Sprintf("%02.0f:%02d", m, s)
	}
	return fmt.Sprintf("%1.0f:%02.0f:%02d", h, m, s)
}

func formatLap(meter float64, sec int64) string {
	lap := float64(sec) / (meter / 1000)
	m := math.Floor(lap / 60)
	s := lap - (60 * m)
	return fmt.Sprintf("%.0f:%02.0f/km", m, s)
}
