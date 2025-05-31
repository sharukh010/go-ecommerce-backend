package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharukh010/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc{

}

func EditHomeAddress() gin.HandlerFunc{

}

func EditWorkAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")
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
		user_id := c.Query("id")
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
		c.IndentedJSON(http.StatusOK,"Successfully Deleted Address")

	}
}