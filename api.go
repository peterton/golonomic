package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type controlMode struct {
	Enabled bool `json:"enabled"`
}

func api() {

	// Home
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data, _ := json.Marshal(map[string]string{
			"name": "I'm a 🤖",
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
			fmt.Println(err)
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
			fmt.Println(err)
			return
		}

		s := newIRSensor("IR-REMOTE")
		quit := make(chan bool)
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
		rc := controlMode{}
		err := json.NewDecoder(r.Body).Decode(&rc)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err)
			return
		}

		s := newIRSensor("IR-SEEK")
		quit := make(chan bool)
		if rc.Enabled {
			log.Println("starting beacon tracking mode")
			go remoteControl(s, quit)
		} else {
			log.Println("stopping beacon tracking mode")
			quit <- true
		}

		w.WriteHeader(204)
	})

	// Start server
	log.Println("Listening and serving HTTP on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
