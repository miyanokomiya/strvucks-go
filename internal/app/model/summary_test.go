package model

import (
	"testing"
	"time"

	"strvucks-go/pkg/swagger"

	"github.com/jinzhu/now"
	"github.com/stretchr/testify/assert"
)

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
		act swagger.DetailedActivity
		exp Summary
		mes string
	}

	data := []Data{
		Data{
			act: swagger.DetailedActivity{
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
			act: swagger.DetailedActivity{
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
			act: swagger.DetailedActivity{
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
