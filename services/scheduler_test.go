package services

import (
	"testing"
)

func TestGetFreeTime(t *testing.T) {
	if GetFreeTime(1) != 7.0 {
		t.Errorf("Expected 7.0, got %f", GetFreeTime(1))
	}
}

func TestGetTotalNumberOfActivities(t *testing.T) {
	if GetTotalNumberOfActivities(7.0) != 3 {
		t.Errorf("Expected 3, got %d", GetTotalNumberOfActivities(7.0))
	}
}

func TestGetAmountOfActivitiesPerDay(t *testing.T) {
	if GetAmountOfActivitiesPerDay(6.0, 3) != 2 {
		t.Errorf("Expected 2, got %d", GetAmountOfActivitiesPerDay(6.0, 3))
	}
}

func TestIsValidScheduleForTheDay(t *testing.T) {
	activities := []string{"museum", "garden"}
	if !IsValidScheduleForTheDay(activities) {
		t.Errorf("Expected true, got false")
	}

	activities = []string{"museum", "garden", "illustrium", "museum"}
	if !IsValidScheduleForTheDay(activities) {
		t.Errorf("Expected false, got true")
	}
}
