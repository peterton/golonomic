package main

import (
	"log"
	"strconv"

	"github.com/ev3go/ev3dev"
)

var irSensorInstance *irSensor

type irSensor struct {
	raw *ev3dev.Sensor
}

// newIRSensor returns the irSensorInstance if it's already initialized,
// otherwise it performs a lazy initialization of the irSensorInstance.
func newIRSensor(mode string) *irSensor {
	if irSensorInstance != nil {
		return irSensorInstance
	}

	s, err := ev3dev.SensorFor("ev3-ports:in4", "lego-ev3-ir")
	if err != nil {
		log.Fatalf("failed to find large IR sensor on in4: %v", err)
	}
	s.SetMode(mode)
	irSensorInstance = &irSensor{
		raw: s,
	}
	return irSensorInstance
}

// getHeading reads heading value from channel 0 of the IR sensor
func (s *irSensor) getHeading() int {
	v, err := s.raw.Value(0)
	if err != nil {
		log.Printf("failed to read IR data channel 0: %v", err)
	}
	value, _ := strconv.Atoi(v)
	return value
}

// getDistance reads distance value from channel 1 of the IR sensor
func (s *irSensor) getDistance() int {
	v, err := s.raw.Value(1)
	if err != nil {
		log.Printf("failed to read IR data channel 1: %v", err)
	}
	value, _ := strconv.Atoi(v)
	return value
}

// getButton reads which button was pressed from channel 0 of the IR sensor
func (s *irSensor) getButton() int {
	v, err := s.raw.Value(0)
	if err != nil {
		log.Printf("failed to read IR data channel 0: %v", err)
	}
	value, _ := strconv.Atoi(v)
	return value
}

func getPower() (v, i, vMin, vMax float64) {
	var err error
	p := ev3dev.PowerSupply("")
	p = ev3dev.PowerSupply(p.String()) // Cache the driver name if not given.

	v, err = p.Voltage()
	if err != nil {
		log.Fatalf("could not read voltage: %v", err)
	}

	i, err = p.Current()
	if err != nil {
		log.Fatalf("could not read current: %v", err)
	}

	vMax, err = p.VoltageMax()
	if err != nil {
		log.Fatalf("could not read max design voltage: %v", err)
	}

	vMin, err = p.VoltageMin()
	if err != nil {
		log.Fatalf("could not read min design voltage: %v", err)
	}

	log.Printf("current power statW: V=%.2fV I=%.0fmA P=%.3fW (designed voltage range:%.2fV-%.2fV)\n", v, i, i*v/1000, vMin/10, vMax/10)
	return v, i, vMax, vMin
}
