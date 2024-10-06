package db

import (
	"context"
	"fmt"

	"github.com/stasdashkevitch/rest-api/cmd/internal/user"
	"github.com/stasdashkevitch/rest-api/cmd/pkg/logging"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (db *db) Create(ctx context.Context, user user.User) (string, error) {
	db.logger.Debug("--create user")
	result, err := db.collection.InsertOne(ctx, user)

	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	db.logger.Debug("--convert insertedID to ObjectID")
	oid, ok := result.InsertedID.(bson.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	db.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to hex. oid", oid)
}

func (db *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to ObjectID. hex: %s, due to error: %v", id, err)
	}

	filter := bson.M{"id": oid}

	result := db.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		return u, fmt.Errorf("failed to find one user by id: %s", id)
	}

	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user (id: %s) from db, due to error: %v", id, err)
	}

	return u, nil
}

func (db *db) Update(ctx context.Context, user user.User) error {

}
func (db *db) Delete(ctx context.Context, id string) error {

}
