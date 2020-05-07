package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// const dbUrl = "mongodb://localhost:27017"

const dbUrl = "mongodb://dbUser:dbUserPassword@cluster0-shard-00-00-wt0uq.gcp.mongodb.net:27017,cluster0-shard-00-01-wt0uq.gcp.mongodb.net:27017,cluster0-shard-00-02-wt0uq.gcp.mongodb.net:27017/test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin&retryWrites=true&w=majority"

var DBCon *mongo.Database

//connect to databse
func Connect() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal(err)
	}
	DBCon = client.Database("ecommerce")
	fmt.Println("Connected to MongoDB!")
}

func GetConnection() *mongo.Database {
	return DBCon
}
