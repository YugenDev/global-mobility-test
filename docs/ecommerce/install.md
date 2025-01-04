# Running the E-commerce Microservice Container Individually

## Prerequisites
- Docker installed on your machine
- MongoDB instance running (either locally or in a container)

## Steps to Run the E-commerce Container Individually

1. **Build the Docker Image**
    Navigate to the `ecommerce/` directory where the `Dockerfile` is located and run the following command to build the Docker image:

    ```sh
    docker build -t ecommerce-service .
    ```

2. **Run the Docker Container**
    Once the image is built, run the container using the following command:

    ```sh
    docker run -d -p 8080:8080 --name ecommerce-container \
      -e MONGO_URI=mongodb://admin:globalMobility@mongodb:27017 \
      -e MONGO_DB_NAME=ecommerce \
      ecommerce-service
    ```

    This command will:
    - Start the container in detached mode (`-d`)
    - Map port 8080 of the host to port 8080 of the container (`-p 8080:8080`)
    - Name the container `ecommerce-container`
    - Set the necessary environment variables for MongoDB connection

3. **Access the Service**
    After the container is running, access the e-commerce service by navigating to `http://localhost:8080` in your web browser.

4. **Stop the Container**
    To stop the running container, use the following command:

    ```sh
    docker stop ecommerce-container
    ```

5. **Remove the Container**
    If you need to remove the container, use the following command:

    ```sh
    docker rm ecommerce-container
    ```

6. **View Logs**
    To view the logs of the running container, use:

    ```sh
    docker logs ecommerce-container
    ```

7. **Enter the Container**
    If you need to enter the running container for debugging or other purposes, use:

    ```sh
    docker exec -it ecommerce-container /bin/sh
    ```

## Next Steps

Refer to the [USAGE](usage.md) documentation for further instructions.
