package dbUtil

import (
	"context"
	"ecommerce/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindVendorByName(vendorName string) bool {
	newVendor := db.Vendor{}
	db.GetConnection().Collection("vendors").FindOne(context.TODO(), bson.D{{"name", vendorName}}).Decode(&newVendor)

	// return true if a document is present and
	// false if the vendor is new vendor
	return false
}

func FindVendorById(vendorId primitive.ObjectID) bool {
	vendorFromDB := db.Vendor{}
	db.GetConnection().Collection("vendors").FindOne(context.TODO(), bson.D{{"_id", vendorId}}).Decode(&vendorFromDB)

	// return true if a document is present and
	// false if the vendor is new vendor
	return false
}

func FindItemById(itemId primitive.ObjectID) bool {
	itemFromDB := db.Item{}
	db.GetConnection().Collection("items").FindOne(context.TODO(), bson.D{{"_id", itemId}}).Decode(&itemFromDB)

	// return true if a document is present and
	// false if the vendor is new vendor
	if itemFromDB.ID.IsZero() {
		return false
	} else {
		return true
	}
}

func Success(message string) db.HTMLResponse {
	response := db.HTMLResponse{}
	response.Error = false
	response.Message = message
	return response
}

func Failure(message string) db.HTMLResponse {
	response := db.HTMLResponse{}
	response.Error = true
	response.Message = message
	return response
}
