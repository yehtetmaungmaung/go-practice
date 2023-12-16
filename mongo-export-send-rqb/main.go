package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://root:Fr0nt!!r@localhost:27017"

type Feature struct {
	Properties Properties `bson:"properties"`
	Geometry   Geometry   `bson:"geometry"`
}

type Geometry struct {
	Type        string     `bson:"type"`
	Coordinates [2]float64 `bson:"coordinates"`
}

type AdditionalInfo struct {
	Township         string `bson:"township"`
	ParentID         string `bson:"parent_id"`
	MaxPorts         string `bson:"max_ports"`
	PortAvailability string `bson:"port_availability"`
}

type Properties struct {
	Type           string         `bson:"type"`
	ID             string         `bson:"id"`
	AdditionalInfo AdditionalInfo `bson:"additional_info"`
	Tags           []string       `bson:"tags"`
}

type Rbqmsg struct {
	Id        string `json:"node_id"`
	EventType string `json:"event_type"`
	Value     `json:"values"`
}

type RbqUpdateMsg struct {
	NodeID      string      `json:"node_id"`
	EventType   string      `json:"event_type"`
	UpdateValue UpdateValue `json:"values"`
}

type UpdateValue struct {
	NodeName         string   `json:"node_name"`
	Lat              float64  `json:"lat"`
	Long             float64  `json:"long"`
	Type             string   `json:"type"`
	MaxPorts         string   `json:"max_port"`
	PortAvailability string   `json:"port_availability"`
	Township         string   `json:"township"`
	ParentName       string   `json:"parent_name"`
	Tags             []string `json:"tags"`
	ConnectionType   string   `json:"connection_type"`
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
	filter := bson.D{{Key: "properties.type", Value: "CA2"}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())

	var count int
	var rbqmsgs []Rbqmsg
	var rbqUpdtMsgs []RbqUpdateMsg

	for cursor.Next(context.TODO()) {
		if count > 5 {
			break
		}
		var feature Feature
		if err := cursor.Decode(&feature); err != nil {
			log.Fatal(err)
		}

		fmt.Println(feature.Properties.AdditionalInfo.MaxPorts)
		// Create an Rbqmsg instance with appropriate data
		rbqmsg := Rbqmsg{
			Id:        feature.Properties.ID,
			EventType: "node_delete",
			Value:     Value{Type: feature.Properties.Type}, // You can set the appropriate type here
		}

		rbqUpdtMsg := RbqUpdateMsg{
			NodeID:    feature.Properties.ID,
			EventType: "node_update",
			UpdateValue: UpdateValue{
				NodeName:         feature.Properties.ID,
				Lat:              feature.Geometry.Coordinates[0],
				Long:             feature.Geometry.Coordinates[1],
				Type:             feature.Properties.Type,
				MaxPorts:         feature.Properties.AdditionalInfo.MaxPorts,
				PortAvailability: feature.Properties.AdditionalInfo.PortAvailability,
				Township:         feature.Properties.AdditionalInfo.Township,
				ParentName:       feature.Properties.AdditionalInfo.ParentID,
				Tags:             feature.Properties.Tags,
				ConnectionType:   "",
			},
		}

		rbqmsgs = append(rbqmsgs, rbqmsg)
		rbqUpdtMsgs = append(rbqUpdtMsgs, rbqUpdtMsg)

		count++
	}

	// Create a JSON file for writing
	jsonFile, err := os.Create("delete-fibermaps-200.json")
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

	updateFile, err := os.Create("update-fibermaps-200.json")
	if err != nil {
		log.Fatal(err)
	}
	defer updateFile.Close()
	jsonData, err = json.Marshal(rbqUpdtMsgs)
	if err != nil {
		log.Fatal(err)
	}
	_, err = updateFile.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}

}
