package services

import (
	"math"
)

const (
	OpeningHour float64 = 10.0
	ClosingHour float64 = 18.0
	TimeToLunch float64 = 1.0
)

var ActivityDurations = map[string]float64{
	"museum":     2.0,
	"garden":     1.0,
	"illustrium": 2.0,
}

const OneDay = 1

func GetFreeTime(numberOfDays int) float64 {
	return ((ClosingHour - OpeningHour) - TimeToLunch) * float64(numberOfDays)
}

func GetTotalNumberOfActivities(freeTimeHours float64) int {
	meanTimePerActivity := math.Ceil((ActivityDurations["garden"] + ActivityDurations["museum"] + ActivityDurations["illustrium"]) / float64(len(ActivityDurations)))
	return int(math.Floor(freeTimeHours / meanTimePerActivity))
}

func GetAmountOfActivitiesPerDay(totalNumberOfActivities float64, numberOfDays int) int {
	amountOfActivitiesPerDay := totalNumberOfActivities / float64(numberOfDays)
	roundedAmountOfActivitiesPerDay := int(math.Ceil(amountOfActivitiesPerDay))
	return roundedAmountOfActivitiesPerDay
}

func IsValidScheduleForTheDay(activitiesType []string) bool {
	totalFreeTime := GetFreeTime(OneDay)
	totalActivityTime := 0.0
	for _, activity := range activitiesType {
		totalActivityTime += ActivityDurations[activity]
	}
	return totalActivityTime <= totalFreeTime
}
