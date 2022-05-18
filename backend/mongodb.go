package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func putTrack(track Track) {
	log.Println("putTrack")
	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err)
	}

	coll := client.Database("stl").Collection("tracks")
	filter := bson.D{{"id", track.Id}}
	replacement := bson.D{{"id", track.Id}, {"name", track.Name}, {"author", track.Author}, {"steps", track.Steps}}
	_, err = coll.ReplaceOne(context.TODO(), filter, replacement)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func postTrack(track Track) {
	log.Println("postTrack")
	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err)
	}

	coll := client.Database("stl").Collection("tracks")
	_, err = coll.InsertOne(context.TODO(), &Track{Id: track.Id, Name: track.Name, Author: track.Author, Steps: track.Steps})
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func getTracks() []Track {
	log.Println("getTracks")
	t := []Track{}

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err)
	}

	coll := client.Database("stl").Collection("tracks")
	filter := bson.D{{"id", bson.D{{"$gte", 0}}}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err)
	}

	if err = cursor.All(context.TODO(), &t); err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	return t
}

func getTrack(id int) Track {
	log.Println("getTrack ", id)
	t := Track{}

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err)
	}

	coll := client.Database("stl").Collection("tracks")
	err = coll.FindOne(context.TODO(), bson.D{{"id", id}}, options.FindOne().SetShowRecordID(true)).Decode(&t)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	return t
}
