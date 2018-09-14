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
