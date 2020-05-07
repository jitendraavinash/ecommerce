package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VendorItem struct {
	ItemId       primitive.ObjectID `bson:"itemId" json:"itemId"`
	Availability bool               `json:"availability"`
}

type Vendor struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `json:"name"`
	ItemsList []VendorItem       `bson:"itemsList" json:"itemsList"`
}

type Item struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `json:"name"`
	Price       float32            `json:"price"`
	Description string             `json:"description"`
}

type VendorItemReqBody struct {
	VendorId     primitive.ObjectID `bson:"vendorId" json:"vendorId"`
	ItemId       primitive.ObjectID `bson:"itemId" json:"itemId"`
	Availability bool               `json:"availability"`
}

type HTMLResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
