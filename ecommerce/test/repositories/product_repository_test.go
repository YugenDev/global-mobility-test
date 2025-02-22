package repositories_test

import (
	"context"
	"errors"
	"testing"

	"github.com/YugenDev/global-mobility-test/internal/models"
	"github.com/YugenDev/global-mobility-test/internal/repositories"
	"github.com/YugenDev/global-mobility-test/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func TestDeleteProduct(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	mockResult := &mongo.DeleteResult{DeletedCount: 1}
	mockCollection.On("DeleteOne", mock.Anything, mock.Anything).Return(mockResult, nil)

	result, err := repo.DeleteProduct(echo.New().NewContext(nil, nil), "test-id")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.DeletedCount)
	mockCollection.AssertExpectations(t)
}

func TestDeleteProduct_EmptyID(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	result, err := repo.DeleteProduct(echo.New().NewContext(nil, nil), "")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "product ID is required", err.Error())
}

func TestGetAllProducts(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	products := []models.Product{
		{ProductID: "1", Name: "Product 1"},
		{ProductID: "2", Name: "Product 2"},
	}

	// Convert products to []interface{}
	docs := make([]interface{}, len(products))
	for i, v := range products {
		docs[i] = v
	}

	// Create a mock cursor
	cursor, err := mongo.NewCursorFromDocuments(docs, nil, nil)
	assert.NoError(t, err)

	mockCollection.On("Find", mock.Anything, mock.Anything).Return(cursor, nil)

	result, err := repo.GetAllProducts()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, products[0].ProductID, result[0].ProductID)
	assert.Equal(t, products[1].ProductID, result[1].ProductID)
	mockCollection.AssertExpectations(t)
}
func TestGetAllProducts_DatabaseError(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	expectedError := errors.New("database error")
	mockCollection.On("Find", mock.Anything, mock.Anything).Return(nil, expectedError)

	result, err := repo.GetAllProducts()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockCollection.AssertExpectations(t)
}

func TestDeleteProduct_DatabaseError(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	expectedError := errors.New("database error")
	mockCollection.On("DeleteOne", mock.Anything, mock.Anything).Return(nil, expectedError)

	result, err := repo.DeleteProduct(echo.New().NewContext(nil, nil), "test-id")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockCollection.AssertExpectations(t)
}

func TestGetAllProducts_CursorError(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	// Create cursor with error
	cursor, err := mongo.NewCursorFromDocuments(nil, errors.New("cursor error"), nil)
	assert.NoError(t, err)

	mockCollection.On("Find", mock.Anything, mock.Anything).Return(cursor, nil)

	result, err := repo.GetAllProducts()

	assert.Error(t, err)
	assert.Nil(t, result)
	mockCollection.AssertExpectations(t)
}
func TestCreateProduct(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	product := &models.Product{
		ProductID: "test-id",
		Name:      "Test Product",
	}

	mockResult := &mongo.InsertOneResult{InsertedID: "test-id"}
	mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(mockResult, nil)

	result, err := repo.CreateProduct(echo.New().NewContext(nil, nil), product)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockResult, result)
	mockCollection.AssertExpectations(t)
}

func TestCreateProduct_NilProduct(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	result, err := repo.CreateProduct(echo.New().NewContext(nil, nil), nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, utils.ErrNullProductData, err)
}

func TestCreateProduct_NilCollection(t *testing.T) {
	repo := repositories.ProductRepository{Collection: nil}
	product := &models.Product{ProductID: "test-id"}

	result, err := repo.CreateProduct(echo.New().NewContext(nil, nil), product)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, utils.ErrDatabaseNotInitialized, err)
}

func TestCreateProduct_DatabaseError(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	product := &models.Product{ProductID: "test-id"}
	expectedError := errors.New("database error")
	mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(nil, expectedError)

	result, err := repo.CreateProduct(echo.New().NewContext(nil, nil), product)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockCollection.AssertExpectations(t)
}
func TestGetProductByID(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	product := models.Product{
		ProductID: "test-id",
		Name:      "Test Product",
	}

	mockSingleResult := mongo.NewSingleResultFromDocument(product, nil, nil)
	mockCollection.On("FindOne", mock.Anything, bson.M{"product_id": "test-id"}).Return(mockSingleResult)

	result, err := repo.GetProductByID("test-id")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, product.ProductID, result.ProductID)
	assert.Equal(t, product.Name, result.Name)
	mockCollection.AssertExpectations(t)
}

func TestGetProductByID_EmptyID(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	result, err := repo.GetProductByID("")

	assert.Error(t, err)
	assert.Equal(t, utils.ErrProductIDRequired, err)
	assert.Equal(t, models.Product{}, result)
}

func TestGetProductByID_DatabaseError(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	expectedError := errors.New("database error")
	mockSingleResult := mongo.NewSingleResultFromDocument(bson.M{}, expectedError, nil)
	mockCollection.On("FindOne", mock.Anything, bson.M{"product_id": "test-id"}).Return(mockSingleResult)

	result, err := repo.GetProductByID("test-id")

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, models.Product{}, result)
	mockCollection.AssertExpectations(t)
}
func TestUpdateProduct(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	product := &models.Product{
		ProductID: "test-id",
		Name:      "Updated Product",
	}

	mockResult := &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"product_id": "test-id"}, mock.Anything).Return(mockResult, nil)

	result, err := repo.UpdateProduct(echo.New().NewContext(nil, nil), "test-id", product)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.MatchedCount)
	assert.Equal(t, int64(1), result.ModifiedCount)
	mockCollection.AssertExpectations(t)
}

func TestUpdateProduct_EmptyID(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	product := &models.Product{
		ProductID: "test-id",
		Name:      "Updated Product",
	}

	result, err := repo.UpdateProduct(echo.New().NewContext(nil, nil), "", product)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, utils.ErrProductIDRequired, err)
}

func TestUpdateProduct_NilProduct(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	result, err := repo.UpdateProduct(echo.New().NewContext(nil, nil), "test-id", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, utils.ErrNullProductData, err)
}

func TestUpdateProduct_DatabaseError(t *testing.T) {
	mockCollection := new(MockCollection)
	repo := repositories.ProductRepository{Collection: mockCollection}

	product := &models.Product{
		ProductID: "test-id",
		Name:      "Updated Product",
	}

	expectedError := errors.New("database error")
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"product_id": "test-id"}, mock.Anything).Return(nil, expectedError)

	result, err := repo.UpdateProduct(echo.New().NewContext(nil, nil), "test-id", product)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockCollection.AssertExpectations(t)
}
