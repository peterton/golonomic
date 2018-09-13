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
	// angles of the motors reletive to the center of the robot
	a1 = 0
	a2 = 120
	a3 = 240
)

// in order to be able to calculate the forces needed to be applied to each motor
// when a direction and speed is given, the components of the force can be found with:
// (x)   (cos(a1*π/180) cos(a2*π/180) cos(a3*π/180))(f1)
// (y) = (sin(a1*π/180) sin(a2*π/180) sin(a3*π/180))(f2)
// (s)   (1             1             1            )(f3)
// if x,y,s are known, we need the inverse of the big matrix to be multiplied with x,y,s
var inverse = mat.NewDense(3, 3, nil)

var motorA *ev3dev.TachoMotor
var motorB *ev3dev.TachoMotor
var motorC *ev3dev.TachoMotor

type moveVectors struct {
	X float64
	Y float64
	S float64
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
}

func initMotor(m string) *ev3dev.TachoMotor {
	if m == "A" || m == "B" || m == "C" {
		motor, err := ev3dev.TachoMotorFor("ev3-ports:out"+m, "lego-ev3-l-motor")
		if err != nil {
			log.Fatalf("failed to find large motor on out%s: %v", m, err)
		}
		err = motor.SetStopAction("brake").Err()
		if err != nil {
			log.Fatalf("failed to set brake stop for large motor on out%s: %v", m, err)
		}
		return motor
	}
	log.Fatalf("specified unknown motor: %s", m)
	return nil
}

func vectorMove(v moveVectors) {
	// relative to the robot, move in direction determined by x,y and angular speed s
	// todo? add abstraction function to provide angle and speed instead of x/y components
	direction := mat.NewDense(3, 1, []float64{v.X, v.Y, v.S})
	force := mat.NewDense(3, 1, nil)
	force.Mul(inverse, direction)

	forceA := int(force.At(0, 0) * float64(motorA.MaxSpeed()))
	forceB := int(force.At(1, 0) * float64(motorB.MaxSpeed()))
	forceC := int(force.At(2, 0) * float64(motorC.MaxSpeed()))
	fmt.Println("forces", forceA, forceB, forceC)

	// just a test
	motorA.SetSpeedSetpoint(forceA).Command("run-forever")
	motorB.SetSpeedSetpoint(forceB).Command("run-forever")
	motorC.SetSpeedSetpoint(forceC).Command("run-forever")
	time.Sleep(2 * time.Second)
	motorA.Command("stop")
	motorB.Command("stop")
	motorC.Command("stop")
}

// Converts Polar r distance at degrees angle to x, y Cartesian
// rounded to 4 decimals
func cartesianToPolar(r, degrees float64) (x, y float64) {

	//fmt.Printf("ctp: r:%v, degrees: %v\n", r, degrees)
	x = math.Cos(degrees*math.Pi/180) * r
	x = math.Round(x*10000) / 10000
	y = r * math.Sin(degrees*math.Pi/180)
	y = math.Round(y*10000) / 10000

	return x, y
}

//moves the robot at degrees angle for speed s
func movePolar(degrees, speed float64) (x, y float64) {

	x, y = cartesianToPolar(1.0, degrees)
	mv := moveVectors{X: x, Y: y, S: speed}
	vectorMove(mv)
	return x, y
}

func main() {
	setupEV3()
	api()

	for {
		h := irSensorInstance.getHeading()
		d := irSensorInstance.getDistance()

		fmt.Println("Heading:", h, " Distance:", d)
		time.Sleep(1 * time.Second)
	}
}
