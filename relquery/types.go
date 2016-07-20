package main

import "gopkg.in/mgo.v2/bson"

// ItemData is what an Item can hold
//  Should be the intersection of the db and transport protocols supported
type ItemData interface{}

// An Item is either an asset, comment, or user.
type Item struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Type    string        `bson:"t" json:"t"` // ItemType.Name
	Version int           `bson:"v" json:"v"`
	Data    ItemData      `bson:"d" json:"d"`
	Rels    []Rel         `bson:"rels,omitempty" json:"rels,omitempty"`
}

// Rel holds an item's relationship to another item.
type Rel struct {
	Name string `bson:"n" json:"n"`   // Name of relationship
	Type string `bson:"t" json:"t"`   // Item Type of target
	ID   string `bson:"id" json:"id"` // Id of target
}

// custom error type for bases that can't be found
type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
