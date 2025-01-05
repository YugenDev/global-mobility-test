package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/YugenDev/global-mobility-test/internal/handlers"
	"github.com/YugenDev/global-mobility-test/internal/models"
	"github.com/YugenDev/global-mobility-test/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetAll() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductService) GetByID(id string) (models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return models.Product{}, args.Error(1)
	}
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductService) CreateProduct(c echo.Context, product *models.Product) error {
	args := m.Called(c, product)
	return args.Error(0)
}

func (m *MockProductService) UpdateProduct(c echo.Context, id string, product *models.Product) error {
	args := m.Called(c, id, product)
	return args.Error(0)
}

func (m *MockProductService) DeleteProduct(c echo.Context, id string) error {
	args := m.Called(c, id)
	return args.Error(0)
}

func TestGetAllProducts(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	products := []models.Product{{ProductID: "1", Name: "Test Product"}}
	mockService.On("GetAll").Return(products, nil)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.GetAllProducts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetProductByIDInvalidID(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/products/:id")
	c.SetParamNames("id")
	c.SetParamValues("")

	if err := handler.GetProductByID(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreateProductInvalidPayload(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	invalidJSON := `{"name": "Test`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(invalidJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.CreateProduct(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreateProductInvalidData(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := &models.Product{
		Name:        "",
		Description: "",
		Price:       -1,
		Stock:       -1,
	}

	productJSON, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(productJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.CreateProduct(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestGetProductByID(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := models.Product{
		ProductID:   "1",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.0,
		Stock:       5,
	}

	mockService.On("GetByID", "1").Return(product, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/products/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handler.GetProductByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var responseProduct models.Product
		err := json.Unmarshal(rec.Body.Bytes(), &responseProduct)
		assert.NoError(t, err)
		assert.Equal(t, product, responseProduct)
	}
}

func TestGetAllProductsError(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	mockService.On("GetAll").Return([]models.Product{}, utils.ErrNoProductsFound)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.GetAllProducts(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}
func TestCreateProductSuccess(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.0,
		Stock:       5,
	}

	mockService.On("GetByID", mock.Anything).Return(models.Product{}, utils.ErrNoProductsFound)
	mockService.On("CreateProduct", mock.Anything, product).Return(nil)

	productJSON, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(productJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.CreateProduct(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
func TestCreateProductIDExists(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := &models.Product{
		ProductID:   "1",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.0,
		Stock:       5,
	}

	existingProduct := models.Product{
		ProductID: "1",
		Name:      "Existing Product",
	}

	mockService.On("GetByID", "1").Return(existingProduct, nil)

	productJSON, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(productJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.CreateProduct(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, utils.ErrProductIDAlreadyExists.Error(), response["message"])
	}
}
func TestGetProductByIDNotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	mockService.On("GetByID", "999").Return(models.Product{}, utils.ErrNoProductsFound)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/products/:id")
	c.SetParamNames("id")
	c.SetParamValues("999")

	if err := handler.GetProductByID(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, utils.ErrNoProductsFound.Error(), response["message"])
	}
}

func TestGetAllProductsInternalError(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	mockService.On("GetAll").Return([]models.Product{}, utils.ErrInternalServer)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.GetAllProducts(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, utils.ErrInternalServer.Error(), response["message"])
	}
}
func TestCreateProductWithErrors(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	mockService.On("CreateProduct", mock.Anything, mock.Anything).Return(utils.ErrInternalServer)

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.0,
		Stock:       5,
	}

	productJSON, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(productJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.CreateProduct(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, utils.ErrInternalServer.Error(), response["message"])
	}
}

func TestCreateProductEmptyNameField(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := &models.Product{
		Name:        "",
		Description: "Test Description",
		Price:       10.0,
		Stock:       5,
	}

	productJSON, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(productJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.CreateProduct(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, utils.ErrProductNameRequired.Error(), response["message"])
	}
}
func TestCreateProductEmptyDescriptionField(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := &models.Product{
		Name:        "Test Product",
		Description: "",
		Price:       10.0,
		Stock:       5,
	}

	productJSON, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(productJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.CreateProduct(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, utils.ErrProductDescriptionRequired.Error(), response["message"])
	}
}
func TestCreateProductInvalidPrice(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       0,
		Stock:       5,
	}

	productJSON, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(productJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.CreateProduct(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, utils.ErrProductPriceInvalid.Error(), response["message"])
	}
}
func TestCreateProductInvalidStock(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.0,
		Stock:       -1,
	}

	productJSON, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(productJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.CreateProduct(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, utils.ErrProductStockInvalid.Error(), response["message"])
	}
}
func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		product        models.Product
		setupMock      func(*MockProductService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name: "Success",
			id:   "1",
			product: models.Product{
				Name:        "Updated Product",
				Description: "Updated Description",
				Price:       20.0,
				Stock:       10,
			},
			setupMock: func(m *MockProductService) {
				m.On("GetByID", "1").Return(models.Product{ProductID: "1"}, nil)
				m.On("UpdateProduct", mock.Anything, "1", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty ID",
			id:             "",
			setupMock:      func(m *MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    utils.ErrProductIDRequired.Error(),
		},
		{
			name: "Product Not Found",
			id:   "999",
			setupMock: func(m *MockProductService) {
				m.On("GetByID", "999").Return(models.Product{}, utils.ErrNoProductsFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedMsg:    utils.ErrNoProductsFound.Error(),
		},
		{
			name: "Invalid Price",
			id:   "1",
			product: models.Product{
				Price: -10,
			},
			setupMock: func(m *MockProductService) {
				m.On("GetByID", "1").Return(models.Product{ProductID: "1"}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    utils.ErrProductPriceInvalid.Error(),
		},
		{
			name: "Invalid Stock",
			id:   "1",
			product: models.Product{
				Price: 10,
				Stock: -1,
			},
			setupMock: func(m *MockProductService) {
				m.On("GetByID", "1").Return(models.Product{ProductID: "1"}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    utils.ErrProductStockInvalid.Error(),
		},
		{
			name: "Attempt to Change ProductID",
			id:   "1",
			product: models.Product{
				ProductID: "2",
				Price:     10,
				Stock:     5,
			},
			setupMock: func(m *MockProductService) {
				m.On("GetByID", "1").Return(models.Product{ProductID: "1"}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    utils.ErrProductIDCannotBeChanged.Error(),
		},
		{
			name: "Internal Server Error",
			id:   "1",
			product: models.Product{
				Price: 10,
				Stock: 5,
			},
			setupMock: func(m *MockProductService) {
				m.On("GetByID", "1").Return(models.Product{ProductID: "1"}, nil)
				m.On("UpdateProduct", mock.Anything, "1", mock.Anything).Return(utils.ErrInternalServer)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    utils.ErrInternalServer.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)
			tt.setupMock(mockService)
			handler := handlers.NewProductHandler(mockService)
			e := echo.New()

			productJSON, _ := json.Marshal(tt.product)
			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(productJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/products/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			err := handler.UpdateProduct(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedMsg != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedMsg, response["message"])
			}
		})
	}
}
func TestDeleteProduct(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockProductService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name: "Success",
			id:   "1",
			setupMock: func(m *MockProductService) {
				m.On("GetByID", "1").Return(models.Product{ProductID: "1"}, nil)
				m.On("DeleteProduct", mock.Anything, "1").Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Empty ID",
			id:             "",
			setupMock:      func(m *MockProductService) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    utils.ErrProductIDRequired.Error(),
		},
		{
			name: "Product Not Found",
			id:   "999",
			setupMock: func(m *MockProductService) {
				m.On("GetByID", "999").Return(models.Product{}, utils.ErrNoProductsFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedMsg:    utils.ErrNoProductsFound.Error(),
		},
		{
			name: "Internal Server Error",
			id:   "1",
			setupMock: func(m *MockProductService) {
				m.On("GetByID", "1").Return(models.Product{ProductID: "1"}, nil)
				m.On("DeleteProduct", mock.Anything, "1").Return(utils.ErrInternalServer)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    utils.ErrInternalServer.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)
			tt.setupMock(mockService)
			handler := handlers.NewProductHandler(mockService)
			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/products/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			err := handler.DeleteProduct(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedMsg != "" {
				var response map[string]string
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedMsg, response["message"])
			}
		})
	}
}

