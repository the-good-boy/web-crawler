package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseConnection struct {
	access     bool
	uri        string
	client     *mongo.Client
	collection *mongo.Collection
}

func (d *DatabaseConnection) connect() {
	fmt.Print("Starting connection...", d.access)
	if d.access {
		// Connect to database
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		d.uri = os.Getenv("MONGODB_URI")
		fmt.Println(d.uri)
		ctx := context.TODO()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri).SetServerAPIOptions(serverAPI))
		if err != nil {
			panic(err)
		}
		if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
			panic(err)
		}
		fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
		d.client = client
		d.collection = d.client.Database("webCrawlerArchive").Collection("webpages")
		fmt.Println(d)
		filter := bson.D{{}}
		// Deletes all documents in the collection
		d.collection.DeleteMany(context.TODO(), filter)
	}
}
func (d *DatabaseConnection) disconnect() {
	if d.access {
		d.client.Disconnect(context.TODO())
	}
}
func (d *DatabaseConnection) insertWebpage(webpage WebPage) {
	if d.access {
		_, err := d.collection.InsertOne(context.TODO(), webpage)
		if err != nil {
			log.Printf("Failed to insert webpage: %v", err)
			return
		}
		return
	}
}
