package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharukh010/go-ecommerce/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrProductIDRequired = errors.New("Product ID required")
	ErrUserIDRequired = errors.New("User ID required")


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

func (app *Application) GetItemFromCart() gin.HandlerFunc{

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