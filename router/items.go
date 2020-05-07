package router

import (
	"context"
	"ecommerce/db"
	"ecommerce/util"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func Item(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		{
			itemsList := getAllItems()
			json.NewEncoder(w).Encode(itemsList)
			break
		}
	case http.MethodPost:
		{
			newItem := editItem(w, req)
			json.NewEncoder(w).Encode(newItem)
			break
		}

	case http.MethodDelete:
		{
			newItem := deleteItem(w, req)
			fmt.Fprintf(w, "Deleted Item: %+v", newItem)
			break
		}
	case http.MethodPut:
		{
			newItem := addItem(w, req)
			fmt.Fprintf(w, "Inserted Item: %+v", newItem)
			break
		}
	default:
		{
			fmt.Fprintf(w, "404 page not found")
		}
	}
}

func getAllItems() []bson.M {
	unWindStage := bson.D{{"$unwind", "$vendorsList"}}
	addToFields := bson.D{{"$addFields", bson.D{{"vendorsList", bson.D{{"$toObjectId", "$vendorsList"}}}}}}
	lookupStage := bson.D{{
		"$lookup", bson.D{
			{"from", "vendors"},
			{"localField", "vendorsList"},
			{"foreignField", "_id"},
			{"as", "vendorsList"},
		},
	}}
	showInfoCursor, err := db.GetConnection().Collection("items").Aggregate(context.TODO(), mongo.Pipeline{unWindStage, addToFields, lookupStage})

	var showsWithInfo []bson.M
	if err = showInfoCursor.All(context.Background(), &showsWithInfo); err != nil {
		panic(err)
	}
	return showsWithInfo
}

func addItem(w http.ResponseWriter, req *http.Request) string {
	newItem := db.Item{}
	err := json.NewDecoder(req.Body).Decode(&newItem)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "database connection error"
	} else if newItem.Name == "" {
		return "Item Name is required"
	} else if newItem.Description == "" {
		newItem.Description = ""
	}

	if dbUtil.FindItemByName(newItem.Name) {
		return "duplicate item"
	} else {
		insertResult, err := db.GetConnection().Collection("items").InsertOne(context.TODO(), newItem)
		if err != nil {
			return "database connection error"
		}

		if str, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
			newItem.ID = str
		}

		// returning inserted document id in string format
		return newItem.ID.Hex()
	}
}

func editItem(w http.ResponseWriter, req *http.Request) string {
	newItem := db.Item{}
	err := json.NewDecoder(req.Body).Decode(&newItem)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "database connection error"
	} else if newItem.ID.IsZero() {
		return "Vendor Id is Required"
	} else if newItem.Name == "" {
		return "Item name is Required"
	}

	if dbUtil.FindItemById(newItem.ID) {
		filter := bson.D{{"_id", newItem.ID}}
		updateQuery := bson.D{{
			"$set", bson.D{
				{"name", newItem.Name},
				{"price", newItem.Price},
				{"description", newItem.Description},
				{"availability", newItem.Availability},
				{"vendorsList", newItem.VendorsList},
			},
		}}

		updateResult, err := db.GetConnection().Collection("items").UpdateOne(context.TODO(), filter, updateQuery)
		if err != nil {
			return "database connection error"
		}
		if updateResult.ModifiedCount == 1 {
			return "update success"
		} else {
			return "update failure"
		}
	} else {
		return "Invalid Item"
	}

}

func deleteItem(w http.ResponseWriter, req *http.Request) bool {
	itemId, ok := req.URL.Query()["itemId"]
	fmt.Println(itemId[0], ok)

	// Delete one document.
	objID, err := primitive.ObjectIDFromHex(itemId[0])
	if err != nil {
		fmt.Println(err)
		return false
	}

	resultDelete, err := db.GetConnection().Collection("items").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(resultDelete.DeletedCount)
	return true
}
