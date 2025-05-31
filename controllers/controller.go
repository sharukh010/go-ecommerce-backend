package controllers

import (
	"context"
	"fmt"
	"log"
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
	generate "github.com/sharukh010/go-ecommerce/tokens"
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

func VerifyPassword(userPassword string, givenPassword string) (valid bool, msg string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword),[]byte(userPassword))
	if err != nil {
		msg = "Login Or Password is Incorrect"
		valid = false
	}
	valid = true
	msg = ""
	return valid,msg
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
	
		token,refreshToken,_ := generate.TokenGenerator(*user.Email,*user.First_Name,*user.Last_Name,user.User_ID)

		user.Token = &token 
		user.Refresh_Token = &refreshToken

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
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel() 
		var user models.User 
		var foundUser models.User 
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return 
		}
		err := UserCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Invalid credentials"})
			return 
		}
		PasswordIsValid,msg := VerifyPassword(*user.Password,*foundUser.Password)
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			fmt.Println(msg)
			return 
		}
		token,refreshToken,_ := generate.TokenGenerator(*foundUser.Email,*foundUser.First_Name,*foundUser.Last_Name,foundUser.User_ID)
		generate.UpdateAllTokens(token,refreshToken,foundUser.User_ID)
		c.JSON(http.StatusFound,foundUser)
	}
}

func ProductViewerAdmin() gin.HandlerFunc{

}

func SearchProduct() gin.HandlerFunc{

}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context){
		var SearchProducts []models.Product
		queryParam := c.Query("name")
		if queryParam == ""{
			log.Println("Query is empty")
			c.Header("Content-Type","Application/json")
			c.JSON(http.StatusNotFound,gin.H{"Error":"Invalid Search Index"})
			c.Abort()
			return 
		}
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		searchQueryDB,err := ProductCollection.Find(ctx,bson.M{"product_name":bson.M{"$regex":queryParam}})

		if err != nil {
			c.IndentedJSON(http.StatusNotFound,"Product not found")
			return 
		}

		err = searchQueryDB.All(ctx,&SearchProducts)
		if err != nil {
			log.Println(err.Error())
			c.IndentedJSON(http.StatusBadRequest,"Invalid Request")
			return 
		}

		defer searchQueryDB.Close(ctx)

		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest,"Invalid Request")
			return 
		}
		c.IndentedJSON(http.StatusOK,SearchProducts)


	}
}