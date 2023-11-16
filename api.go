package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sentinel/geospatial/storage"
	"github.com/sentinel/geospatial/types"

	//	"github.com/lib/pq"
	"time"
)

type APIServer struct {
	ambulanceStorage storage.AmbulanceStorage
	router           *fiber.App
	listenAddr       string
}

func NewAPIServer(listenAddr string, storage storage.AmbulanceStorage) *APIServer {
	return &APIServer{
		listenAddr:       listenAddr,
		ambulanceStorage: storage,
		router:           fiber.New(),
	}
}

// start server and handle different routes and middleware
func (api *APIServer) Run() error {

	// routes
	// create ambulance routes
	api.router.Post("api/v1/create-ambulance", api.createAmbulance)
	// fetch neabry ambulances route
	api.router.Get("api/v1/nearby", api.fetchNearby)

	//	start server and listen incoming requests
	return api.router.Listen(api.listenAddr)
}

// fetch ambulances near the user's location
func (api *APIServer) fetchNearby(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// parse user location
	userLocation := new(types.Location)
	err := c.BodyParser(userLocation)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid user address")
	}

	//	fetch nearby location from location storage
	nearbyAmbulances, err := api.ambulanceStorage.FetchNearby(ctx, userLocation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to fetch nearby")
	}

	//	return a list of ambulances found
	return c.Status(fiber.StatusOK).JSON(nearbyAmbulances)
}

// handles create new ambulance request
func (api *APIServer) createAmbulance(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//parse ambulance json data into ambulance struct
	var ambulanceData = new(types.Ambulance)
	err := c.BodyParser(ambulanceData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid ambulance data")
	}

	//	new ambulance
	ambulance := new(types.Ambulance)
	ambulance.ID, ambulance.AmbulanceID = storage.CreateID()
	ambulance.CreatedAt = time.Now().Local()
	ambulance.RegistrationNumber = ambulanceData.RegistrationNumber
	ambulance.Status = "OFFLINE"
	ambulance.Location = ambulanceData.Location

	//persist ambulance data to data store
	err = api.ambulanceStorage.CreateAmbulance(ctx, ambulance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to create ambulance")
	}
	//return success if ambulance persisted as error
	return c.Status(fiber.StatusOK).JSON("ambulance created successfully")
}
