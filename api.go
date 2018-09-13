package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
		s := newIRSensor()
		data, _ := json.Marshal(map[string]interface{}{
			"heading":  s.getHeading(),
			"distance": s.getDistance(),
		})
		w.WriteHeader(200)
		w.Write(data)
	})

	router.POST("/vectormove", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var v moveVectors
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		vectorMove(v)
		w.WriteHeader(200)
	})

	// Start server
	log.Println("Listening and serving HTTP on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
