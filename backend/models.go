package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CreatedAt time.Time          `json:"created_at"`
//UpdatedAt time.Time          `json:"updated_at"`
//Text      string             `json:"text"`

type Track struct {
	_ID    primitive.ObjectID `json:"_id"`
	Id     int                `json:"id"`
	Name   string             `json:"name"`
	Author string             `json:"author"`
	Steps  int                `json:"steps"`
}
