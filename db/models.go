package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vendor struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `json:"name"`
	Address  string             `json:"address"`
	Location string             `json:"location"`
}

type Item struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `json:"name"`
	Price        float32            `json:"price"`
	Description  string             `json:"description"`
	Availability bool               `json:"availability"`
	VendorsList  []string           `bson:"vendorsList" json:"vendorsList"`
}
