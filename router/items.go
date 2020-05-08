package router

import (
	"context"
	"ecommerce/db"
	"ecommerce/util"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func Item(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		{
			//add item
			result := addItem(req)
			json.NewEncoder(w).Encode(result)
			break
		}
	case http.MethodPost:
		{
			result := editItem(req)
			json.NewEncoder(w).Encode(result)
			break
		}
	case http.MethodGet:
		{
			itemsList := getAllItems()
			json.NewEncoder(w).Encode(itemsList)
			break
		}
	case http.MethodDelete:
		{
			result := deleteItem(req)
			json.NewEncoder(w).Encode(result)
			break
		}
	default:
		{
			fmt.Fprintf(w, "404 page not found")
		}
	}
}

func addItem(req *http.Request) db.HTMLResponse {
	newItem := db.Item{}
	err := json.NewDecoder(req.Body).Decode(&newItem)

	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if newItem.Name == "" {
		return dbUtil.Failure("Item Name is required")
	}

	insertResult, err := db.GetConnection().Collection("items").InsertOne(context.TODO(), newItem)
	if err != nil {
		return dbUtil.Failure(err.Error())
	}

	if str, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		newItem.ID = str
	}

	// returning inserted document id in string format
	return dbUtil.Success(newItem.ID.Hex())
}

func editItem(req *http.Request) db.HTMLResponse {
	newItem := db.Item{}
	err := json.NewDecoder(req.Body).Decode(&newItem)

	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if newItem.ID.IsZero() {
		return dbUtil.Failure("Item Id is Required")
	} else if newItem.Name == "" {
		return dbUtil.Failure("Item name is Required")
	}

	filter := bson.D{{"_id", newItem.ID}}
	updateQuery := bson.D{{
		"$set", bson.D{
			{"name", newItem.Name},
			{"price", newItem.Price},
			{"description", newItem.Description},
		},
	}}

	updateResult, err := db.GetConnection().Collection("items").UpdateOne(context.TODO(), filter, updateQuery)
	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if updateResult.ModifiedCount == 1 && updateResult.MatchedCount == 1 {
		return dbUtil.Success("update success")
	} else if updateResult.ModifiedCount == 0 {
		return dbUtil.Failure("update failure")
	}
	return dbUtil.Failure("Item Not Found")

}

func deleteItem(req *http.Request) db.HTMLResponse {
	itemId, _ := req.URL.Query()["itemId"]

	objID, err := primitive.ObjectIDFromHex(itemId[0])
	if err != nil {
		return dbUtil.Failure(err.Error())
	}

	resultDelete, err := db.GetConnection().Collection("items").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return dbUtil.Failure(err.Error())
	} else if resultDelete.DeletedCount == 1 {
		return dbUtil.Success("deletion success")
	}
	return dbUtil.Failure("deletion failure")
}

func getAllItems() []db.Item {
	cur, _ := db.GetConnection().Collection("items").Find(context.Background(), bson.D{})
	defer cur.Close(context.Background())

	eachItem := db.Item{}
	ItemsList := []db.Item{}

	for cur.Next(context.Background()) {
		cur.Decode(&eachItem)
		ItemsList = append(ItemsList, eachItem)
	}

	return ItemsList
}
