package main

import (
	"math"
	"strings"
	"testing"
	"time"
)

const epsilon = 0.0001

const exampleLocation = `
{
    "features": [
        {
            "geometry": {
                "coordinates": [
                    -121.9859701,
                    37.3684139
                ],
                "type": "Point"
            },
            "properties": {
                "accuracyInMeters": 27,
                "id": "-8943391558775469088",
                "photoHeight": 96,
                "photoUrl": "https://latitude.google.com/latitude/apps/badge/api?type=photo&photo=YwbS6zwBAAA.n3jg8QXnZTeXjIeHVjTTig.PjRlBn_guhKWXBqcTtVpXg",
                "photoWidth": 96,
                "placardHeight": 59,
                "placardUrl": "https://latitude.google.com/latitude/apps/badge/api?type=photo_placard&photo=YwbS6zwBAAA.n3jg8QXnZTeXjIeHVjTTig.PjRlBn_guhKWXBqcTtVpXg&moving=true&stale=true&lod=1&format=png",
                "placardWidth": 56,
                "reverseGeocode": "Santa Clara, CA, USA",
                "timeStamp": 1361165342
            },
            "type": "Feature"
        }
    ],
    "type": "FeatureCollection"
}
`

func TestFeatureParsing(t *testing.T) {
	features, err := parseFeatures(strings.NewReader(exampleLocation))
	if err != nil {
		t.Fatalf("Error parsing features: %v", err)
	}
	if len(features) != 1 {
		t.Fatalf("Expected to parse one feature, parsed %v", features)
	}

	f := features[0]

	if f.Geometry.Type != "Point" {
		t.Errorf("Expected a Point, got %v", f.Geometry.Type)
	}

	if f.Properties.ReverseGeocode != "Santa Clara, CA, USA" {
		t.Errorf("Unexpected geocode: %v", f.Properties.ReverseGeocode)
	}
	if f.Properties.AccuracyInMeters != 27 {
		t.Errorf("Incorrect accuracy: %v", f.Properties.AccuracyInMeters)
	}
	ts := f.Timestamp().Format(time.RFC3339Nano)
	if ts != "2013-02-18T05:29:02Z" {
		t.Errorf("Incorrect timestamp: %v", ts)
	}

	if math.Abs(f.Geometry.Latitude()-37.3684139) > epsilon {
		t.Errorf("Invalid latitude: %v", f.Geometry.Latitude())
	}
	if math.Abs(f.Geometry.Longitude() - -121.9859701) > epsilon {
		t.Errorf("Invalid longitude: %v", f.Geometry.Longitude())
	}
}
