package configs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ServerPort = ":8080"

func ConnectDB() *mongo.Client {
	DBConnectionInfo := getDBConnectionInfo()
	credentials := options.Credential{
		Username:      DBConnectionInfo["username"],
		Password:      DBConnectionInfo["password"],
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "admin",
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(DBConnectionInfo["dburi"]).SetAuth(credentials))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// DB Client instance
var DB *mongo.Client = ConnectDB()

// GetCollection getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("usersDB").Collection(collectionName)

	return collection
}
