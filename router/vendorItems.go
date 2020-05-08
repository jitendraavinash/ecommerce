package router

import (
	"context"
	"ecommerce/db"
	"ecommerce/util"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

func VendorItems(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		{
			// add an Item to Vendor List
			result := addItemToVendor(req)
			json.NewEncoder(w).Encode(result)
			break
		}
	case http.MethodPost:
		{
			// edit an Item in Vendor List
			result := editItemInVendor(req)
			json.NewEncoder(w).Encode(result)
			break
		}
	case http.MethodDelete:
		{
			// delete an Item in Vendor List
			result := deleteItemInVendor(req)
			json.NewEncoder(w).Encode(result)
			break
		}
	case http.MethodGet:
		{
			vendorItemList := getVendorItemList(req)
			json.NewEncoder(w).Encode(vendorItemList)
			break
		}
	default:
		{
			fmt.Fprintf(w, "404 page not found")
		}
	}
}

func addItemToVendor(req *http.Request) db.HTMLResponse {
	newVendor := db.VendorItemReqBody{}
	err := json.NewDecoder(req.Body).Decode(&newVendor)

	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if newVendor.VendorId.IsZero() {
		return dbUtil.Failure("Vendor Id is Required")
	} else if newVendor.ItemId.IsZero() {
		return dbUtil.Failure("Item Id is Required")
	}

	if dbUtil.FindItemById(newVendor.ItemId) {

		matchCondition := bson.D{
			{"_id", newVendor.VendorId},
			{
				"itemsList.itemId", bson.D{
					{"$nin", [5]primitive.ObjectID{newVendor.ItemId}},
				},
			},
		}
		updateQuery := bson.D{{
			"$push", bson.D{
				{"itemsList", bson.D{
					{"itemId", newVendor.ItemId},
					{"availability", newVendor.Availability},
				}},
			},
		}}

		updateResult, err := db.GetConnection().Collection("vendors").UpdateOne(context.TODO(), matchCondition, updateQuery)
		fmt.Println(updateResult.MatchedCount, updateResult.ModifiedCount)
		if err != nil {
			return dbUtil.Failure(err.Error())
		} else if updateResult.ModifiedCount == 1 && updateResult.MatchedCount == 1 {
			return dbUtil.Success("update success")
		}
		return dbUtil.Failure("Item Already present")
	}
	return dbUtil.Failure("Item is not available")

}

func editItemInVendor(req *http.Request) db.HTMLResponse {
	newVendor := db.VendorItemReqBody{}
	err := json.NewDecoder(req.Body).Decode(&newVendor)

	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if newVendor.VendorId.IsZero() {
		return dbUtil.Failure("Vendor Id is Required")
	} else if newVendor.ItemId.IsZero() {
		return dbUtil.Failure("Item Id is Required")
	}

	if dbUtil.FindItemById(newVendor.ItemId) {

		matchCondition := bson.D{
			{"_id", newVendor.VendorId},
			{
				"itemsList", bson.D{
					{"$elemMatch", bson.D{
						{"itemId", newVendor.ItemId},
					}},
				},
			},
		}
		updateQuery := bson.D{{
			"$set", bson.D{
				{"itemsList.$.availability", newVendor.Availability},
			},
		}}

		updateResult, err := db.GetConnection().Collection("vendors").UpdateOne(context.TODO(), matchCondition, updateQuery)
		fmt.Println(updateResult.MatchedCount, updateResult.ModifiedCount)
		if err != nil {
			return dbUtil.Failure(err.Error())
		} else if updateResult.ModifiedCount == 1 && updateResult.MatchedCount == 1 {
			return dbUtil.Success("update success")
		}
		return dbUtil.Failure("Item Already present")
	}
	return dbUtil.Failure("Item is not available")

}

func deleteItemInVendor(req *http.Request) db.HTMLResponse {
	vendorId, _ := req.URL.Query()["vendorId"]
	itemId, _ := req.URL.Query()["itemId"]

	// Delete one document.
	vendorObjID, _ := primitive.ObjectIDFromHex(vendorId[0])
	itemObjID, _ := primitive.ObjectIDFromHex(itemId[0])
	fmt.Println(vendorObjID, itemObjID)
	if vendorObjID.IsZero() {
		return dbUtil.Failure("Vendor Id is Required")
	} else if itemObjID.IsZero() {
		return dbUtil.Failure("Item Id is Required")
	}

	matchCondition := bson.D{{"_id", vendorObjID}}
	updateQuery := bson.D{{
		"$pull", bson.D{
			{"itemsList", bson.D{
				{"itemId", itemObjID},
			}},
		},
	}}

	updateResult, err := db.GetConnection().Collection("vendors").UpdateOne(context.TODO(), matchCondition, updateQuery)
	fmt.Println(updateResult.MatchedCount, updateResult.ModifiedCount)
	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if updateResult.ModifiedCount == 1 && updateResult.MatchedCount == 1 {
		return dbUtil.Success("update success")
	}
	return dbUtil.Failure("Item Already present")
}

func getVendorItemList(req *http.Request) db.HTMLResponse {
	itemId, _ := req.URL.Query()["itemId"]
	itemObjID, err := primitive.ObjectIDFromHex(itemId[0])
	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if itemObjID.IsZero() {
		return dbUtil.Failure("Item Id is Required")
	}

	findQuery := bson.D{
		{
			"itemsList", bson.D{
				{"$elemMatch", bson.D{
					{"itemId", itemObjID},
				}},
			},
		},
	}
	findOptions := options.Find()
	findOptions.SetProjection(bson.D{{"itemsList", 0}})
	cur, err := db.GetConnection().Collection("vendors").Find(context.TODO(), findQuery, findOptions)
	if err != nil {
		dbUtil.Failure(err.Error())
	}

	var results []*db.Vendor
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem db.Vendor
		cur.Decode(&elem)
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		dbUtil.Failure(err.Error())
	}
	cur.Close(context.TODO())

	out, _ := json.Marshal(results)

	return dbUtil.Success(string(out))
}
