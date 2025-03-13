package products

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Description string             `json:"description" bson:"description" validate:"required"`
	Price       float64            `json:"price" bson:"price" validate:"required"`
	Image       string             `json:"image" bson:"image"`
}

type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	Total    int `json:"total"`
}

type SearchParams struct {
	Query    string   `json:"query" form:"q"`
	MinPrice *float64 `json:"min_price" form:"min_price"`
	MaxPrice *float64 `json:"max_price" form:"max_price"`
}
