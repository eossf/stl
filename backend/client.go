package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Client struct {
	Uri         string
	MongoClient *mongo.Client
	Err         error
	Uuid        string
	Locked      bool
}

//func getClient() *Client {
//	// get next available client
//	countClients := 0
//	countLockedClients := 0
//	addClient := true
//	for _, c := range clients {
//		if c != nil {
//			if c.Locked {
//				countLockedClients++
//			} else {
//				addClient = false
//				break
//			}
//			countClients++
//			addClient = true
//		}
//	}
//	if addClient && countClients > MAX {
//		log.Fatal("No more connection available ", countClients, "/", MAX, " reached")
//	}
//	if addClient && countClients <= MAX {
//		//MONGODB_URI=mongodb://stluser:stluser@localhost:27017/stl?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+1.3.1
//		c := NewClient(os.Getenv("MONGODB_URI"))
//		clients = append(clients, c)
//		log.Println("New Client created: ", c.Uuid)
//	} else {
//		log.Println("Use existing connection")
//	}
//	return clients[countClients]
//}

func NewClient(uri string) *Client {

	c := Client{
		Uri:    uri,
		Uuid:   uuid.NewString(),
		Err:    nil,
		Locked: false,
	}

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	c.MongoClient = client

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	return &c
}
