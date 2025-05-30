package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	client,err := mongo.Connect(ctx,options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(),nil)
	if err != nil {
		log.Println("Failed to Connect to mongodb")
		return nil 
	}
	fmt.Println("Successfully Connected to the mongodb")
	return client
	
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, CollectionName string) *mongo.Collection{
	var UserCollection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return UserCollection
}

func ProductData(client *mongo.Client, CollectionName string) *mongo.Collection{
	var ProductCollection = client.Database("Ecommerce").Collection(CollectionName)
	return ProductCollection
}