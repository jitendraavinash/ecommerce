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

func Vendor(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		{
			vendorsList := getAllVendors()
			json.NewEncoder(w).Encode(vendorsList)
			break
		}
	case http.MethodPost:
		{
			newVendor := editVendor(w, req)
			fmt.Fprintf(w, "Updated Vendor: %+v", newVendor)
			break
		}

	case http.MethodDelete:
		{
			newVendor := deleteVendor(w, req)
			fmt.Fprintf(w, "Deleted Vendor: %+v", newVendor)
			break
		}
	case http.MethodPut:
		{
			newVendor := addVendor(w, req)
			fmt.Fprintf(w, newVendor)
			break
		}
	default:
		{
			fmt.Fprintf(w, "404 page not found")
		}
	}
}

func getAllVendors() []bson.M {

	addToFields := bson.D{{"$addFields", bson.D{{"_id", bson.M{"$toString": "$_id"}}}}}
	ProjectionStage := bson.D{{"$project", bson.D{{"itemList.vendorsList", 0}}}}
	lookupStage := bson.D{{
		"$lookup", bson.D{
			{"from", "items"},
			{"localField", "_id"},
			{"foreignField", "vendorsList"},
			{"as", "itemList"},
		},
	}}
	showInfoCursor, err := db.GetConnection().Collection("vendors").Aggregate(context.TODO(), mongo.Pipeline{addToFields, lookupStage, ProjectionStage})
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(context.Background(), &showsWithInfo); err != nil {
		panic(err)
	}
	fmt.Println(showsWithInfo)

	return showsWithInfo
}

func addVendor(w http.ResponseWriter, req *http.Request) string {
	newVendor := db.Vendor{}
	err := json.NewDecoder(req.Body).Decode(&newVendor)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "database connection error"
	} else if newVendor.Location == "" {
		return "Location not present"
	} else if newVendor.Name == "" {
		return "Location not present"
	} else if newVendor.Address == "" {
		newVendor.Address = ""
	}

	if dbUtil.FindVendorByName(newVendor.Name) {
		return "duplicate vendor"
	} else {
		insertResult, err := db.GetConnection().Collection("vendors").InsertOne(context.TODO(), newVendor)
		if err != nil {
			return "database connection error"
		}
		if str, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
			newVendor.ID = str
		}
		return "Inserted Vendor with ID: " + newVendor.ID.Hex()
	}
}

func editVendor(w http.ResponseWriter, req *http.Request) string {
	newVendor := db.Vendor{}
	err := json.NewDecoder(req.Body).Decode(&newVendor)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "database connection error"
	} else if newVendor.ID.IsZero() {
		return "Vendor Id is Required"
	} else if newVendor.Location == "" {
		return "Location is Required"
	} else if newVendor.Name == "" {
		return "Vendor Name is Required"
	} else if newVendor.Address == "" {
		newVendor.Address = ""
	}

	if dbUtil.FindVendorById(newVendor.ID) {
		filter := bson.D{{"_id", newVendor.ID}}
		updateQuery := bson.D{{
			"$set", bson.D{
				{"name", newVendor.Name},
				{"location", newVendor.Location},
				{"address", newVendor.Address},
			},
		}}

		updateResult, err := db.GetConnection().Collection("vendors").UpdateOne(context.TODO(), filter, updateQuery)
		if err != nil {
			return "database connection error"
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
		return "update success"
	} else {
		return "Invalid vendor"
	}

}

func deleteVendor(w http.ResponseWriter, req *http.Request) bool {
	vendorId, ok := req.URL.Query()["vendorId"]
	fmt.Println(vendorId[0], ok)

	// Delete one document.
	objID, err := primitive.ObjectIDFromHex(vendorId[0])
	if err != nil {
		return false
	}

	resultDelete, err := db.GetConnection().Collection("vendors").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(resultDelete.DeletedCount)
	return true
}
