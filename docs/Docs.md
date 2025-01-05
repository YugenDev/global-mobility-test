# Global Mobility Apex Technical Test

This documentation provides an overview of the endpoints available in the Global Mobility Apex Technical Test. The test includes two main services: Ecommerce and SpaceAPI, as well as an API Gateway to consume both services.

## Ecommerce CRUD

Basic CRUD operations for the e-commerce microservice.

### Add a product

- **Method:** POST
- **URL:** `http://localhost:8080/products`
- **Description:** This endpoint allows you to add a product. You can either manually specify a product ID or omit the `product_id` in the request body to have it automatically generated as a UUID.
- **Request Body:**
    ```json
    {
        "product_id": "123",
        "name": "Testing 2",
        "description": "yeeehaaaaaa",
        "price": 29.99,
        "stock": 50
        
    }
    ```

### Get all products

- **Method:** GET
- **URL:** `http://localhost:8080/products`
- **Description:** This endpoint allows you to get a list of all products.

### Get a product by ID

- **Method:** GET
- **URL:** `http://localhost:8080/products/{id}`
- **Description:** This endpoint allows you to retrieve a specific product by searching using its ID.

### Delete a product by ID

- **Method:** DELETE
- **URL:** `http://localhost:8080/products/{id}`
- **Description:** This endpoint allows you to delete a specific product by searching using its ID.

### Modify a product

- **Method:** PUT
- **URL:** `http://localhost:8080/products/{id}`
- **Description:** This endpoint allows you to modify a specific product by searching using its ID.
- **Request Body:**
    ```json
    {
        "name": "Testing change",
        "description": "yeeehaaaaaa",
        "price": 29.99,
        "stock": 50
    }
    ```

### Healthcheck

- **Method:** GET
- **URL:** `http://localhost:8080`
- **Description:** This endpoint allows you to check if the service is up.

## SPACE API

Operations to consume the NASA Space API for retrieving astronomical photos.

### Healthcheck

- **Method:** GET
- **URL:** `http://localhost:8000`
- **Description:** This endpoint allows you to check if the service is up.

### Get APOD JSON

- **Method:** GET
- **URL:** `http://localhost:8000/apod`
- **Description:** This endpoint retrieves an astronomical picture along with its information in JSON format.

### Get APOD HTML

- **Method:** GET
- **URL:** `http://localhost:8000/apod/html`
- **Description:** This endpoint retrieves an astronomical picture along with its information in HTML format.

## API-GATEWAY

Consume both the E-commerce and Space APIs through the Traefik API Gateway.

### SPACE API

#### Get APOD JSON

- **Method:** GET
- **URL:** `http://localhost:80/api/space-api/apod`

#### Get APOD HTML

- **Method:** GET
- **URL:** `http://localhost:80/api/space-api/apod/html`

#### Healthcheck

- **Method:** GET
- **URL:** `http://localhost:80/api/space-api`

### ECOMMERCE

#### Add a product

- **Method:** POST
- **URL:** `http://localhost:80/api/ecommerce/products`
- **Request Body:**
    ```json
    {
        "product_id": "123",
        "name": "Testing",
        "description": "This is an outstanding description",
        "price": 29.99,
        "stock": 50
    }
    ```

#### Get all products

- **Method:** GET
- **URL:** `http://localhost:80/api/ecommerce/products`
- **Headers:**
    ```json
    {
        "Content-Type": "application/json"
    }
    ```

#### Get a product by ID

- **Method:** GET
- **URL:** `http://localhost:80/api/ecommerce/products/{id}`

#### Delete a product by ID

- **Method:** DELETE
- **URL:** `http://localhost:80/api/ecommerce/products/{id}`

#### Modify a product

- **Method:** PUT
- **URL:** `http://localhost:80/api/ecommerce/products/{id}`
- **Request Body:**
    ```json
    {
        "name": "Testing change",
        "description": "yeeehaaaaaa",
        "price": 29.99,
        "stock": 50
    }
    ```

#### Healthcheck

- **Method:** GET
- **URL:** `http://localhost:80/api/ecommerce`