## Endpoints

### Base URL & Healthcheck
`http://localhost:8080`

### Product Routes

#### Create a Product
- **URL**: `/products`
- **Method**: `POST`
- **Request Body**:
    ```json
    {
      "product_id": "id123",
      "name": "Product Name",
      "description": "Product Description",
      "price": 100.0,
      "stock": 50
    }
    ```
    - **Note**: If you omit the `"product_id"` field, a UUID will be generated automatically.
- **Response**:
    - **201 Created**: Product created successfully
    - **400 Bad Request**: Invalid request payload

#### Get All Products
- **URL**: `/products`
- **Method**: `GET`
- **Response**:
    - **200 OK**: List of products
    - **404 Not Found**: No products found

#### Get Product by ID
- **URL**: `/products/:id`
- **Method**: `GET`
- **Response**:
    - **200 OK**: Product details
    - **404 Not Found**: Product not found

#### Update a Product
- **URL**: `/products/:id`
- **Method**: `PUT`
- **Request Body**:
    ```json
    {
      "name": "Updated Product Name",
      "description": "Updated Product Description",
      "price": 150.0,
      "stock": 30
    }
    ```
- **Response**:
    - **200 OK**: Product updated successfully
    - **400 Bad Request**: Invalid request payload
    - **404 Not Found**: Product not found

#### Delete a Product
- **URL**: `/products/:id`
- **Method**: `DELETE`
- **Response**:
    - **204 No Content**: Product deleted successfully
    - **404 Not Found**: Product not found