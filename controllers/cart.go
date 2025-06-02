package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharukh010/go-ecommerce/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func AddToCart() gin.HandlerFunc{
	
}

func RemoveItem() gin.HandlerFunc{

}

func GetItemFromCart() gin.HandlerFunc{

}

func (app *Application) BuyFromCart() gin.HandlerFunc{
	return func(c *gin.Context){
		userQueryID := c.Query("id")
		if userQueryID == ""{
			c.JSON(http.StatusBadRequest,"UserID is required")
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
			c.JSON(http.StatusBadRequest,"User ID required")
			return 
		}
		productQueryID := c.Query("pid")
		if productQueryID == ""{
			c.JSON(http.StatusBadRequest,"Product ID required")
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