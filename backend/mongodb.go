package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

var clients = make([]*Client, 0, MAX)

const MAX = 99

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

/*func getTracks() []Track {
	t := []Track{}
	c := getClient()
	c.Locked = true
	var result bson.M
	coll := c.MongoClient.Database("stl").Collection("tracks")
	err := coll.FindOne(context.TODO(), bson.D{{"id", 1}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found \n")
		return t
	}
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

	log.Println("getTracks session UUID:", c.Uuid)
	//c.Locked = false
	return t
}*/

func getTrack(id int) Track {
	t := Track{}
	client := getClient(os.Getenv("MONGODB_URI"))

	//c.Locked = true
	var result bson.M
	coll := client.Database("stl").Collection("tracks") //c.MongoClient.Database("stl").Collection("tracks")
	coll.FindOne(context.TODO(), bson.D{{"id", 1}}).Decode(&result)
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

	//log.Println("getTracks session UUID:", c.Uuid)
	//c.Locked = false
	return t
}
