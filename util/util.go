package dbUtil

import (
	"context"
	"ecommerce/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func GetModifiedResult(err error, updateResult *mongo.UpdateResult, message string) db.HTMLResponse {
	if err != nil {
		return Failure(err.Error())
	} else if updateResult.ModifiedCount == 1 && updateResult.MatchedCount == 1 {
		return Success("update success")
	} else if updateResult.ModifiedCount == 0 {
		return Failure("update failure")
	}
	return Failure(message)
}
