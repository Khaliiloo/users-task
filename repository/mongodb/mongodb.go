package mongodb

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
	"users-task/configs"
)

type MongoDB struct {
	client *mongo.Client
}

func init() {

	db, err := Connect()
	userCollection := db.getCollection("usersDB", "users")

	indexName, err := userCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		configs.Logger.Errorf("couldn't create unique index for `id` field: %q\n", err)
	} else {
		configs.Logger.Printf("unique index %v for `id` field was created\n", indexName)
	}

	indexName, err = userCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		configs.Logger.Errorf("couldn't create unique index for `email` field: %q\n", err)
	} else {
		configs.Logger.Printf("unique index %v for `email` field was created\n", indexName)
	}

}
func Connect() (*MongoDB, error) {
	DBConnectionInfo := getDBConnectionInfo()
	credentials := options.Credential{
		Username:      DBConnectionInfo["username"],
		Password:      DBConnectionInfo["password"],
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "admin",
	}

	var db MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error = nil
	db.client, err = mongo.Connect(ctx, options.Client().ApplyURI(DBConnectionInfo["dburi"]).SetAuth(credentials))

	if configs.Logger == nil {
		configs.Logger = configs.NewLogger("log.txt")
	}
	if err != nil {
		configs.Logger.Fatal(err)
	}

	err = db.client.Ping(ctx, nil)
	if err != nil {
		configs.Logger.Fatal(err)
	}
	configs.Logger.Println("Connected to MongoDB")

	return &db, nil
}

func (db *MongoDB) getCollection(dbName, collName string) *mongo.Collection {
	database := db.client.Database(dbName)
	if database == nil {
		log.Printf("database %v is not existed", dbName)
		return nil
	}
	collection := database.Collection(collName)
	if collection == nil {
		log.Printf("collection %v is not existed", collName)
		return nil
	}
	return collection
}

func getDBConnectionInfo() map[string]string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return map[string]string{"dburi": os.Getenv("DBURI"),
		"username": os.Getenv("USERNAME"),
		"password": os.Getenv("PASSWORD")}
}

/*

	func (db *MongoDB) List(dest interface{}, query string) error {

		return nil
	}


func Get(ctx context.Context, id int) error {
	err := dbGet(ID, dest, query)
	if err != nil {
		return err
	}
	return nil
}
*/
/*
func Create(data interface{}, query string) error {
	err := dbCreate(data, query)
	if err != nil {
		return err
	}
	return nil
}

func Update(ID string, data interface{}, query string) error {
	err := dbUpdate(ID, data, query)
	if err != nil {
		return err
	}
	return nil
}

func Delete(ID string, query string) error {
	err := dbDelete(ID, query)
	if err != nil {
		return err
	}
	return nil
}
*/
