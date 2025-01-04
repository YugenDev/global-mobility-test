# SpaceAPI Installation Guide

## Overview
SpaceAPI is a RESTful API service designed to manage space-related data. This guide provides instructions on how to install, run, and use the SpaceAPI service.

## Prerequisites
- Docker installed on your machine
- MongoDB instance running (either locally or in a container)
- NASA API Key (for accessing NASA's APOD service)

## Installation

### Steps to Run the SpaceAPI Container Individually

1. **Build the Docker Image**
    Navigate to the directory 

spaceAPI

 where the 

Dockerfile

 is located and run the following command to build the Docker image:
    ```sh
    docker build -t spaceapi-service .
    ```

2. **Run the Docker Container**
    Once the image is built, you can run the container using the following command:
    ```sh
    docker run -d -p 8000:8000 --name spaceapi-container \
      -e NASA_API_KEY=your_nasa_api_key
    ```
    This command will start the container in detached mode (`-d`), map port 8000 of the host to port 8000 of the container (`-p 8000:8000`), and name the container spaceapi-container. It also sets the necessary environment variable NASA API Key. 
    
    In case you dont have one get your Nasa API key [HERE](https://api.nasa.gov/)

3. **Access the Service**
    After the container is running, you can access the SpaceAPI service by navigating to `http://localhost:8000` in your web browser.

4. **Stopping the Container**
    To stop the running container, use the following command:
    ```sh
    docker stop spaceapi-container
    ```

5. **Removing the Container**
    If you need to remove the container, you can do so with the following command:
    ```sh
    docker rm spaceapi-container
    ```

6. **Viewing Logs**
    To view the logs of the running container, use:
    ```sh
    docker logs spaceapi-container
    ```

7. **Entering the Container**
    If you need to enter the running container for debugging or other purposes, use:
    ```sh
    docker exec -it spaceapi-container /bin/sh
    ```