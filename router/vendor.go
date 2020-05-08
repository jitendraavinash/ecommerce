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
	case http.MethodPut:
		{
			// add vendor
			result := addVendor(req)
			json.NewEncoder(w).Encode(result)
			break
		}
	case http.MethodPost:
		{
			// edit vendor
			result := editVendor(req)
			json.NewEncoder(w).Encode(result)
			break
		}
	case http.MethodGet:
		{
			vendorsList := getAllVendors()
			json.NewEncoder(w).Encode(vendorsList)
			break
		}
	case http.MethodDelete:
		{
			result := deleteVendor(req)
			json.NewEncoder(w).Encode(result)
			break
		}
	default:
		{
			fmt.Fprintf(w, "404 page not found")
		}
	}
}

func addVendor(req *http.Request) db.HTMLResponse {
	newVendor := db.Vendor{}
	err := json.NewDecoder(req.Body).Decode(&newVendor)

	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if newVendor.Name == "" {
		return dbUtil.Failure("Vendor Name is Required")
	}
	// to make sure no entry gets to database while creating a vendor
	newVendor.ItemsList = []db.VendorItem{}
	insertResult, err := db.GetConnection().Collection("vendors").InsertOne(context.TODO(), newVendor)
	if err != nil {
		return dbUtil.Failure(err.Error())
	}
	if str, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		newVendor.ID = str
	}
	return dbUtil.Success(newVendor.ID.Hex())
}

func editVendor(req *http.Request) db.HTMLResponse {
	newVendor := db.Vendor{}
	err := json.NewDecoder(req.Body).Decode(&newVendor)

	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if newVendor.ID.IsZero() {
		return dbUtil.Failure("Vendor Id is Required")
	} else if newVendor.Name == "" {
		return dbUtil.Failure("Vendor Name is Required")
	}

	filter := bson.D{{"_id", newVendor.ID}}
	updateQuery := bson.D{{
		"$set", bson.D{
			{"name", newVendor.Name},
		},
	}}

	updateResult, err := db.GetConnection().Collection("vendors").UpdateOne(context.TODO(), filter, updateQuery)
	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if updateResult.ModifiedCount == 1 && updateResult.MatchedCount == 1 {
		return dbUtil.Success("update success")
	} else if updateResult.ModifiedCount == 0 {
		return dbUtil.Failure("update failure")
	}
	return dbUtil.Failure("Vendor Not Found")

}

func deleteVendor(req *http.Request) db.HTMLResponse {
	vendorId, _ := req.URL.Query()["vendorId"]

	// Delete one document.
	objID, err := primitive.ObjectIDFromHex(vendorId[0])
	if err != nil {
		return dbUtil.Failure(err.Error())
	}

	resultDelete, err := db.GetConnection().Collection("vendors").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if resultDelete.DeletedCount == 1 {
		return dbUtil.Success("deletion success")
	}
	return dbUtil.Failure("deletion failure")
}

func getAllVendors() []bson.M {
	unWindStage := bson.D{
		{"$unwind", bson.D{
			{"path", "$itemsList"},
			{"preserveNullAndEmptyArrays", true},
		}},
	}
	lookupStage := bson.D{{
		"$lookup", bson.D{
			{"from", "items"},
			{"localField", "itemsList.itemId"},
			{"foreignField", "_id"},
			{"as", "itemsList"},
		},
	}}
	showInfoCursor, err := db.GetConnection().Collection("vendors").Aggregate(context.TODO(), mongo.Pipeline{unWindStage, lookupStage})
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(context.Background(), &showsWithInfo); err != nil {
		panic(err)
	}
	fmt.Println(showsWithInfo)

	return showsWithInfo
}
