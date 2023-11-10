package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://root:Fr0nt!!r@localhost:27017"

type Feature struct {
	Properties
}

type Properties struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Rbqmsg struct {
	Id        string `json:"node_id"`
	EventType string `json:"event_type"`
	Value     `json:"values"`
}

type Value struct {
	Type string `json:"type"`
}

func main() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("mapDev").Collection("features_access_fttx")
	filter := bson.D{{Key: "properties.type", Value: "CPE"}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())

	var count int
	var rbqmsgs []Rbqmsg

	for cursor.Next(context.TODO()) {
		if count > 20000 {
			break
		}
		var feature Feature
		if err := cursor.Decode(&feature); err != nil {
			log.Fatal(err)
		}

		// Create an Rbqmsg instance with appropriate data
		rbqmsg := Rbqmsg{
			Id:        feature.Properties.ID,
			EventType: "node_delete",
			Value:     Value{Type: feature.Properties.Type}, // You can set the appropriate type here
		}

		rbqmsgs = append(rbqmsgs, rbqmsg)

		count++
	}

	// Create a JSON file for writing
	jsonFile, err := os.Create("delete-fibermaps.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// Marshal the slice of Rbqmsg into a JSON array and write it to the file
	jsonData, err := json.Marshal(rbqmsgs)
	if err != nil {
		log.Fatal(err)
	}
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
}
