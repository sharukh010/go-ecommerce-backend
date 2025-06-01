package database

import (
	"context"
	"errors"
	"log"

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

func AddProductToCart() {

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

func InstantBuyer() {

}
