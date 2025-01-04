## Installation Guide

### Prerequisites
- Ensure you have Docker and Docker Compose installed on your machine. You can download Docker from [here](https://www.docker.com/get-started).

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

3. The application should now be running. You can access it at `http://localhost:80/api/ecommerce`.

### Usage Guide
- This backend is meant to be used through its API-Gateway (Traefik). You can view Traefik's dashboard at `http://localhost:8081/dashboard`.

- For detailed API documentation, Postman collections, and more usage details, refer to the [documentation module](docs/index.md).

### Additional Notes
- Ensure to check the `.env` file for any environment variables that need to be configured before running the application.
- To stop the application, use:
    ```sh
    docker-compose down
    ```