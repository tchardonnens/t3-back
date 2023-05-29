package services

import (
	"t3/m/v2/models"
	"testing"
)

// TestDistance tests the distance function
func TestDistance(t *testing.T) {
	p := models.Point{X: 0, Y: 0}
	q := models.Point{X: 3, Y: 4}
	got := distance(p, q)
	want := 5.0
	if got != want {
		t.Errorf("distance(p, q) = %v; want %v", got, want)
	}
}

// TestGetNearestPoint tests the getNearestPoint function
func TestGetNearestPoint(t *testing.T) {
	points := []models.Point{
		{Id: 1, X: 0, Y: 0},
		{Id: 2, X: 3, Y: 4},
		{Id: 3, X: 1, Y: 1},
	}
	currentPoint := models.Point{X: 0, Y: 0}
	gotPoint, gotIndex := getNearestPoint(currentPoint, points)
	wantPoint := models.Point{Id: 1, X: 0, Y: 0}
	wantIndex := 0
	if gotPoint != wantPoint || gotIndex != wantIndex {
		t.Errorf("getNearestPoint(currentPoint, points) = %v, %v; want %v, %v", gotPoint, gotIndex, wantPoint, wantIndex)
	}
}
