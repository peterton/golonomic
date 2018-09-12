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

type tachoMotor struct {
	*ev3dev.TachoMotor
	maxSpeed int
}

var motorA tachoMotor
var motorB tachoMotor
var motorC tachoMotor

func setupInverse() {
	data := []float64{
		math.Cos(a1 * math.Pi / 180), math.Cos(a2 * math.Pi / 180), math.Cos(a3 * math.Pi / 180),
		math.Sin(a1 * math.Pi / 180), math.Sin(a2 * math.Pi / 180), math.Sin(a3 * math.Pi / 180),
		1, 1, 1}
	matrix := mat.NewDense(3, 3, data)

	err := inverse.Inverse(matrix)
	if err != nil {
		log.Fatalf("failed to inverse matrix: %v", err)
	}

	fn := mat.Formatted(inverse, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("i = %v", fn)
}

func setupMotors() {
	var err error

	motorA.TachoMotor, err = ev3dev.TachoMotorFor("ev3-ports:outA", "lego-ev3-l-motor")
	if err != nil {
		log.Fatalf("failed to find large motor on outA: %v", err)
	}
	err = motorA.SetStopAction("brake").Err()
	if err != nil {
		log.Fatalf("failed to set brake stop for large motor on outA: %v", err)
	}
	motorA.maxSpeed, err = motorA.MaxSpeed()

	motorB.TachoMotor, err = ev3dev.TachoMotorFor("ev3-ports:outB", "lego-ev3-l-motor")
	if err != nil {
		log.Fatalf("failed to find large motor on outB: %v", err)
	}
	err = motorB.SetStopAction("brake").Err()
	if err != nil {
		log.Fatalf("failed to set brake stop for large motor on outB: %v", err)
	}
	motorB.maxSpeed, err = motorB.MaxSpeed()

	motorC.TachoMotor, err = ev3dev.TachoMotorFor("ev3-ports:outC", "lego-ev3-l-motor")
	if err != nil {
		log.Fatalf("failed to find large motor on outC: %v", err)
	}
	err = motorC.SetStopAction("brake").Err()
	if err != nil {
		log.Fatalf("failed to set brake stop for large motor on outC: %v", err)
	}
	motorC.maxSpeed, err = motorC.MaxSpeed()
}

func move(x, y, z float64) {
	direction := mat.NewDense(1, 3, []float64{x, y, z})
	force := mat.NewDense(1, 3, nil)
	force.Mul(direction, inverse)

	fn := mat.Formatted(force, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("force = %v", fn)

	// just a test
	motorA.SetSpeedSetpoint(50 * motorA.maxSpeed / 100).Command("run-forever")
	time.Sleep(time.Second / 2)
	motorA.Command("stop")
}

func main() {
	setupInverse()
	setupMotors()

	move(0, 1, 0)
}
