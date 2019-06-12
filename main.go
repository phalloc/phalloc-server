// Copyright 2019 Tai Groot. All rights reserved.
// Please see the license file at the root of the project.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//	"github.com/bvinc/go-sqlite-lite/sqlite3" // database import

	"github.com/gorilla/mux"
)

// An Event contains all the data correlated with a received packet
type Event struct {
	TS      string   `json:"ts"`     // timestamp
	Client  string   `json:"client"` // who sent this event
	Metrics []Metric `json:"data"`   // array of metrics
	Secret  string   `json:"APIKEY"`
}

type Metric struct {
	Mac       string `json:"mac"`
	TimeStamp string `json:"ts"`
	RSSI      string `json:"rssi"`
}

// used to store macs as dictionary lookups
type MAC struct {
	ID      string    `json:"mac"`
	Payload []LocData `json:"locData"`
}

// used as element of lookup slice
type LocData struct {
	RSSI string `json:"rssi"`
	TS   string `json:"timestamp"`
}

var mod, length int
var macs map[string]MAC
var APIKEY string

// Get all events
func getEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(macs)
}

// Add new event
func createEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	_ = json.NewDecoder(r.Body).Decode(&event)
	if event.Secret != APIKEY {
		http.Error(w, "Not authorized", 401)
		return
	}
	fmt.Println(event)
	for _, ele := range event.Metrics {
		elem, ok := macs[ele.Mac]
		if ok {
			var data LocData
			data.RSSI = ele.RSSI
			data.TS = ele.TimeStamp
			elem.Payload = append(elem.Payload, data)
			macs[ele.Mac] = elem
			fmt.Println("ok")

		} else {
			var data LocData
			var collection MAC
			data.RSSI = ele.RSSI
			data.TS = ele.TimeStamp
			collection.ID = ele.Mac
			collection.Payload = append(collection.Payload, data)
			macs[ele.Mac] = collection
			fmt.Println("added")
		}
	}
	//TODO save new events to database

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Ok")
}

func init() {
	if APIKEY == "" {
		log.Fatalf("Must compile with an API Key!")
	}
	macs = make(map[string]MAC)
}

// Main function
func main() {
	mod = 0
	length = 0

	//TODO database connection & slice init
	// conn, err := sqlite3.Open("mydatabase.db")
	//if err != nil {
	//	os.Exit(1)
	//}
	// defer conn.Close()
	// conn.BusyTimeout(5 * time.Second)

	// Init router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/events", getEvents).Methods("GET")
	r.HandleFunc("/events", createEvent).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}
