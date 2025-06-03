package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/sharukh010/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct = errors.New("can't find product")
	ErrCantDecodeProducts = errors.New("can't decode product")
	ErrUserIDIsNotValid = errors.New("user is not valid")
	ErrCantUpdateUser = errors.New("cannot add product to cart")
	ErrCantRemoveItem = errors.New("cannot remove item from cart")
	ErrCantGetItem = errors.New("cannot get item from cart")
	ErrCantBuyCartItem = errors.New("cannot update the purchase")
	ErrCartIsEmpty = errors.New("cannot buy cart is empty")
)

func AddProductToCart(ctx context.Context,prodCollection,userCollection *mongo.Collection,productID primitive.ObjectID,userID string) error {
	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err.Error())
		return ErrUserIDIsNotValid
	}
	
	var product models.ProductUser
	err = prodCollection.FindOne(ctx,bson.M{"_id":productID}).Decode(&product)
	if err != nil {
		log.Println(err.Error())
		return ErrCantFindProduct
	}
	filter := bson.M{"_id":id}
	update := bson.M{"$push":bson.M{"usercart":product}}

	_,err = userCollection.UpdateOne(ctx,filter,update)
	if err != nil {
		log.Println(err.Error())
		return ErrCantUpdateUser
	}
	return nil 
}

func RemoveCartItem(ctx context.Context,userCollection *mongo.Collection,productID primitive.ObjectID,userID string)error {
	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err.Error())
		return ErrUserIDIsNotValid
	}
	filter := bson.D{primitive.E{Key:"_id",Value:id}}
	update := bson.M{"$pull":bson.M{"usercart":bson.M{"_id":productID}}}
	_,err = userCollection.UpdateMany(ctx,filter,update)
	if err != nil {
		return ErrCantRemoveItem
	}
	return nil 

}

func BuyItemFromCart(ctx context.Context,userCollection *mongo.Collection,userID string) error {
	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err.Error())
		return ErrUserIDIsNotValid
	}
	var getUser models.User 
	var orderCart models.Order 
	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Orderered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser,0)
	orderCart.Payment_Method.COD = true 
	unwind := bson.D{{Key:"$unwind",Value:bson.D{primitive.E{Key:"path",Value:"$usercart"}}}}
	grouping := bson.D{{Key:"$group",Value:bson.D{primitive.E{Key:"_id",Value:"$_id"},{Key:"total",Value:bson.D{{Key:"$sum",Value:"$usercart.price"}}}}}}
	currentResults,err := userCollection.Aggregate(ctx,mongo.Pipeline{unwind,grouping})
	if err != nil {
		log.Println(err.Error())
		return ErrCantBuyCartItem
	}
	var getUserCart []bson.M 
	if err = currentResults.All(ctx,&getUserCart);err != nil{
		log.Println(err.Error())
		return ErrCantBuyCartItem
	}
	var totalPrice int32
	//check this //
	for _,userItem := range getUserCart {
		price := userItem["total"]
		totalPrice = price.(int32)
	}
	err = userCollection.FindOne(ctx,bson.D{primitive.E{Key:"_id",Value:id}}).Decode(&getUser)
	if err != nil {
		log.Println(err.Error())
		return ErrUserIDIsNotValid
	}

	if len(getUser.UserCart) == 0 {
		return ErrCartIsEmpty
	}

	
	orderCart.Order_Cart = getUser.UserCart
	orderCart.Price = int(totalPrice)

	filter := bson.D{primitive.E{Key:"_id",Value:id}}
	update := bson.D{{Key:"$push",Value:bson.D{primitive.E{Key:"orders",Value:orderCart}}}}
	_,err = userCollection.UpdateMany(ctx,filter,update)
	if err != nil {
		log.Println(err.Error())
		return ErrCantUpdateUser
	}
	

	userCartEmpty := make([]models.ProductUser,0)
	filtered := bson.D{primitive.E{Key:"_id",Value:id}}
	updated := bson.D{{Key:"$set",Value:bson.D{primitive.E{Key:"usercart",Value:userCartEmpty}}}}
	_,err = userCollection.UpdateOne(ctx,filtered,updated)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil 

}

func InstantBuyer(ctx context.Context,prodCollection,userCollection *mongo.Collection,productID primitive.ObjectID,userID string) error{
	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err.Error())
		return ErrUserIDIsNotValid
	}
	var productDetails models.ProductUser
	var orderDetails models.Order 
	orderDetails.Order_ID = primitive.NewObjectID()
	orderDetails.Orderered_At = time.Now()
	orderDetails.Order_Cart = make([]models.ProductUser,0)
	orderDetails.Payment_Method.COD = true 
	err = prodCollection.FindOne(ctx,bson.D{primitive.E{Key:"_id",Value:productID}}).Decode(&productDetails)
	if err != nil {
		log.Println(err.Error())
		return ErrCantFindProduct
	}
	orderDetails.Price = productDetails.Price
	filter := bson.D{primitive.E{Key:"_id",Value:id}}
	update := bson.D{{Key:"$push",Value:bson.D{primitive.E{Key:"orders",Value:orderDetails}}}}
	_,err = userCollection.UpdateOne(ctx,filter,update)
	if err != nil {
		log.Println(err.Error())
		return ErrCantBuyCartItem
	}
	return nil 
}
