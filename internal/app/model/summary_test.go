package model

import (
	"strings"
	"testing"
	"time"

  "github.com/miyanokomiya/strava-client-go"
	"github.com/jinzhu/now"
	"github.com/stretchr/testify/assert"
)

func TestFirstOrInit(t *testing.T) {
	tx := DB().Begin()
	defer tx.Rollback()

	tx = tx.Save(&User{AthleteID: 1})
	tx = tx.Save(&User{AthleteID: 2})

	initial := Summary{}
	if err := initial.FirstOrInit(tx, 1).Error; err != nil {
		t.Fatal("cannot init summary", err)
	}

	assert.Equal(t, int64(1), initial.AthleteID, "init summary")

	tx = tx.Create(&Summary{AthleteID: 2, LatestDistance: 10})

	first := Summary{}
	if err := first.FirstOrInit(tx, 2).Error; err != nil {
		t.Fatal("cannot find summary", err)
	}

	assert.Equal(t, 10.0, first.LatestDistance, "find summary")
}

func TestSave(t *testing.T) {
	tx := DB().Begin()
	defer tx.Rollback()

	tx = tx.Save(&User{AthleteID: 1})
	tx = tx.Save(&User{AthleteID: 2})

	initial := Summary{AthleteID: 1}
	if err := initial.Save(tx).Error; err != nil {
		t.Fatal("cannot create summary", err)
	}

	assert.Equal(t, int64(1), initial.AthleteID, "create summary")

	exist := Summary{AthleteID: 2, LatestDistance: 10}
	tx = tx.Create(&exist)

	first := Summary{AthleteID: 2, LatestDistance: 20}
	if err := first.Save(tx).Error; err != nil {
		t.Fatal("cannot update summary", err)
	}

	assert.Equal(t, 20.0, first.LatestDistance, "update summary")
}

func TestMigrate(t *testing.T) {
	format := "2006-01-02 15:04:05"
	baseTime, _ := time.Parse(format, "2019-09-10 23:36:00")

	summary := Summary{
		LatestDistance:            1,
		LatestMovingTime:          2,
		LatestTotalElevationGain:  3,
		LatestCalories:            4,
		MonthBaseDate:             now.New(baseTime).BeginningOfMonth(),
		MonthlyCount:              5,
		MonthlyDistance:           6,
		MonthlyMovingTime:         7,
		MonthlyTotalElevationGain: 8,
		MonthlyCalories:           9,
		WeekBaseDate:              now.New(baseTime).BeginningOfWeek(),
		WeeklyCount:               10,
		WeeklyDistance:            11,
		WeeklyMovingTime:          12,
		WeeklyTotalElevationGain:  13,
		WeeklyCalories:            14,
	}

	type Data struct {
		act strava.DetailedActivity
		exp Summary
		mes string
	}

	data := []Data{
		Data{
			act: strava.DetailedActivity{
				StartDate:          baseTime.AddDate(0, 0, 1),
				Distance:           100,
				MovingTime:         200,
				TotalElevationGain: 300,
				Calories:           400,
			},
			exp: Summary{
				LatestDistance:            100,
				LatestMovingTime:          200,
				LatestTotalElevationGain:  300,
				LatestCalories:            400,
				MonthBaseDate:             now.New(baseTime).BeginningOfMonth(),
				MonthlyCount:              6,
				MonthlyDistance:           106,
				MonthlyMovingTime:         207,
				MonthlyTotalElevationGain: 308,
				MonthlyCalories:           409,
				WeekBaseDate:              now.New(baseTime).BeginningOfWeek(),
				WeeklyCount:               11,
				WeeklyDistance:            111,
				WeeklyMovingTime:          212,
				WeeklyTotalElevationGain:  313,
				WeeklyCalories:            414,
			},
			mes: "same month, week => summate month, week",
		},
		Data{
			act: strava.DetailedActivity{
				StartDate:          baseTime.AddDate(0, 0, 7),
				Distance:           100,
				MovingTime:         200,
				TotalElevationGain: 300,
				Calories:           400,
			},
			exp: Summary{
				LatestDistance:            100,
				LatestMovingTime:          200,
				LatestTotalElevationGain:  300,
				LatestCalories:            400,
				MonthBaseDate:             now.New(baseTime).BeginningOfMonth(),
				MonthlyCount:              6,
				MonthlyDistance:           106,
				MonthlyMovingTime:         207,
				MonthlyTotalElevationGain: 308,
				MonthlyCalories:           409,
				WeekBaseDate:              now.New(baseTime.AddDate(0, 0, 7)).BeginningOfWeek(),
				WeeklyCount:               1,
				WeeklyDistance:            100,
				WeeklyMovingTime:          200,
				WeeklyTotalElevationGain:  300,
				WeeklyCalories:            400,
			},
			mes: "same month, different week => summate month, replace week",
		},
		Data{
			act: strava.DetailedActivity{
				StartDate:          baseTime.AddDate(0, 1, 0),
				Distance:           100,
				MovingTime:         200,
				TotalElevationGain: 300,
				Calories:           400,
			},
			exp: Summary{
				LatestDistance:            100,
				LatestMovingTime:          200,
				LatestTotalElevationGain:  300,
				LatestCalories:            400,
				MonthBaseDate:             now.New(baseTime.AddDate(0, 1, 0)).BeginningOfMonth(),
				MonthlyCount:              1,
				MonthlyDistance:           100,
				MonthlyMovingTime:         200,
				MonthlyTotalElevationGain: 300,
				MonthlyCalories:           400,
				WeekBaseDate:              now.New(baseTime.AddDate(0, 1, 0)).BeginningOfWeek(),
				WeeklyCount:               1,
				WeeklyDistance:            100,
				WeeklyMovingTime:          200,
				WeeklyTotalElevationGain:  300,
				WeeklyCalories:            400,
			},
			mes: "different month => replace month, week",
		},
	}

	for _, d := range data {
		assert.Equal(t, d.exp, summary.Migrate(&d.act), d.mes)
	}
}

func TestMigrateBySummary(t *testing.T) {
	format := "2006-01-02 15:04:05"
	baseTime, _ := time.Parse(format, "2019-09-10 23:36:00")

	summary := Summary{
		LatestDistance:            1,
		LatestMovingTime:          2,
		LatestTotalElevationGain:  3,
		LatestCalories:            4,
		MonthBaseDate:             now.New(baseTime).BeginningOfMonth(),
		MonthlyCount:              5,
		MonthlyDistance:           6,
		MonthlyMovingTime:         7,
		MonthlyTotalElevationGain: 8,
		MonthlyCalories:           9,
		WeekBaseDate:              now.New(baseTime).BeginningOfWeek(),
		WeeklyCount:               10,
		WeeklyDistance:            11,
		WeeklyMovingTime:          12,
		WeeklyTotalElevationGain:  13,
		WeeklyCalories:            14,
	}

	type Data struct {
		act strava.SummaryActivity
		exp Summary
		mes string
	}

	data := []Data{
		Data{
			act: strava.SummaryActivity{
				StartDate:          baseTime.AddDate(0, 0, 1),
				Distance:           100,
				MovingTime:         200,
				TotalElevationGain: 300,
			},
			exp: Summary{
				LatestDistance:            100,
				LatestMovingTime:          200,
				LatestTotalElevationGain:  300,
				LatestCalories:            0,
				MonthBaseDate:             now.New(baseTime).BeginningOfMonth(),
				MonthlyCount:              6,
				MonthlyDistance:           106,
				MonthlyMovingTime:         207,
				MonthlyTotalElevationGain: 308,
				MonthlyCalories:           9,
				WeekBaseDate:              now.New(baseTime).BeginningOfWeek(),
				WeeklyCount:               11,
				WeeklyDistance:            111,
				WeeklyMovingTime:          212,
				WeeklyTotalElevationGain:  313,
				WeeklyCalories:            14,
			},
			mes: "same month, week => summate month, week",
		},
	}

	for _, d := range data {
		assert.Equal(t, d.exp, summary.MigrateBySummary(&d.act), d.mes)
	}
}

func TestGenerateText(t *testing.T) {
	format := "2006-01-02 15:04:05"
	baseTime, _ := time.Parse(format, "2019-09-10 23:36:00")

	s := Summary{
		LatestDistance:    1100,
		LatestMovingTime:  121,
		LatestCalories:    1400,
		MonthBaseDate:     now.New(baseTime).BeginningOfMonth(),
		MonthlyCount:      3,
		MonthlyDistance:   2100,
		MonthlyMovingTime: 182,
		MonthlyCalories:   2400,
		WeekBaseDate:      now.New(baseTime).BeginningOfWeek(),
		WeeklyCount:       2,
		WeeklyDistance:    3100,
		WeeklyMovingTime:  243,
		WeeklyCalories:    3400,
	}

	exp := []string{
		"New Act: 1.10km 2min 1400kcal ",
		"Weekly: 3.10km 4min 3400kcal (2) ",
		"Monthly: 2.10km 3min 2400kcal (3) ",
		"https://www.strava.com/activities/999",
	}

	assert.Equal(t, strings.Join(exp, "\n"), s.GenerateText(999), "generate text")
}
