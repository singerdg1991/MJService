package database

import (
	"context"
	"fmt"
	"github.com/hoitek/Kit/retry"
	logger "github.com/hoitek/Logger"
	"github.com/hoitek/Maja-Service/config"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// MongoDB is a global variable for the mongoDB connection
var MongoDB *mongo.Client

// ConnectMongoDB connects to mongoDB
func ConnectMongoDB() *mongo.Client {
	// Get config
	var (
		HOST     = config.AppConfig.DatabaseMongoDBHost
		USER     = config.AppConfig.DatabaseMongoDBUser
		PASSWORD = config.AppConfig.DatabaseMongoDBPass
		DB_NAME  = config.AppConfig.DatabaseMongoDBName
		PORT     = config.AppConfig.DatabaseMongoDBPort
	)

	// Try to get mongoDB connection for N amount of time
	client, err := retry.Get(func() (*mongo.Client, error) {
		ctx := context.Background()
		uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", USER, PASSWORD, HOST, PORT, DB_NAME)
		log.Printf("Connecting to MongoDB: %s", uri)
		clientOptions := options.Client().ApplyURI(uri)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			return nil, err
		}
		return client, nil
	}, 2, 3)

	// Handle error
	if err != nil {
		logger.Error(err)
		panic(err)
	} else {
		// Log when connection succeed
		logger.Info("Connected to MongoDB!")
	}

	// Set global connection
	MongoDB = client

	return client
}
