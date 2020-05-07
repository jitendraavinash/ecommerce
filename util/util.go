package dbUtil

import (
	"context"
	"ecommerce/db"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindVendorByName(vendorName string) bool {
	newVendor := db.Vendor{}
	db.GetConnection().Collection("vendors").FindOne(context.TODO(), bson.D{{"name", vendorName}}).Decode(&newVendor)

	// return true if a document is present and
	// false if the vendor is new vendor
	if (db.Vendor{}) != newVendor {
		return true
	} else {
		return false
	}
}

func FindVendorById(vendorId primitive.ObjectID) bool {
	vendorFromDB := db.Vendor{}
	db.GetConnection().Collection("vendors").FindOne(context.TODO(), bson.D{{"_id", vendorId}}).Decode(&vendorFromDB)

	// return true if a document is present and
	// false if the vendor is new vendor
	if (db.Vendor{}) != vendorFromDB {
		return true
	} else {
		return false
	}
}

func FindItemByName(itemName string) bool {
	newItem := db.Item{}
	db.GetConnection().Collection("items").FindOne(context.TODO(), bson.D{{"name", itemName}}).Decode(&newItem)

	// return true if item with same name is present and
	// false if the item is not present
	if newItem.ID.IsZero() {
		return false
	} else {
		return true
	}
}

func FindItemById(itemId primitive.ObjectID) bool {
	itemFromDB := db.Item{}
	db.GetConnection().Collection("items").FindOne(context.TODO(), bson.D{{"_id", itemId}}).Decode(&itemFromDB)

	// return true if a document is present and
	// false if the vendor is new vendor

	// Item id is zero so that means no item is present

	fmt.Println(itemFromDB)
	if itemFromDB.ID.IsZero() {
		return false
	} else {
		return true
	}
}
