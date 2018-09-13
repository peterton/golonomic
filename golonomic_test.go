package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCartesianToPolar(t *testing.T) {

	tables := []struct {
		r       float64
		degrees float64
		x       float64
		y       float64
	}{
		{1, 45, 0.7071, 0.7071},
		{1, 90, 0, 1},
		{5, 130, -3.2139, 3.8302},
		{5, 180, -5, 0},
		{5, 220, -3.8302, -3.2139},
		{5, 270, 0, -5},
		{5, 310, 3.2139, -3.8302},
	}

	for _, table := range tables {
		x, y := cartesianToPolar(table.r, table.degrees)
		assert.Equal(t, table.x, x)
		assert.Equal(t, table.y, y)
	}
}
