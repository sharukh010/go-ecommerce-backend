package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/sharukh010/go-ecommerce/database"
	"github.com/sharukh010/go-ecommerce/models"
)

var UserCollection *mongo.Collection = database.UserData(database.Client,"Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client,"Products")
var Validate = validator.New()

func HashPassword(password string) (string,error) {
	bytes,err := bcrypt.GenerateFromPassword([]byte(password),14)

	if err != nil {
		return "",err
	}
	
	return string(bytes),nil 
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

}

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		var user models.User 
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return 
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
			return 
		}

		count,err := UserCollection.CountDocuments(ctx,bson.M{"email":user.Email})

		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return 
		}

		if count > 0{
			c.JSON(http.StatusBadRequest,gin.H{"error":"User already exists"})
		}

		count,err = UserCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})

		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err})
			return 
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest,gin.H{"error":"Phone number all ready used"})
			return 
		}

		password,err := HashPassword(*user.Password)
		if err!= nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return 
		}
		user.Password = &password

		user.Created_At,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.Updated_At,_ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		/*
		need to write token logic
		
		*/

		user.UserCart = make([]models.ProductUser,0)
		user.Address_Details = make([]models.Address,0)
		user.Order_Status = make([]models.Order,0)
		_,insertErr := UserCollection.InsertOne(ctx,user);
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":insertErr.Error()})
			return 
		}

		c.JSON(http.StatusCreated,"Successfully Signed Up!!")


	}
}

func Login() gin.HandlerFunc{

}

func ProductViewerAdmin() gin.HandlerFunc{

}

func SearchProduct() gin.HandlerFunc{

}

func SearchProductByQuery() gin.HandlerFunc {
	
}