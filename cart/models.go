package cart

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id" validate:"required"`
	Quantity  int                `json:"quantity" bson:"quantity" validate:"required,min=1"`
}

type Cart struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
	Items  []CartItem         `json:"items" bson:"items"`
}

type Order struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	CartID    primitive.ObjectID `json:"cart_id" bson:"cart_id"`
	Total     float64            `json:"total" bson:"total"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
