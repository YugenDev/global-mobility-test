package services_test

import (
	"errors"
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
func TestCreateProductWithExistingID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		product      *models.Product
		mockBehavior func()
		expectedErr  error
	}{
		{
			name: "Existing Product ID",
			product: &models.Product{
				ProductID:   "existing-id",
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
				Stock:       10,
			},
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "existing-id").Return(
					models.Product{ProductID: "existing-id"},
					nil,
				)
			},
			expectedErr: utils.ErrProductIDAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{} // Reset mock
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.CreateProduct(c, tt.product)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestUpdateProductWithInvalidID(t *testing.T) {
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
			name: "Invalid Product ID Change",
			id:   "123",
			product: &models.Product{
				ProductID: "456", // Different from original ID
				Name:      "Test Product",
			},
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "123").Return(models.Product{
					ProductID: "123",
					Name:      "Original Product",
				}, nil)
			},
			expectedErr: utils.ErrProductIDCannotBeChanged,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{} // Reset mock
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.UpdateProduct(c, tt.id, tt.product)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestCreateProductEmptyID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		product      *models.Product
		mockBehavior func()
		expectedErr  error
	}{
		{
			name: "Empty ID Gets Generated",
			product: &models.Product{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
				Stock:       10,
			},
			mockBehavior: func() {
				mockRepo.On("CreateProduct", mock.Anything, mock.MatchedBy(func(p *models.Product) bool {
					return p.ProductID != ""
				})).Return(&mongo.InsertOneResult{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{}
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.CreateProduct(c, tt.product)

			assert.Equal(t, tt.expectedErr, err)
			assert.NotEmpty(t, tt.product.ProductID)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestGetByIDEmptyID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		productID    string
		expectedErr  error
		expectedData models.Product
	}{
		{
			name:         "Empty ID",
			productID:    "",
			expectedErr:  utils.ErrProductIDRequired,
			expectedData: models.Product{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := service.GetByID(tt.productID)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedData, product)
		})
	}
}
func TestDeleteProductEmptyProductID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		productID    string
		mockBehavior func()
		expectedErr  error
	}{
		{
			name:      "Empty ProductID from GetProductByID",
			productID: "test-id",
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "test-id").Return(
					models.Product{ProductID: ""},
					nil,
				)
			},
			expectedErr: utils.ErrNoProductsFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{}
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.DeleteProduct(c, tt.productID)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestEmptyID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		productID    string
		mockBehavior func()
		expectedErr  error
	}{
		{
			name:         "Empty ID on Delete",
			productID:    "",
			mockBehavior: func() {},
			expectedErr:  utils.ErrProductIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{}
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.DeleteProduct(c, tt.productID)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestGetByIDWithError(t *testing.T) {
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
			name:      "Repository Error",
			productID: "test-id",
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "test-id").Return(
					models.Product{},
					errors.New("repository error"),
				)
			},
			expectedErr:  errors.New("repository error"),
			expectedData: models.Product{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{}
			tt.mockBehavior()

			product, err := service.GetByID(tt.productID)

			assert.Equal(t, tt.expectedErr.Error(), err.Error())
			assert.Equal(t, tt.expectedData, product)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestCreateProductWithRepositoryError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		product      *models.Product
		mockBehavior func()
		expectedErr  error
	}{
		{
			name: "Repository Error",
			product: &models.Product{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100,
				Stock:       10,
			},
			mockBehavior: func() {
				mockRepo.On("CreateProduct", mock.Anything, mock.Anything).Return(
					&mongo.InsertOneResult{},
					errors.New("repository error"),
				)
			},
			expectedErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{}
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.CreateProduct(c, tt.product)

			assert.Equal(t, tt.expectedErr.Error(), err.Error())
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestUpdateProductWithRepositoryError(t *testing.T) {
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
			name: "Repository Error on Update",
			id:   "test-id",
			product: &models.Product{
				Name:  "Updated Product",
				Price: 100,
			},
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "test-id").Return(models.Product{
					ProductID: "test-id",
					Name:      "Original Product",
				}, nil)
				mockRepo.On("UpdateProduct", mock.Anything, "test-id", mock.Anything).Return(
					&mongo.UpdateResult{},
					errors.New("repository error"),
				)
			},
			expectedErr: errors.New("repository error"),
		},
		{
			name: "Repository Error on GetProductByID",
			id:   "test-id",
			product: &models.Product{
				Name: "Updated Product",
			},
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "test-id").Return(
					models.Product{},
					errors.New("repository error"),
				)
			},
			expectedErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{}
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.UpdateProduct(c, tt.id, tt.product)

			assert.Equal(t, tt.expectedErr.Error(), err.Error())
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestDeleteProductWithRepositoryError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		productID    string
		mockBehavior func()
		expectedErr  error
	}{
		{
			name:      "Repository Error on Delete",
			productID: "test-id",
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "test-id").Return(models.Product{
					ProductID: "test-id",
				}, nil)
				mockRepo.On("DeleteProduct", mock.Anything, "test-id").Return(
					nil,
					errors.New("repository error"),
				)
			},
			expectedErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{}
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.DeleteProduct(c, tt.productID)

			assert.Equal(t, tt.expectedErr.Error(), err.Error())
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestGetAllWithRepositoryError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		mockBehavior func()
		expectedErr  error
		expectedData []models.Product
	}{
		{
			name: "Repository Error",
			mockBehavior: func() {
				mockRepo.On("GetAllProducts").Return(
					[]models.Product{},
					errors.New("repository error"),
				)
			},
			expectedErr:  errors.New("repository error"),
			expectedData: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{} // Reset mock
			tt.mockBehavior()

			products, err := service.GetAll()
			assert.Equal(t, tt.expectedErr.Error(), err.Error())
			assert.Equal(t, tt.expectedData, products)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestUpdateProductPartialFields(t *testing.T) {
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
			name: "Update Only Name",
			id:   "test-id",
			product: &models.Product{
				Name: "Updated Name",
			},
			mockBehavior: func() {
				existingProduct := models.Product{
					ProductID:   "test-id",
					Name:        "Original Name",
					Description: "Original Description",
					Price:       100,
					Stock:       10,
				}
				mockRepo.On("GetProductByID", "test-id").Return(existingProduct, nil)
				mockRepo.On("UpdateProduct", mock.Anything, "test-id", mock.MatchedBy(func(p *models.Product) bool {
					return p.Name == "Updated Name" &&
						p.Description == "Original Description" &&
						p.Price == 100 &&
						p.Stock == 10
				})).Return(&mongo.UpdateResult{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Update Only Description",
			id:   "test-id",
			product: &models.Product{
				Description: "Updated Description",
			},
			mockBehavior: func() {
				existingProduct := models.Product{
					ProductID:   "test-id",
					Name:        "Original Name",
					Description: "Original Description",
					Price:       100,
					Stock:       10,
				}
				mockRepo.On("GetProductByID", "test-id").Return(existingProduct, nil)
				mockRepo.On("UpdateProduct", mock.Anything, "test-id", mock.MatchedBy(func(p *models.Product) bool {
					return p.Name == "Original Name" &&
						p.Description == "Updated Description" &&
						p.Price == 100 &&
						p.Stock == 10
				})).Return(&mongo.UpdateResult{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{} // Reset mock
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.UpdateProduct(c, tt.id, tt.product)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestUpdateProductMultipleFields(t *testing.T) {
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
			name: "Update Multiple Fields",
			id:   "test-id",
			product: &models.Product{
				Name:        "Updated Name",
				Description: "Updated Description",
				Price:       200,
				Stock:       20,
			},
			mockBehavior: func() {
				existingProduct := models.Product{
					ProductID:   "test-id",
					Name:        "Original Name",
					Description: "Original Description",
					Price:       100,
					Stock:       10,
				}
				mockRepo.On("GetProductByID", "test-id").Return(existingProduct, nil)
				mockRepo.On("UpdateProduct", mock.Anything, "test-id", mock.MatchedBy(func(p *models.Product) bool {
					return p.Name == "Updated Name" &&
						p.Description == "Updated Description" &&
						p.Price == 200 &&
						p.Stock == 20
				})).Return(&mongo.UpdateResult{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{}
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.UpdateProduct(c, tt.id, tt.product)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestUpdateProductPriceAndStockValidations(t *testing.T) {
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
			name: "Update Price and Stock Within Valid Range",
			id:   "test-id",
			product: &models.Product{
				Price: 150.50,
				Stock: 25,
			},
			mockBehavior: func() {
				existingProduct := models.Product{
					ProductID:   "test-id",
					Name:        "Original Name",
					Description: "Original Description",
					Price:       100,
					Stock:       10,
				}
				mockRepo.On("GetProductByID", "test-id").Return(existingProduct, nil)
				mockRepo.On("UpdateProduct", mock.Anything, "test-id", mock.MatchedBy(func(p *models.Product) bool {
					return p.Name == "Original Name" &&
						p.Description == "Original Description" &&
						p.Price == 150.50 &&
						p.Stock == 25
				})).Return(&mongo.UpdateResult{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Update With Invalid Price",
			id:   "test-id",
			product: &models.Product{
				Price: -50,
				Stock: 25,
			},
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "test-id").Return(models.Product{
					ProductID: "test-id",
					Price:     100,
					Stock:     10,
				}, nil)
			},
			expectedErr: utils.ErrProductPriceInvalid,
		},
		{
			name: "Update With Invalid Stock",
			id:   "test-id",
			product: &models.Product{
				Price: 150.50,
				Stock: -10,
			},
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "test-id").Return(models.Product{
					ProductID: "test-id",
					Price:     100,
					Stock:     10,
				}, nil)
			},
			expectedErr: utils.ErrProductStockInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{}
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.UpdateProduct(c, tt.id, tt.product)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestUpdateProductFieldValidations(t *testing.T) {
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
			name:    "Empty Field Updates Maintain Original Values",
			id:      "test-id",
			product: &models.Product{
				// Empty fields should not modify original values
			},
			mockBehavior: func() {
				existingProduct := models.Product{
					ProductID:   "test-id",
					Name:        "Original Name",
					Description: "Original Description",
					Price:       100,
					Stock:       10,
				}
				mockRepo.On("GetProductByID", "test-id").Return(existingProduct, nil)
				mockRepo.On("UpdateProduct", mock.Anything, "test-id", mock.MatchedBy(func(p *models.Product) bool {
					return p.Name == "Original Name" &&
						p.Description == "Original Description" &&
						p.Price == 100 &&
						p.Stock == 10
				})).Return(&mongo.UpdateResult{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Zero Values Don't Override Original Non-Zero Values",
			id:   "test-id",
			product: &models.Product{
				Name:  "New Name",
				Price: 0, // Zero price should not update
				Stock: 0, // Zero stock should not update
			},
			mockBehavior: func() {
				existingProduct := models.Product{
					ProductID:   "test-id",
					Name:        "Original Name",
					Description: "Original Description",
					Price:       100,
					Stock:       10,
				}
				mockRepo.On("GetProductByID", "test-id").Return(existingProduct, nil)
				mockRepo.On("UpdateProduct", mock.Anything, "test-id", mock.MatchedBy(func(p *models.Product) bool {
					return p.Name == "New Name" &&
						p.Description == "Original Description" &&
						p.Price == 100 &&
						p.Stock == 10
				})).Return(&mongo.UpdateResult{}, nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{}
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := service.UpdateProduct(c, tt.id, tt.product)

			assert.Equal(t, tt.expectedErr, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestGeneralRepositoryErrorPropagation(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := services.NewProductService(mockRepo)

	tests := []struct {
		name         string
		operation    string
		setup        func(c echo.Context) error
		mockBehavior func()
		expectedErr  error
	}{
		{
			name:      "GetAll Repository Error Propagation",
			operation: "GetAll",
			setup: func(c echo.Context) error {
				_, err := service.GetAll()
				return err
			},
			mockBehavior: func() {
				mockRepo.On("GetAllProducts").Return(
					[]models.Product{},
					errors.New("database connection error"),
				)
			},
			expectedErr: errors.New("database connection error"),
		},
		{
			name:      "Update Repository Error on Get",
			operation: "Update",
			setup: func(c echo.Context) error {
				return service.UpdateProduct(c, "test-id", &models.Product{
					Name: "Test Product",
				})
			},
			mockBehavior: func() {
				mockRepo.On("GetProductByID", "test-id").Return(
					models.Product{},
					errors.New("database connection error"),
				)
			},
			expectedErr: errors.New("database connection error"),
		},
		{
			name:      "Create Repository Error",
			operation: "Create",
			setup: func(c echo.Context) error {
				return service.CreateProduct(c, &models.Product{
					Name:        "Test Product",
					Description: "Test Description",
					Price:       100,
					Stock:       10,
				})
			},
			mockBehavior: func() {
				mockRepo.On("CreateProduct", mock.Anything, mock.Anything).Return(
					&mongo.InsertOneResult{},
					errors.New("database connection error"),
				)
			},
			expectedErr: errors.New("database connection error"),
		},
		{
			name:      "Delete Repository Error",
			operation: "Delete",
			setup: func(c echo.Context) error {
				mockRepo.On("GetProductByID", "test-id").Return(
					models.Product{ProductID: "test-id"},
					nil,
				)
				mockRepo.On("DeleteProduct", mock.Anything, "test-id").Return(
					nil,
					errors.New("database connection error"),
				)
				return service.DeleteProduct(c, "test-id")
			},
			mockBehavior: func() {},
			expectedErr:  errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock = mock.Mock{} // Reset mock
			tt.mockBehavior()

			c := echo.New().NewContext(nil, nil)
			err := tt.setup(c)

			assert.Equal(t, tt.expectedErr.Error(), err.Error())
			mockRepo.AssertExpectations(t)
		})
	}
}
