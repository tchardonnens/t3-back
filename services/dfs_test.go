package services

import (
	"t3/m/v2/models"
	"testing"
)

func testAddNeighbors(t *testing.T) {
	cityA := &models.Site{Name: "cityA"}
	cityB := &models.Site{Name: "cityB"}

	AddNeighbour(cityA, cityB)

	if len(cityA.Neighbours) != 1 {
		t.Errorf("expected to have 1 neighbor, got %d", len(cityA.Neighbours))
	}

	if len(cityB.Neighbours) != 1 {
		t.Errorf("expected to have 1 neighbor, got %d", len(cityB.Neighbours))
	}

	if cityA.Neighbours[0] != cityB {
		t.Errorf("Expected site1's neighbour to be site2, got %v", cityA.Neighbours[0])
	}

	if cityB.Neighbours[0] != cityA {
		t.Errorf("Expected site2's neighbour to be site1, got %v", cityB.Neighbours[0])
	}
}
