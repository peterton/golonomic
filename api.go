package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type botMode struct {
	Enabled bool   `json:"enabled"`
	IRMode  string `json:"irmode"`
	BotMode string `json:"botmode"`
}

type drivePattern struct {
	Pattern string `json:"pattern"`
	// could add things like
}

func api() {
	quit := make(chan bool)
	currentBotMode := botMode{}

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

	// enable/stop bot mode ("beacon, remote,pattern")
	router.POST("/botmode", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		if r.Body == nil {
			w.WriteHeader(400)
			return
		}

		bm := botMode{}
		err := json.NewDecoder(r.Body).Decode(&bm)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		s := newIRSensor(bm.IRMode)

		if currentBotMode.Enabled != true {
			// no bot currently enabled
			if bm.Enabled {
				// request to enable new botMode
				currentBotMode.Enabled = true
				currentBotMode.IRMode = bm.IRMode
				currentBotMode.BotMode = bm.BotMode
				go remoteControl(s, quit)
			} else {
				// request to disable exiting BotMode
				log.Printf("Attempt to disable, but nothing enabled %s:%s. Ignoring...\n", bm.BotMode, bm.IRMode)
			}
		} else if currentBotMode.BotMode == bm.BotMode {
			// call is for currently enabled botMode
			if bm.Enabled {
				// it's already enabled...
				log.Printf("botMode:IRMode already enabled: %s:%s. Ignoring...\n", bm.BotMode, bm.IRMode)
			} else {
				// dislable it
				currentBotMode.Enabled = false
				currentBotMode.IRMode = ""
				currentBotMode.BotMode = ""
				quit <- true
			}
		} else {
			//there is a bot enabled and request is for different bot
			if bm.Enabled {
				// disable current bot and then enable new one
				currentBotMode.Enabled = false
				currentBotMode.IRMode = ""
				currentBotMode.BotMode = ""
				quit <- true
				currentBotMode.Enabled = true
				currentBotMode.IRMode = bm.IRMode
				currentBotMode.BotMode = bm.BotMode
				go remoteControl(s, quit)
			} else {
				// request to disable for a different BotMode
				log.Printf("Attempt to disable different botMode (%s:%s) from current (%s:%s).Ignoring...\n", bm.BotMode, bm.IRMode, currentBotMode.BotMode, currentBotMode.IRMode)
			}
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
			vectorMove(moveVector{X: 1, Y: 0, W: 0})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: 0, Y: 1, W: 0})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: -1, Y: 0, W: 0})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: 0, Y: -1, W: 0})
			time.Sleep(2 * time.Second)
			stopMotors()
		case "roundedsquare":
			vectorMove(moveVector{X: 1, Y: 0, W: 1})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: 0, Y: 1, W: 1})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: -1, Y: 0, W: 1})
			time.Sleep(2 * time.Second)
			vectorMove(moveVector{X: 0, Y: -1, W: 1})
			time.Sleep(2 * time.Second)
			stopMotors()
		default:
			w.WriteHeader(400)
			return
		}

		w.WriteHeader(204)
	})

	// Handle CORS
	corsRouter := cors.Default().Handler(router)

	// Start server
	log.Println("Listening and serving HTTP on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsRouter))
}
