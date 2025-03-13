package cart

import (
	"context"
	"go-gin-e-comm/common"
	"go-gin-e-comm/products"
	"log/slog"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CartRepository interface {
	GetCart(c context.Context, userID primitive.ObjectID) (*Cart, error)
	AddProductToCart(c context.Context, userID, productID primitive.ObjectID, quantity int) error
	RemoveProductFromCart(c context.Context, userID, productID primitive.ObjectID) error
	Checkout(c context.Context, userID primitive.ObjectID) (*Order, error)
}

type cartRepository struct {
	db  *mongo.Database
	log *slog.Logger
}

func NewCartRepository(db *mongo.Database, log *slog.Logger) CartRepository {
	return &cartRepository{db: db, log: log}
}

func (r *cartRepository) GetCart(c context.Context, userID primitive.ObjectID) (*Cart, error) {
	var cart Cart
	err := r.db.Collection("carts").FindOne(c, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCartNotFound
		}
		r.log.Error("Failed to get cart", "error", err.Error())
		return nil, common.ErrDatabase
	}
	return &cart, nil
}

func (r *cartRepository) AddProductToCart(c context.Context, userID, productID primitive.ObjectID, quantity int) error {
	cart, err := r.GetCart(c, userID)
	if err != nil {
		if err == ErrCartNotFound {
			cart = &Cart{UserID: userID, Items: []CartItem{}}
		} else {
			return err
		}
	}
	productExists := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Quantity += quantity
			productExists = true
			break
		}
	}
	if !productExists {
		cart.Items = append(cart.Items, CartItem{
			ProductID: productID,
			Quantity:  quantity,
		})
	}
	_, err = r.db.Collection("carts").UpdateOne(c,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"items": cart.Items}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		r.log.Error("Failed to add product to cart", "error", err.Error())
		return common.ErrDatabase
	}
	return nil
}

func (r *cartRepository) RemoveProductFromCart(c context.Context, userID, productID primitive.ObjectID) error {
	cart, err := r.GetCart(c, userID)
	if err != nil {
		return err
	}
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = slices.Delete(cart.Items, i, i+1)
			break
		}
	}
	if len(cart.Items) == 0 {
		_, err = r.db.Collection("carts").DeleteOne(c, bson.M{"user_id": userID})
		if err != nil {
			r.log.Error("Failed to remove cart", "error", err.Error())
			return common.ErrDatabase
		}
		return nil
	}
	_, err = r.db.Collection("carts").UpdateOne(c,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"items": cart.Items}},
	)
	if err != nil {
		r.log.Error("Failed to remove product from cart", "error", err.Error())
		return common.ErrDatabase
	}
	return nil
}

func (r *cartRepository) Checkout(c context.Context, userID primitive.ObjectID) (*Order, error) {
	cart, err := r.GetCart(c, userID)
	if err != nil {
		return nil, err
	}
	total := 0.0
	for _, item := range cart.Items {
		var product products.Product
		err := r.db.Collection("products").FindOne(c, bson.M{"_id": item.ProductID}).Decode(&product)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				r.log.Error("Product not found", "product_id", item.ProductID.Hex())
				return nil, ErrProductInCartNotFound
			}
			r.log.Error("Failed to get product", "error", err.Error())
			return nil, common.ErrDatabase
		}
		total += float64(item.Quantity) * product.Price
	}
	order := &Order{
		UserID:    userID,
		CartID:    cart.ID,
		Total:     total,
		CreatedAt: time.Now(),
	}
	_, err = r.db.Collection("orders").InsertOne(c, order)
	if err != nil {
		r.log.Error("Failed to create order", "error", err.Error())
		return nil, common.ErrDatabase
	}
	_, err = r.db.Collection("carts").DeleteOne(c, bson.M{"user_id": userID})
	if err != nil {
		r.log.Error("Failed to remove cart after checkout", "error", err.Error())
		return nil, common.ErrDatabase
	}
	return order, nil
}
