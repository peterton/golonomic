package main

import (
	"fmt"
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

	fmt.Println("Starting test for cartesianToPolar...")
	i := 0
	for _, table := range tables {
		i++
		x, y := cartesianToPolar(table.r, table.degrees)
		assert.Equal(t, table.x, x)
		assert.Equal(t, table.y, y)
		//fmt.Printf("Cartesian of (r:%v,degrees:%v) is (x:%v,y:%v)\n", table.r, table.degrees, x, y)
		// if x != table.x {
		// 	t.Errorf("ERROR: Cartesian of (%v,%v) was incorrect, got: (x:%v,y:%v), want: (x:%v, y:%v).", table.r, table.degrees, x, y, table.x, table.y)
		// } else {
		//    "CORRECT: Cartesian of (%v,%v) was incorrect, got: (x:%v,y:%v), want: (x:%v, y:%v).", table.r, table.degrees, x, y, table.x, table.y)
		//
		// }
	}
	fmt.Printf("\nTests for cartesianToPolar completed; %v tests run\n", i)
}

func TestMovePolar(t *testing.T) {
	tables := []struct {
		degrees float64
		x       float64
		y       float64
	}{
		{45, 0.7071, 0.7071},
		{90, 0, 1},
		{130, -0.6428, 0.766},
		{180, -1, 0},
		{220, -0.766, -0.6428},
		{270, 0, -1},
		{310, 0.6428, -0.766},
	}
	fmt.Println("Starting test for movePolar...")

	i := 0
	for _, table := range tables {
		i++
		x, y := movePolar(table.degrees, 10)
		assert.Equal(t, table.x, x)
		assert.Equal(t, table.y, y)
		//fmt.Printf("Move of degrees:%v returned (x:%v,y:%v)\n", table.degrees, x, y)
	}
	fmt.Printf("\nTests for movePolar completed; %v tests run\n", i)
}
