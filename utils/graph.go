package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Point struct {
	X float64
	Y float64
}

func distance(a, b Point) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

func findClosestCentroid(p Point, centroids []Point) int {
	minDistance := math.MaxFloat64
	minIndex := -1

	for i, centroid := range centroids {
		dist := distance(p, centroid)
		if dist < minDistance {
			minDistance = dist
			minIndex = i
		}
	}

	return minIndex
}

func computeNewCentroids(points []Point, assignments []int, K int) []Point {
	newCentroids := make([]Point, K)
	counts := make([]int, K)

	for i := 0; i < len(points); i++ {
		centroidIndex := assignments[i]
		newCentroids[centroidIndex].X += points[i].X
		newCentroids[centroidIndex].Y += points[i].Y
		counts[centroidIndex]++
	}

	for i := 0; i < K; i++ {
		newCentroids[i].X /= float64(counts[i])
		newCentroids[i].Y /= float64(counts[i])
	}

	return newCentroids
}

func kmeans(points []Point, K int, maxIterations int) ([]Point, []int) {
	rand.Seed(time.Now().UnixNano())

	// Initialize centroids randomly
	centroids := make([]Point, K)
	for i := 0; i < K; i++ {
		centroids[i] = points[rand.Intn(len(points))]
	}

	assignments := make([]int, len(points))

	for iter := 0; iter < maxIterations; iter++ {
		// Assign each point to the closest centroid
		for i, p := range points {
			assignments[i] = findClosestCentroid(p, centroids)
		}

		// Compute new centroids based on the assignments
		newCentroids := computeNewCentroids(points, assignments, K)

		// If the centroids didn't change, we're done
		changed := false
		for i := 0; i < K; i++ {
			if distance(centroids[i], newCentroids[i]) > 1e-6 {
				changed = true
				break
			}
		}

		if !changed {
			break
		}

		centroids = newCentroids
	}

	return centroids, assignments
}

func main() {
	points := []Point{
		{1, 1}, {1, 2}, {2, 1}, {2, 2},
		{8, 8}, {8, 9}, {9, 8}, {9, 9},
	}

	K := 2
	maxIterations := 100

	centroids, assignments := kmeans(points, K, maxIterations)

	fmt.Println("Centroids:", centroids)
	fmt.Println("Assignments:", assignments)
}
