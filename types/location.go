package types

import "time" 

type Location struct {
	Type string `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}