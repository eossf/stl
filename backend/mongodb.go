package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

func postTrack(track Track) {
	// -------------------------------------------------------------------
	// ------------- mongodb session -------------------------------------
	server := os.Getenv("MONGODB_HOST")
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	session.Login(&mgo.Credential{Username: os.Getenv("MONGODB_USER"), Password: os.Getenv("MONGODB_PASSWORD"), Source: "stl"})
	c := session.DB("stl").C("track")
	// -------------------------------------------------------------------
	// -------------------------------------------------------------------

	err = c.Insert(&Track{track.Id, track.Name, track.Author, track.Steps})
	if err != nil {
		log.Fatal(err)
	}
}

func getTracks() []Track {
	// -------------------------------------------------------------------
	// ------------- mongodb session -------------------------------------
	server := os.Getenv("MONGODB_HOST")
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	session.Login(&mgo.Credential{Username: os.Getenv("MONGODB_USER"), Password: os.Getenv("MONGODB_PASSWORD"), Source: "stl"})
	c := session.DB("stl").C("track")
	// -------------------------------------------------------------------
	// -------------------------------------------------------------------
	result := []Track{}
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func getTrack(id string) Track {
	// -------------------------------------------------------------------
	// ------------- mongodb session -------------------------------------
	server := os.Getenv("MONGODB_HOST")
	session, err := mgo.Dial(server)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	session.Login(&mgo.Credential{Username: os.Getenv("MONGODB_USER"), Password: os.Getenv("MONGODB_PASSWORD"), Source: "stl"})
	c := session.DB("stl").C("track")
	// -------------------------------------------------------------------
	// -------------------------------------------------------------------
	result := Track{}
	err = c.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func testMongodb() {
	getTrack("1")
	log.Println("Mongodb OK server: ", os.Getenv("MONGODB_HOST"))
}
