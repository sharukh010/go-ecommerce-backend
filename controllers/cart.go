package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharukh010/go-ecommerce/database"
	"github.com/sharukh010/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrProductIDRequired = errors.New("product ID required")
	ErrUserIDRequired = errors.New("user ID required")


)
type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection,userCollection *mongo.Collection) *Application{
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc{
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == ""{
			c.JSON(http.StatusBadRequest,ErrProductIDRequired.Error())
			return 
		}
		userQueryID := c.Query("userid")
		if userQueryID == ""{
			c.JSON(http.StatusBadRequest,ErrUserIDRequired.Error())
			return 
		}
		productID,err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,err.Error())
			return 
		}
		var ctx,cancel = context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx,app.prodCollection,app.userCollection,productID,userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,err.Error())
			return 
		}
		c.JSON(http.StatusOK,"Successfully Added to the cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc{
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == ""{
			c.JSON(http.StatusBadRequest,ErrProductIDRequired.Error())
			return 
		}

		userQueryID := c.Query("userID")
		if userQueryID == ""{
			c.JSON(http.StatusBadRequest,ErrUserIDRequired.Error())
			return 
		}

		productID,err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,err.Error())
			return 
		}
		var ctx,cancel = context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()
		err = database.RemoveCartItem(ctx,app.userCollection,productID,userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,err.Error())
			return 
		}
		c.JSON(http.StatusOK,"Successfully removed from cart")
	}
}

func GetItemFromCart() gin.HandlerFunc{
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == ""{
			c.JSON(http.StatusBadRequest,ErrUserIDRequired.Error())
			return 
		}
		userID,err := primitive.ObjectIDFromHex(userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,err.Error())
			return 
		}

		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		var filledCart models.User
		err = UserCollection.FindOne(ctx,bson.M{"_id":userID}).Decode(&filledCart)
		if err != nil {
			c.JSON(http.StatusInternalServerError,err.Error())
			return 
		}
		filter_match := bson.D{{Key:"$match",Value:bson.D{primitive.E{Key:"_id",Value: userID}}}}
		unwind := bson.D{{Key:"$unwind",Value:bson.D{primitive.E{Key:"path",Value: "$usercart"}}}}
		grouping := bson.D{{Key:"$group",Value:bson.D{primitive.E{Key:"_id",Value: "$_id"},{Key:"total",Value:bson.D{primitive.E{Key:"$sum",Value:"$usercart.price"}}}}}}

		pointCursor,err := UserCollection.Aggregate(ctx,mongo.Pipeline{filter_match,unwind,grouping})

		if err != nil {
			log.Println(err.Error())
		}
		var listing []bson.M 
		if err = pointCursor.All(ctx,&listing); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError,"Error occured while fetching cart items")
		}

		for _,json := range listing{
			c.JSON(http.StatusOK,json["total"])
			c.JSON(http.StatusOK,filledCart.UserCart)
		}

	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc{
	return func(c *gin.Context){
		userQueryID := c.Query("id")
		if userQueryID == ""{
			c.JSON(http.StatusBadRequest,ErrUserIDRequired.Error())
			return 
		}
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx,app.userCollection,userQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,err.Error())
			return 
		}
		c.JSON(http.StatusOK,"Successfully placed the Order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc{
	return func(c *gin.Context){
		userQueryID := c.Query("userid")
		if userQueryID == ""{
			c.JSON(http.StatusBadRequest,ErrUserIDRequired.Error())
			return 
		}
		productQueryID := c.Query("pid")
		if productQueryID == ""{
			c.JSON(http.StatusBadRequest,ErrProductIDRequired.Error())
			return 
		}
		productID,err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,err.Error())
			return 
		}

		var ctx,cancel = context.WithTimeout(context.Background(),5*time.Second)
		defer cancel() 

		err = database.InstantBuyer(ctx,app.prodCollection,app.userCollection,productID,userQueryID);
		if err != nil {
			c.JSON(http.StatusBadRequest,err.Error())
			return 
		}
		c.JSON(http.StatusOK,"Successfully placed the order")
	}
}