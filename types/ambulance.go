package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Ambulance struct {
	ID primitive.ObjectID `bson:"_id"`
	AmbulanceID string `json:"ambulance_id"`
	RegistrationNumber string `json:"registration_number"`
	CreatedAt time.Time `json:"created_at"`
	Location Location `bson:"location" json:"location"`
	Status string `json:"status"` // AVAILABLE | OFFLINE | ONROUTE
}