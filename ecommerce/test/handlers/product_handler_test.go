package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/YugenDev/global-mobility-test/internal/handlers"
	"github.com/YugenDev/global-mobility-test/internal/models"
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

func TestUpdateProductInvalidData(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := &models.Product{
		Price: -1,
		Stock: -1,
	}

	productJSON, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(productJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/products/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if err := handler.UpdateProduct(c); assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreateProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.0,
		Stock:       5,
	}

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

func TestUpdateProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	product := &models.Product{
		Name:  "Updated Product",
		Price: 20.0,
		Stock: 10,
	}

	mockService.On("UpdateProduct", mock.Anything, "1", product).Return(nil)

	productJSON, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(productJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/products/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handler.UpdateProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := handlers.NewProductHandler(mockService)
	e := echo.New()

	mockService.On("DeleteProduct", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/products/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handler.DeleteProduct(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
