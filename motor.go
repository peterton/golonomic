package main

import (
	"log"
	"math"

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

type moveVector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
}

func initMotor(m string) *ev3dev.TachoMotor {
	if m == "A" || m == "B" || m == "C" {
		motor, err := ev3dev.TachoMotorFor("ev3-ports:out"+m, "lego-ev3-l-motor")
		if err != nil {
			log.Fatalf("failed to find large motor on out%v: %v", m, err)
		}
		err = motor.SetStopAction("brake").Err()
		if err != nil {
			log.Fatalf("failed to set brake stop for large motor on out%v: %v", m, err)
		}
		return motor
	}
	log.Fatalf("specified unknown motor: %s", m)
	resetMotors()
	return nil
}

func stopMotors() {
	motorA.Command("stop")
	motorB.Command("stop")
	motorC.Command("stop")
}

func resetMotors() {
	motorA.Command("reset")
	motorB.Command("reset")
	motorC.Command("reset")
}

func vectorMove(v moveVector) {
	// if vector is 0,0,0 - do stop
	// if anything else, just run-forever
	if v.X == 0 && v.Y == 0 && v.W == 0 {
		stopMotors()
	} else {
		// relative to the robot, move in direction determined by x,y and angular speed s
		// todo? add abstraction function to provide angle and speed instead of x/y components
		direction := mat.NewDense(3, 1, []float64{v.X, v.Y, v.W})
		force := mat.NewDense(3, 1, nil)
		force.Mul(inverse, direction)

		forceA := int(force.At(0, 0) * float64(motorA.MaxSpeed()))
		forceB := int(force.At(1, 0) * float64(motorB.MaxSpeed()))
		forceC := int(force.At(2, 0) * float64(motorC.MaxSpeed()))
		log.Println("forces", forceA, forceB, forceC)

		// just a test
		motorA.SetSpeedSetpoint(forceA).Command("run-forever")
		motorB.SetSpeedSetpoint(forceB).Command("run-forever")
		motorC.SetSpeedSetpoint(forceC).Command("run-forever")
	}
}

// Converts Polar r distance at degrees angle to x, y Cartesian
// rounded to 4 decimals
func cartesianToPolar(r, degrees float64) (x, y float64) {

	//fmt.Printf("ctp: r:%v, degreeW: %v\n", r, degrees)
	x = math.Cos(degrees*math.Pi/180) * r
	x = math.Round(x*10000) / 10000
	y = r * math.Sin(degrees*math.Pi/180)
	y = math.Round(y*10000) / 10000

	return x, y
}

//moves the robot at degrees angle with angular speed s
func movePolar(degrees, angularSpeed float64) {
	x, y := cartesianToPolar(1.0, degrees)
	mv := moveVector{X: x, Y: y, W: angularSpeed}
	vectorMove(mv)
}

func remoteControl(s *irSensor, quit chan bool) {
	for {
		select {
		case <-quit:
			stopMotors()
			return
		default:
			/*
				hold remote as:
				  [C]
				[A] [D]
				  [B]
			*/
			mv := moveVector{}
			btn := s.getButton()
			log.Printf("RC mode: button %v pressed", btn)
			switch btn {
			case 1: // A
				mv = moveVector{X: 1, Y: 0, W: 0}
			case 2: // B
				mv = moveVector{X: 0, Y: -1, W: 0}
			case 3: // C
				mv = moveVector{X: 0, Y: 1, W: 0}
			case 4: // D
				mv = moveVector{X: -1, Y: 0, W: 0}
			case 5: // A+C
				mv = moveVector{X: 1, Y: 1, W: 0}
			case 8: // B+D
				mv = moveVector{X: -1, Y: -1, W: 0}
			case 9: // BACK, rotate
				mv = moveVector{X: 0, Y: 0, W: 1}
			case 10: // A+B
				mv = moveVector{X: 1, Y: -1, W: 0}
			case 11: // C+D
				mv = moveVector{X: -1, Y: 1, W: 0}
			default: // do nothing
				mv = moveVector{X: 0, Y: 0, W: 0}
			}
			vectorMove(mv)
		}
	}
}

func beaconTracker(s *irSensor, quit chan bool) {
	lastHeading := 0
	for {
		select {
		case <-quit:
			stopMotors()
			return
		default:
			heading := s.getHeading()
			distance := s.getDistance()
			log.Printf("Beacon found in heading %v at distance %v", heading, distance)

			// distance doesn't really matter, we need heading
			// ir sensor is placed at 180 degress (x = 0, y = -1)
			if distance == -128 || distance == 100 {
				// if no beacon found, rotate slowly (x = 0, y = 0, s = 0.5)
				// turn into the direction we last saw the beacon
				if lastHeading < 0 {
					lastHeading = -1
				} else {
					lastHeading = 1
				}
				mv := moveVector{X: 0, Y: 0, W: .5 * float64(lastHeading)}
				vectorMove(mv)
			} else {
				// heading ranges from -25 to 25; what are these values?
				// 0..+25 is 180..>180 degress, 0..-25 is 180..<180 degrees
				//movePolar(180, float64(heading)/25)
				// turn based on heading
				// drive based on distance (speed?) -- go slow!
				mv := moveVector{X: 0, Y: float64(distance) / -100, W: float64(heading) / -25}
				vectorMove(mv)
			}
			lastHeading = heading
		}
	}
}
