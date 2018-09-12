package main

import (
	"fmt"
	"log"
	"math"
	"time"

    "github.com/ev3go/ev3dev"
	"gonum.org/v1/gonum/mat"
)

const (
	a1 = 0
	a2 = 120
	a3 = 240
)

var inverse = mat.NewDense(3, 3, nil)

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
	fmt.Printf("i = %v", fn)
}

func testMotors() {
	outA, err := ev3dev.TachoMotorFor("ev3-ports:outA", "lego-ev3-l-motor")
	if err != nil {
		log.Fatalf("failed to find medium motor on outA: %v", err)
	}
	err = outA.SetStopAction("brake").Err()
	if err != nil {
		log.Fatalf("failed to set brake stop for medium motor on outA: %v", err)
	}
	maxMedium := outA.MaxSpeed()

	outA.SetSpeedSetpoint(50 * maxMedium / 100).Command("run-forever")
	time.Sleep(time.Second / 2)
	outA.Command("stop")
}

func move(x, y, z float64) {
	direction := mat.NewDense(1, 3, []float64{x, y, z})
	force := mat.NewDense(1, 3, nil)
	force.Mul(direction, inverse)

	fn := mat.Formatted(force, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("force = %v", fn)
}

func main() {
	setupInverse()
	testMotors()

	move(0, 1, 0)	
}