package main

import (
	"encoding/json"
	"io"
	"time"
)

type Geometry struct {
	Coordinates []float64
	Type        string
}

func (g Geometry) Latitude() float64 {
	return g.Coordinates[1]
}

func (g Geometry) Longitude() float64 {
	return g.Coordinates[0]
}

type Properties struct {
	AccuracyInMeters int
	ReverseGeocode   string
	TimeStamp        int
}

type Feature struct {
	Geometry   Geometry
	Properties Properties
}

func (f Feature) Timestamp() time.Time {
	return time.Unix(int64(f.Properties.TimeStamp), 0).UTC()
}

// Parse the features response from google
func parseFeatures(r io.Reader) ([]Feature, error) {
	d := json.NewDecoder(r)
	rv := struct {
		Features []Feature
	}{}
	err := d.Decode(&rv)
	return rv.Features, err
}
