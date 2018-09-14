package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type controlMode struct {
	Enabled bool `json:"enabled"`
}

type drivePattern struct {
	Pattern string `json:"pattern"`
	// could add things like
}

func api() {
	quit := make(chan bool)

	// Home
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data, _ := json.Marshal(map[string]string{
			"name":    "I'm a 🤖",
			"version": getVersion(),
		})
		w.WriteHeader(200)
		w.Write(data)
	})

	// Sensor Readings
	router.GET("/sensor", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		s := newIRSensor("IR-SEEK")
		data, _ := json.Marshal(map[string]interface{}{
			"heading":  s.getHeading(),
			"distance": s.getDistance(),
		})
		w.WriteHeader(200)
		w.Write(data)
	})

	// Power Readings
	router.GET("/power", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		v, i, vMax, vMin := getPower()
		data, _ := json.Marshal(map[string]interface{}{
			"V":    v,
			"I":    i,
			"P":    i * v / 1000,
			"vMin": vMin / 10,
			"vMax": vMax / 10,
		})
		w.WriteHeader(200)
		w.Write(data)
	})

	// Move based on vector (direction+speed)
	router.POST("/vectormove", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if r.Body == nil {
			w.WriteHeader(400)
			return
		}
		v := moveVector{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}
		vectorMove(v)
		w.WriteHeader(204)
	})

	// Put bot in RC mode
	router.POST("/rc", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if r.Body == nil {
			w.WriteHeader(400)
			return
		}
		rc := controlMode{}
		err := json.NewDecoder(r.Body).Decode(&rc)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		s := newIRSensor("IR-REMOTE")

		if rc.Enabled {
			log.Println("starting remote control mode")
			go remoteControl(s, quit)
		} else {
			log.Println("stopping remote control mode")
			quit <- true
		}

		w.WriteHeader(204)
	})

	// Put bot in beacon tracking mode
	router.POST("/beacon", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if r.Body == nil {
			w.WriteHeader(400)
			return
		}
		beacon := controlMode{}
		err := json.NewDecoder(r.Body).Decode(&beacon)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		s := newIRSensor("IR-SEEK")
		if beacon.Enabled {
			log.Println("starting beacon tracking mode")
			go beaconTracker(s, quit)
		} else {
			log.Println("stopping beacon tracking mode")
			quit <- true
		}

		w.WriteHeader(204)
	})

	// Drive a pattern
	router.POST("/pattern", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if r.Body == nil {
			w.WriteHeader(400)
			return
		}
		dp := drivePattern{}
		err := json.NewDecoder(r.Body).Decode(&dp)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		switch dp.Pattern {
		case "square":
			vectorMove(moveVector{X: 1, Y: 0, S: 0})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: 0, Y: 1, S: 0})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: -1, Y: 0, S: 0})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: 0, Y: -1, S: 0})
			time.Sleep(2 * time.Second)
			stopMotors()
		case "roundedsquare":
			vectorMove(moveVector{X: 1, Y: 0, S: 1})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: 0, Y: 1, S: 1})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: -1, Y: 0, S: 1})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: 0, Y: -1, S: 1})
			time.Sleep(2 * time.Second)
			stopMotors()
		default:
			w.WriteHeader(400)
			return
		}

		w.WriteHeader(204)
	})

	// Start server
	log.Println("Listening and serving HTTP on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
