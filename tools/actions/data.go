package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Geometry struct {
	Coordinates []float64
	Type        string
}

type latitude float64
type longitude float64

func (l latitude) String() string {
	if l > 0 {
		return fmt.Sprintf("%vN", float64(l))
	} else if l < 0 {
		return fmt.Sprintf("%vS", float64(0.0-l))
	}
	return "0"
}

func (l longitude) String() string {
	if l > 0 {
		return fmt.Sprintf("%vE", float64(l))
	} else if l < 0 {
		return fmt.Sprintf("%vW", float64(0.0-l))
	}
	return "0"
}

func (g Geometry) Latitude() latitude {
	return latitude(g.Coordinates[1])
}

func (g Geometry) Longitude() longitude {
	return longitude(g.Coordinates[0])
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
