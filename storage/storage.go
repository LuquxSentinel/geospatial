package storage

import (
	"context"

	"github.com/sentinel/geospatial/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AmbulanceStorage interface {
	FetchNearby(ctx context.Context, userLocation *types.Location) ([]*types.Ambulance, error)
	CreateAmbulance(ctx context.Context, ambulance *types.Ambulance) error
	UpdateLocation(ctx context.Context, ambulanceID string, location *types.Location) error
	UpdateStatus(ctx context.Context, ambulanceID string, status string) error
}

type MongoDBAmbulanceStorage struct {
	collection *mongo.Collection
}

func NewMongoDBAmbulanceStorage(collection *mongo.Collection) (*MongoDBAmbulanceStorage, error) {
	indexModel := mongo.IndexModel{
		Keys: primitive.D{primitive.E{Key: "location.coordinates", Value: "2dsphere"}},
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return nil, err
	}
	return &MongoDBAmbulanceStorage{
		collection: collection,
	}, nil
}

func (store *MongoDBAmbulanceStorage) FetchNearby(ctx context.Context, userLocation *types.Location) ([]*types.Ambulance, error) {

	filter := primitive.D{
		primitive.E{Key: "location.coordinates", Value: primitive.D{
			primitive.E{Key: "$near", Value: primitive.D{
				primitive.E{Key: "$geometry", Value: userLocation},
				primitive.E{Key: "$maxDistance", Value: 500},
			},
			},
		}}}

	var ambulances []*types.Ambulance
	result, err := store.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = result.All(ctx, &ambulances)
	if err != nil {
		return nil, err
	}

	return ambulances, nil
}

func (store *MongoDBAmbulanceStorage) CreateAmbulance(ctx context.Context, ambulance *types.Ambulance) error {
	_, err := store.collection.InsertOne(ctx, ambulance)
	if err != nil {
		return err
	}

	return nil
}

func (store *MongoDBAmbulanceStorage) UpdateLocation(ctx context.Context, ambulanceID string, location *types.Location) error {
	return nil
}

func (store *MongoDBAmbulanceStorage) UpdateStatus(ctx context.Context, ambulanceID string, status string) error {
	return nil
}

func CreateID() (primitive.ObjectID, string) {
	objectID := primitive.NewObjectID()
	stringID := objectID.Hex()

	return objectID, stringID
}
