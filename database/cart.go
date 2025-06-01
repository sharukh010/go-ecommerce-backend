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

func BuyItemFromCart() {

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
	filter2 := bson.D{primitive.E{Key:"_id",Value:id}}
	update2 := bson.M{"$push":bson.M{"orders.$[].order_list":productDetails}}
	_,err = userCollection.UpdateOne(ctx,filter2,update2)
	if err != nil {
		log.Println(err.Error())
		return ErrCantBuyCartItem
	}
	return nil 
}
