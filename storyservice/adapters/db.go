package adapters

import (
	"context"
	"log"
	"storyservice/helper"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client
var dbConnector sync.Once

func GetDbClient() *mongo.Client {
	dbConnector.Do(func() {
		initDb()
	})

	return dbClient
}

func initDb() {
	var err error
	uri := helper.GetDsn()
	dbClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Println("Error on db connection")
		log.Fatal(err)
	} else {
		log.Println("Connected to database")
	}
}
