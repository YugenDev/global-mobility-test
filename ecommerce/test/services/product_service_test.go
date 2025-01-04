package services_test

import (
	"testing"

	"github.com/YugenDev/global-mobility-test/internal/models"
	"github.com/YugenDev/global-mobility-test/internal/services"
	"github.com/YugenDev/global-mobility-test/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CreateProduct(c echo.Context, product *models.Product) (*mongo.InsertOneResult, error) {
	args := m.Called(c, product)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockProductRepository) GetProductByID(id string) (models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAllProducts() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateProduct(c echo.Context, id string, product *models.Product) (*mongo.UpdateResult, error) {
	args := m.Called(c, id, product)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockProductRepository) DeleteProduct(c echo.Context, id string) (*mongo.DeleteResult, error) {
	args := m.Called(c, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func TestCreateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		product      *models.Product
		mockBehavior func()
		expectedErr  error
	}{
		{
			name: "Valid Product",
			product: &models.Product{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
				Stock:       10,
			},
			mockBehavior: func() {
				mockRepo.On("CreateProduct", mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Missing Name",
			product: &models.Product{
				Description: "Test Description",
				Price:       100,
				Stock:       10,
			},
			mockBehavior: func() {},
			expectedErr:  utils.ErrProductNameRequired,
		},
		{
			name: "Missing Description",
			product: &models.Product{
				Name:  "Test Product",
				Price: 100,
				Stock: 10,
			},
			mockBehavior: func() {},
			expectedErr:  utils.ErrProductDescriptionRequired,
		},
		{
			name: "Invalid Price",
			product: &models.Product{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       -10,
				Stock:       10,
			},
			mockBehavior: func() {},
			expectedErr:  utils.ErrProductPriceInvalid,
		},
		{
			name: "Invalid Stock",
			product: &models.Product{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
				Stock:       -5,
			},
			mockBehavior: func() {},
			expectedErr:  utils.ErrProductStockInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			c := echo.New().NewContext(nil, nil)
			err := service.CreateProduct(c, tt.product)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetAll(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		mockBehavior func()
		expectedErr  error
		expectedData []models.Product
	}{
		{
			name: "No Products Found",
			mockBehavior: func() {
				mockRepo.On("GetAllProducts").Return([]models.Product{}, utils.ErrNoProductsFound)
			},
			expectedErr:  utils.ErrNoProductsFound,
			expectedData: nil,
		},
		{
			name: "Products Found",
			mockBehavior: func() {
				products := []models.Product{
					{ProductID: "1", Name: "Product 1"},
					{ProductID: "2", Name: "Product 2"},
				}
				mockRepo.On("GetAllProducts").Return(products, nil)
			},
			expectedErr:  nil,
			expectedData: []models.Product{{ProductID: "1", Name: "Product 1"}, {ProductID: "2", Name: "Product 2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{} // Reset mock
			tt.mockBehavior()

			products, err := service.GetAll()
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedData, products)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetByID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		productID    string
		mockBehavior func()
		expectedErr  error
		expectedData models.Product
	}{
		{
			name:      "Product Not Found",
			productID: "nonexistent",
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "nonexistent").Return(models.Product{}, utils.ErrNoProductsFound)
			},
			expectedErr:  utils.ErrNoProductsFound,
			expectedData: models.Product{},
		},
		{
			name:      "Product Found",
			productID: "existing-id",
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "existing-id").Return(
					models.Product{
						ProductID: "existing-id",
						Name:      "Test Product",
					},
					nil,
				)
			},
			expectedErr: nil,
			expectedData: models.Product{
				ProductID: "existing-id",
				Name:      "Test Product",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{} // Reset mock
			tt.mockBehavior()

			product, err := service.GetByID(tt.productID)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedData, product)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		id           string
		product      *models.Product
		mockBehavior func()
		expectedErr  error
	}{
		{
			name: "Valid Update",
			id:   "123",
			product: &models.Product{
				Name:        "Updated Product",
				Description: "Updated Description",
				Price:       200,
				Stock:       20,
			},
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "123").Return(models.Product{Name: "Product 1"}, nil)
				mockRepo.On("UpdateProduct", mock.Anything, "123", mock.Anything).Return(&mongo.UpdateResult{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Missing ID",
			id:   "",
			product: &models.Product{
				Name:        "Updated Product",
				Description: "Updated Description",
				Price:       200,
				Stock:       20,
			},
			mockBehavior: func() {},
			expectedErr:  utils.ErrProductIDRequired,
		},
		{
			name: "Invalid Price",
			id:   "123",
			product: &models.Product{
				Price: -10,
			},
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "123").Return(models.Product{Name: "Product 1"}, nil)
			},
			expectedErr: utils.ErrProductPriceInvalid,
		},
		{
			name: "Invalid Stock",
			id:   "123",
			product: &models.Product{
				Stock: -5,
			},
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "123").Return(models.Product{Name: "Product 1"}, nil)
			},
			expectedErr: utils.ErrProductStockInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()
			c := echo.New().NewContext(nil, nil)
			err := service.UpdateProduct(c, tt.id, tt.product)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		productID    string
		mockBehavior func()
		expectedErr  error
	}{
		{
			name:      "Valid Delete",
			productID: "valid-id",
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "valid-id").Return(models.Product{ProductID: "valid-id"}, nil)
				mockRepo.On("DeleteProduct", mock.Anything, "valid-id").Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
			},
			expectedErr: nil,
		},
		{
			name:      "Product Not Found",
			productID: "nonexistent",
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "nonexistent").Return(models.Product{}, utils.ErrNoProductsFound)
			},
			expectedErr: utils.ErrNoProductsFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{} // Reset mock
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.DeleteProduct(c, tt.productID)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
