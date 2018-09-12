package main

import (
	"fmt"
	"log"
	"math"

	"github.com/ev3go/ev3dev"
	"gonum.org/v1/gonum/mat"
)

const (
	a1 = 0
	a2 = 120
	a3 = 240
)

var inverse = mat.NewDense(3, 3, nil)

type Motor struct {
	device   *ev3dev.TachoMotor
	maxSpeed int
}
var motor map[string]*Motor

func setupInverse() {
	data := []float64{
		math.Cos(a1*math.Pi/180), math.Cos(a2*math.Pi/180), math.Cos(a3*math.Pi/180),
		math.Sin(a1*math.Pi/180), math.Sin(a2*math.Pi/180), math.Sin(a3*math.Pi/180),
		1, 1, 1}
	matrix := mat.NewDense(3, 3, data)

	err := inverse.Inverse(matrix)
	if err != nil {
		log.Fatalf("failed to inverse matrix: %v", err)
	}

	fn := mat.Formatted(inverse, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("n = %v", fn)
}

func setupMotors() {
	var err error

	motor["A"].device, err = ev3dev.TachoMotorFor("outA", "lego-ev3-l-motor")
	if err != nil {
		log.Fatalf("failed to find large motor on port A: %v", err)
	}
	motor["B"].device, err = ev3dev.TachoMotorFor("outB", "lego-ev3-l-motor")
	if err != nil {
		log.Fatalf("failed to find large motor on port B: %v", err)
	}
	motor["C"].device, err = ev3dev.TachoMotorFor("outC", "lego-ev3-l-motor")
	if err != nil {
		log.Fatalf("failed to find large motor on port C: %v", err)
	}

	err = motor["A"].device.SetStopAction("brake").Err()
	if err != nil {
		log.Fatalf("failed to set brake stop for large motor on port A: %v", err)
	}
	err = motor["B"].device.SetStopAction("brake").Err()
	if err != nil {
		log.Fatalf("failed to set brake stop for large motor on port B: %v", err)
	}
	err = motor["C"].device.SetStopAction("brake").Err()
	if err != nil {
		log.Fatalf("failed to set brake stop for large motor on port C: %v", err)
	}

	motor["A"].maxSpeed = motor["A"].device.MaxSpeed()
	motor["B"].maxSpeed = motor["B"].device.MaxSpeed()
	motor["C"].maxSpeed = motor["C"].device.MaxSpeed()
}

func move(x, y, z float64) {
	direction := mat.NewDense(1, 3, []float64{x, y, z})
	force := mat.NewDense(1, 3, nil)
	force.Mul(direction, inverse)
	fmt.Print(force, "\n")
}

func main() {
	setupInverse()
	setupMotors()
	
	move(0, 1, 0)	
}