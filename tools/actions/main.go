package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.SetFlags(0)

	res, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatalf("Error issuing http request: %v", err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("HTTP error: %v", res.Status)
	}
	defer res.Body.Close()

	fs, err := parseFeatures(res.Body)
	if err != nil {
		log.Fatalf("Error parsing features: %v", err)
	}

	f := fs[0]

	log.Printf("Current location: %v/%v (%v)",
		f.Geometry.Latitude(), f.Geometry.Longitude(),
		f.Properties.ReverseGeocode)
	log.Printf("As of %v (%v ago)", f.Timestamp(),
		time.Now().Sub(f.Timestamp()))
}
