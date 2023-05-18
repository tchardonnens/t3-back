package services

import (
	"math"
)

const (
	OpeningHour int = 9
	ClosingHour int = 17
	TimeToLunch int = 1
)

var ActivityDurations = map[string]int{
	"museum":     2,
	"garden":     1,
	"illustrium": 2,
}

const OneDay = 1

func GetFreeTime(numberOfDays int) int {
	return ((ClosingHour - OpeningHour) - TimeToLunch) * numberOfDays
}

func GetTotalNumberOfActivities(freeTime int) int {
	meanTimePerActivity := (ActivityDurations["garden"] + ActivityDurations["museum"] + ActivityDurations["illustrium"]) / 3
	return freeTime / meanTimePerActivity
}

func GetAmountOfActivitiesPerDay(totalNumberOfActivities, numberOfDays int) int {
	amountOfActivitiesPerDay := totalNumberOfActivities / numberOfDays
	roundedAmountOfActivitiesPerDay := int(math.Floor(float64(amountOfActivitiesPerDay)))
	return roundedAmountOfActivitiesPerDay
}

func IsValidScheduleForTheDay(activitiesType []string) bool {
	totalFreeTime := GetFreeTime(OneDay)
	totalActivityTime := 0
	for _, activity := range activitiesType {
		totalActivityTime += ActivityDurations[activity]
	}
	return totalActivityTime <= totalFreeTime
}
