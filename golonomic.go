package main

import (
	"log"
	"math"
	"strings"

	"gonum.org/v1/gonum/mat"
)

var (
	commit  string
	builtAt string
	builtBy string
	builtOn string
)

func getVersion() string {
	s := []string{
		"commit:", commit,
		"built @", builtAt,
		"by", builtBy,
		"on", builtOn}
	return strings.Join(s, " ")
}

func setupEV3() {
	data := []float64{
		math.Cos(a1 * math.Pi / 180), math.Cos(a2 * math.Pi / 180), math.Cos(a3 * math.Pi / 180),
		math.Sin(a1 * math.Pi / 180), math.Sin(a2 * math.Pi / 180), math.Sin(a3 * math.Pi / 180),
		1, 1, 1}
	matrix := mat.NewDense(3, 3, data)

	err := inverse.Inverse(matrix)
	if err != nil {
		log.Fatalf("failed to inverse matrix: %v", err)
	}

	motorA = initMotor("A")
	motorB = initMotor("B")
	motorC = initMotor("C")

	getPower()
}

func main() {
	log.Print("version: ", getVersion())

	setupEV3()
	api()
}
