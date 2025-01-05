## Installation Guide

### Prerequisites
- Ensure Docker and Docker Compose are installed on your machine. Download Docker from [here](https://www.docker.com/get-started).

### Steps to Install
1. Clone the repository to your local machine:
    ```sh
    git clone https://github.com/YugenDev/global-mobility-test.git
    cd global-mobility-test
    ```

2. Build and start the application using Docker Compose:
    ```sh
    docker-compose up --build
    ```

3. The application should now be running. Access it at:
    - `http://localhost:80/api/ecommerce`
    - `http://localhost:80/api/space-api`

### Usage Guide
- This backend is intended to be used through its API-Gateway (Traefik). View Traefik's dashboard at [http://localhost:8081/dashboard](http://localhost:8081/dashboard).

## Healthchecks
- Ecommerce `http://localhost:80/api/ecommerce`
- SpaceAPI `http://localhost:80/api/space-api`

- For detailed API documentation, Postman collections, and more usage details, refer to the [documentation module](docs/Docs.md).

### Additional Notes
- Check the `.env` file for any environment variables that need to be configured before running the application.
    ```sh
    NASA_API_KEY=YourAPIKey
    ```
    Get your own NASA API KEY [here](https://api.nasa.gov/).

- To stop the application, use:
    ```sh
    docker-compose down
    ```