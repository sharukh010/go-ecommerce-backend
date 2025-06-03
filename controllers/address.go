package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharukh010/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("userid")
		if user_id == ""{
			c.JSON(http.StatusBadRequest,gin.H{"error":"User ID required"})
			return 
		}
		user_ID,err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return 
		}
		var address models.Address
		address.Address_id = primitive.NewObjectID()
		if err = c.BindJSON(&address); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return 
		}

		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel() 

		match_filter := bson.D{{Key:"$match",Value: bson.D{primitive.E{Key:"_id",Value: user_ID}}}}
		unwind := bson.D{{Key:"$unwind",Value:bson.D{primitive.E{Key:"path",Value: "$address"}}}}
		group := bson.D{{Key:"$group",Value: bson.D{primitive.E{Key:"_id",Value: "$address_id"},{Key:"count",Value: bson.D{primitive.E{Key: "$sum",Value: 1}}}}}}

		pointCursor,err := UserCollection.Aggregate(ctx,mongo.Pipeline{match_filter,unwind,group})

		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Something went wrong"})
			return 
		}

		var addressInfo []bson.M 
		if err = pointCursor.All(ctx,&addressInfo); err != nil {
			panic(err)
		}

		var size int32 
		for _,address_no := range addressInfo{
			count := address_no["count"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{primitive.E{Key:"_id",Value:user_ID}}
			update := bson.D{{Key:"$push",Value: bson.D{primitive.E{Key:"address",Value: address}}}}
			_,err := UserCollection.UpdateOne(ctx,filter,update)
			if err != nil {
				fmt.Println(err)
			}
		}else{
			c.JSON(http.StatusBadRequest,"More than 2 Addresses are Not allowed")
			return 
		}
		c.JSON(http.StatusOK,"Successfully added the Address")

	}
}

func EditHomeAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("userid")
		if user_id == ""{
			c.JSON(http.StatusNotFound,gin.H{"error":"User Id required"})
			return 
		}
		user_ID,err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		}
		var editHomeAddress models.Address
		if err := c.BindJSON(&editHomeAddress); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return 
		}
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key:"_id",Value: user_ID}}
		update := bson.D{{Key:"$set",Value : bson.D{primitive.E{Key:"address.0.house_name",Value:editHomeAddress.House},{Key: "address.0.street_name",Value: editHomeAddress.Street},{Key:"address.0.city_name",Value: editHomeAddress.City},{Key:"address.0.pin_code",Value: editHomeAddress.Pincode}}}}

		_,err = UserCollection.UpdateOne(ctx,filter,update)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Something Went Wrong"})
			return 
		}
		c.JSON(http.StatusOK,"Successfully Updated the Home Address")
	}

}

func EditWorkAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("userid")
		if user_id == ""{
			c.JSON(http.StatusNotFound,gin.H{"error":"User ID required"})
			return 
		}
		user_ID,err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return 
		}
		var editWorkAddress models.Address
		if err := c.BindJSON(&editWorkAddress); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return 
		}
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key:"_id",Value:user_ID}}
		update := bson.D{{Key:"$set",Value:bson.D{primitive.E{Key:"address.1.house_name",Value: editWorkAddress.House},{Key:"address.1.street_name",Value:editWorkAddress.Street},{Key:"address.1.city_name",Value: editWorkAddress.City},{Key:"address.1.pin_code",Value: editWorkAddress.Pincode}}}}
		_,err = UserCollection.UpdateOne(ctx,filter,update)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Something Went wrong"})
			return 
		}
		c.JSON(http.StatusOK,"Successfully Updated the Work Address")

	}
}

func DeleteAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("userid")
		if user_id == ""{
			c.JSON(http.StatusBadRequest,gin.H{"error":"User ID required"})
			c.Abort()
			return 
		}
		addresses := make([]models.Address,0)
		user_ID,err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,"Internal Server Error")
			return 
		}
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key:"_id",Value:user_ID}}
		update := bson.D{{Key:"$set",Value:bson.D{primitive.E{Key:"address",Value:addresses}}}}

		_,err = UserCollection.UpdateOne(ctx,filter,update)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound,"User not found with given ID")
			return 
		}
		c.IndentedJSON(http.StatusOK,"Successfully Deleted Addresses")

	}
}