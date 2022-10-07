package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseService struct {
	ctx      context.Context
	client   *mongo.Client
	database string
}

func NewService() *DatabaseService {
	return &DatabaseService{}
}

// Connects the service instance to the MongoDB.
func (s *DatabaseService) Connect(host, username, password, database string, port int) {
	s.database = database
	const connString = "mongodb://%v:%v@%v:%v/%v"
	dsn := fmt.Sprintf(connString,
		username, password, host, port, database)
	clientOptions := options.Client().ApplyURI(dsn)
	mongoClient, err := mongo.Connect(s.ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	s.client = mongoClient

	if err = s.client.Ping(s.ctx, nil); err != nil {
		log.Fatal(err)
	}
}

// Returns the collection for the connected database.
func (s *DatabaseService) GetCollection(collection string) *mongo.Collection {
	return s.client.Database(s.database).Collection(collection)
}
