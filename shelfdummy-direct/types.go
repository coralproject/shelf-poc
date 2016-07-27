package main

import (
	"github.com/cayleygraph/cayley/quad"
	"gopkg.in/mgo.v2/bson"
)

// Item is used to encode coral item data.
type Item struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Type    string        `bson:"t" json:"t"`
	Version int           `bson:"v" json:"v"`
	Data    ItemData      `bson:"d" json:"d"`
}

// ItemData is used to encode some nested dummy data.
type ItemData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Body string `json:"body,omitempty"`
}

// Job is used to encode data to be sent to sponge.
type Job struct {
	Data  Item
	Type  string
	Quads []quad.Quad
}
