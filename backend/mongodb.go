package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

/*func postTrack(track Track) {
	c := getClient()
	c.Locked = true
	coll := s.getCollection()
	err := coll.Insert(&Track{track.Id, track.Name, track.Author, track.Steps})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("postTrack session UUID:", s.Uuid)
	c.Locked = false
}*/

func getTracks() []Track {
	t := []Track{}

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("stl").Collection("tracks")
	filter := bson.D{{"id", bson.D{{"$gte", 1}}}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)
	}

	return t
}

func getTrack(id int) Track {
	t := Track{}

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var result bson.M
	coll := client.Database("stl").Collection("tracks")
	err = coll.FindOne(context.TODO(), bson.D{{"id", 1}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found \n")
	}
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

	return t
}
