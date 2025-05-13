# bugutv-autosign

This project automates the sign-in process for the bugutv.vip website to earn points.

## Prerequisites

-   Go 1.24 or higher
-   Docker (if you want to build and run the Docker image)

## Usage

### Configuration

Set the following environment variables:

-   `BUGUTV_USERNAME`: Your bugutv.vip username.
-   `BUGUTV_PASSWORD`: Your bugutv.vip password.

### Running Locally

1.  Clone the repository:

    ```bash
    git clone https://github.com/kingcanfish/bugutv.git
    cd bugutv
    ```

2.  Set the environment variables in your shell or `.env` file.

3.  Run the program:

    ```bash
    go run main.go
    ```

### Running with Docker

1.  Build the Docker image:

    ```bash
    docker build -t bugutv-autosign .
    ```

2.  Run the Docker container, passing the environment variables:

    ```bash
    docker run -d --name bugutv-autosign \
      -e BUGUTV_USERNAME=<your_username> \
      -e BUGUTV_PASSWORD=<your_password> \
      bugutv-autosign
    ```

    Replace `<your_username>` and `<your_password>` with your actual bugutv.vip credentials.

### Docker Image

The project includes a `Dockerfile` for containerizing the application.  The `.github/workflows/docker-build-push.yml` file configures a GitHub Actions workflow to automatically build and push Docker images to Docker Hub when changes are pushed to the `main` branch.  The workflow uses secrets `DOCKER_IMAGE_NAME`, `DOCKER_USERNAME`, and `DOCKER_PASSWORD` to configure the Docker image name and authentication.

## GitHub Actions

The project includes a GitHub Actions workflow (`.github/workflows/docker-build-push.yml`) that automates the building and pushing of the Docker image to Docker Hub.

## Disclaimer

This project is intended for personal use only. Use it at your own risk. I am not responsible for any consequences resulting from the use of this tool.