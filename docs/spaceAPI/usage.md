## Endpoints

### Base URL & Healthcheck
`http://localhost:8000`

### Swagger Documentation
`http://localhost:8000/docs`

### Space Data Routes

#### Get Astronomy Picture of the Day (APOD)
- **URL**: `/apod`
- **Method**: `GET`
- **Response**:
    ```json
    {
        "request": "request",
        "title": "title",
        "date": "date",
        "explanation": "explanation",
        "url": "url",
        "media_type": "media_type"
    }
    ```
- **Status Codes**:
    - **200 OK**: Successfully returns the APOD data.
    - **500 Internal Server Error**: Error occurred while fetching the data.

#### Get APOD as HTML
- **URL**: `/apod/html`
- **Method**: `GET`
- **Response**:
    - **200 OK**: Successfully returns the APOD data rendered as HTML.
    - **500 Internal Server Error**: Error occurred while fetching the data.
