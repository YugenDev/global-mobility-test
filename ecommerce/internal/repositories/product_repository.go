package repositories

import (
	"context"
	"log"
	"time"

	"github.com/YugenDev/global-mobility-test/internal/config"
	"github.com/YugenDev/global-mobility-test/internal/models"
	"github.com/YugenDev/global-mobility-test/internal/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IProductRepository interface {
	CreateProduct(c echo.Context, product *models.Product) (*mongo.InsertOneResult, error)
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id string) (models.Product, error)
	UpdateProduct(c echo.Context, id string, product *models.Product) (*mongo.UpdateResult, error)
	DeleteProduct(c echo.Context, id string) (*mongo.DeleteResult, error)
}

type MongoCollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

type ProductRepository struct {
	Collection MongoCollection
}

var _ IProductRepository = (*ProductRepository)(nil)

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		Collection: config.GetCollection("products"),
	}
}

func (r *ProductRepository) CreateProduct(c echo.Context, product *models.Product) (*mongo.InsertOneResult, error) {
	if r.Collection == nil {
		return nil, utils.ErrDatabaseNotInitialized
	}
	if product == nil {
		return nil, utils.ErrNullProductData
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	result, err := r.Collection.InsertOne(ctx, product)
	if err != nil {
		log.Println("Error creating product: ", err)
		return nil, err
	}

	return result, nil
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var products []models.Product
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Error getting products: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			log.Println("Error decoding product: ", err)
			continue
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor error: ", err)
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetProductByID(id string) (models.Product, error) {
	if id == "" {
		return models.Product{}, utils.ErrProductIDRequired
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product models.Product
	err := r.Collection.FindOne(ctx, bson.M{"product_id": id}).Decode(&product)
	if err != nil {
		log.Println("Error getting product by ID: ", err)
		return models.Product{}, err
	}

	return product, nil
}

func (r *ProductRepository) UpdateProduct(c echo.Context, id string, product *models.Product) (*mongo.UpdateResult, error) {
	if id == "" {
		return nil, utils.ErrProductIDRequired
	}
	if product == nil {
		return nil, utils.ErrNullProductData
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product.UpdatedAt = time.Now()

	result, err := r.Collection.UpdateOne(ctx, bson.M{"product_id": id}, bson.M{"$set": product})
	if err != nil {
		log.Println("Error updating product: ", err)
		return nil, err
	}

	return result, nil
}

func (r *ProductRepository) DeleteProduct(c echo.Context, id string) (*mongo.DeleteResult, error) {
	if id == "" {
		return nil, utils.ErrProductIDRequired
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.Collection.DeleteOne(ctx, bson.M{"product_id": id})
	if err != nil {
		log.Println("Error deleting product: ", err)
		return nil, err
	}

	return result, nil
}
