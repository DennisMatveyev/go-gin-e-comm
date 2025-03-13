package products

import (
	"context"
	"go-gin-e-comm/common"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	CreateProduct(c context.Context, product *Product) error
	GetProducts(c context.Context, p Pagination) ([]Product, *Pagination, error)
	GetProductByID(id primitive.ObjectID) (*Product, error)
	SearchProducts(c context.Context, params SearchParams) ([]Product, error)
}

type productRepository struct {
	db  *mongo.Database
	log *slog.Logger
}

func NewProductRepository(db *mongo.Database, log *slog.Logger) ProductRepository {
	return &productRepository{db: db, log: log}
}

func (r *productRepository) GetProducts(c context.Context, p Pagination) ([]Product, *Pagination, error) {
	var products []Product
	skip := (p.Page - 1) * p.PageSize
	total, err := r.db.Collection("products").CountDocuments(c, bson.M{})
	if err != nil {
		r.log.Error("Failed to count products", "error", err.Error())
		return nil, nil, common.ErrDatabase
	}
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(p.PageSize)).
		SetSort(bson.M{"_id": -1})

	cursor, err := r.db.Collection("products").Find(c, bson.M{}, opts)
	if err != nil {
		r.log.Error("Failed to get products", "error", err.Error())
		return nil, nil, common.ErrDatabase
	}
	defer cursor.Close(c)

	if err := cursor.All(c, &products); err != nil {
		r.log.Error("Failed to decode products", "error", err.Error())
		return nil, nil, common.ErrDatabase
	}
	pagination := &Pagination{
		Page:     p.Page,
		PageSize: p.PageSize,
		Total:    int(total),
	}
	return products, pagination, nil
}

func (r *productRepository) GetProductByID(id primitive.ObjectID) (*Product, error) {
	var product Product
	err := r.db.Collection("products").FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	if err != nil {
		r.log.Error("Failed to get product by ID", "error", err.Error())
		return nil, common.ErrDatabase
	}
	return &product, nil
}

func (r *productRepository) CreateProduct(c context.Context, product *Product) error {
	_, err := r.db.Collection("products").InsertOne(c, product)
	if err != nil {
		r.log.Error("Failed to create product", "error", err.Error())
		return common.ErrDatabase
	}
	return nil
}

func (r *productRepository) SearchProducts(c context.Context, params SearchParams) ([]Product, error) {
	var products []Product
	conditions := []bson.M{
		{
			"$or": []bson.M{
				{"name": bson.M{"$regex": params.Query, "$options": "i"}},
				{"description": bson.M{"$regex": params.Query, "$options": "i"}},
			},
		},
	}
	if params.MinPrice != nil || params.MaxPrice != nil {
		priceFilter := bson.M{}
		if params.MinPrice != nil {
			priceFilter["$gte"] = *params.MinPrice
		}
		if params.MaxPrice != nil {
			priceFilter["$lte"] = *params.MaxPrice
		}
		if len(priceFilter) > 0 {
			conditions = append(conditions, bson.M{"price": priceFilter})
		}
	}
	filter := bson.M{"$and": conditions}

	cursor, err := r.db.Collection("products").Find(c, filter)
	if err != nil {
		r.log.Error("Failed to search products", "error", err.Error())
		return nil, common.ErrDatabase
	}
	defer cursor.Close(c)

	if err := cursor.All(c, &products); err != nil {
		r.log.Error("Failed to decode products", "error", err.Error())
		return nil, common.ErrDatabase
	}
	return products, nil
}
