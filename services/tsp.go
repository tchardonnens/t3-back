package services

import (
	"fmt"
	"math"
	"t3/m/v2/models"
)

func distance(p, q models.Point) float64 {
	return math.Sqrt(math.Pow(q.X-p.X, 2) + math.Pow(q.Y-p.Y, 2))
}

func getNearestPoint(currentpoint models.Point, points []models.Point) (models.Point, int) {
	minIndex := -1
	minDist := math.MaxFloat64

	for i, point := range points {
		if !point.Visited {
			dist := distance(currentpoint, point)
			if dist < minDist {
				minIndex = i
				minDist = dist
			}
		}
	}

	if minIndex == -1 {
		return models.Point{}, -1
	}

	return points[minIndex], minIndex
}

func planTour(points []models.Point, maxPerDay int, maxDist float64) [][]models.Point {
	days := [][]models.Point{}
	visitedPoints := 0

	for visitedPoints < len(points) {
		day := []models.Point{}
		currentPoint := points[0] // or wherever you want to start
		for len(day) < maxPerDay {
			nearestPoint, nearestIndex := getNearestPoint(currentPoint, points)
			if nearestIndex == -1 || distance(currentPoint, nearestPoint) > maxDist {
				break
			}
			day = append(day, nearestPoint)
			points[nearestIndex].Visited = true
			currentPoint = nearestPoint
			visitedPoints++
		}

		days = append(days, day)
	}

	return days
}

func TSP(points []models.Point, maxPerDay int) [][]models.Point {
	maxDist := 10.0
	tour := planTour(points, maxPerDay, maxDist)

	for i, day := range tour {
		fmt.Println("Day", i+1)
		for _, point := range day {
			fmt.Println("  Visit:", point.Name)
		}
	}

	return tour
}
