package main

import (
	"fmt"
	"log"
	"math"

	"gonum.org/v1/gonum/mat"
)

const (
	a1 = 0
	a2 = 120
	a3 = 240
)

func main() {
	data := []float64{
		math.Cos(a1*math.Pi/180), math.Cos(a2*math.Pi/180), math.Cos(a3*math.Pi/180),
		math.Sin(a1*math.Pi/180), math.Sin(a2*math.Pi/180), math.Sin(a3*math.Pi/180),
		1, 1, 1}
	matrix := mat.NewDense(3, 3, data)

	inverse := mat.NewDense(3, 3, nil)
	err := inverse.Inverse(matrix)
	if err != nil {
		log.Fatalf("failed to inverse matrix: %v", err)
	}

	fn := mat.Formatted(inverse, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("n = %v", fn)
}