package services

import (
	"math"
	"t3/m/v2/models"
)

func distance(p, q models.Point) float64 {
	return math.Hypot(p.X-q.X, p.Y-q.Y)
}

func nearestNeighbourAlgorithm(points []models.Point) []models.Point {
	visited := make([]bool, len(points))
	// Start from the first point
	tour := []models.Point{points[0]}
	visited[0] = true

	for i := 1; i < len(points); i++ {
		last := tour[i-1]
		next := -1
		for j := 0; j < len(points); j++ {
			if !visited[j] && (next == -1 || distance(last, points[j]) < distance(last, points[next])) {
				next = j
			}
		}
		tour = append(tour, points[next])
		visited[next] = true
	}

	return tour
}

func TSP(points []models.Point) []models.Point {
	tour := nearestNeighbourAlgorithm(points)
	return tour
}
