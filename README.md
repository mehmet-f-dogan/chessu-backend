# ChessU - Chess Course Platform - Backend

This repository contains the backend code for a Chess Course Platform developed using Go programming language. The backend utilizes the GORM library for interacting with a PostgreSQL database. The application is hosted on an EC2 instance and integrates with an S3 bucket for file storage. The server is built using the Fiber framework to handle incoming HTTP requests.

## Features

- CRUD operations for managing chess courses, lessons, and user enrollments
- Secure authentication and authorization
- Integration with PostgreSQL database using GORM
- File storage and retrieval using an S3 bucket
- RESTful API endpoints for communication with the frontend

## Technologies Used

The backend of the Chess Course Platform is built using the following technologies:

- Go: A programming language used for building high-performance applications.
- GORM: An ORM library for interacting with databases in Go.
- PostgreSQL: An open-source relational database management system.
- EC2: Amazon Elastic Compute Cloud for hosting the backend server.
- S3 Bucket: Amazon Simple Storage Service for securely storing and retrieving files.
- Fiber: A web framework for building fast and efficient HTTP APIs in Go.

## Getting Started

To run the Chess Course Platform backend locally, follow these steps:

1. Clone this repository to your local machine.
2. Install Go and set up your Go development environment.
3. Install the necessary dependencies by running the following command:

   ```
   go get -u ./...
   ```

4. Set up a PostgreSQL database and update the database connection details in the configuration file.
5. Configure the necessary environment variables. You will need to provide the following:

   - PostgreSQL database connection details: Update the database URL, username, password, and other required information.
   - S3 Bucket details: Configure the S3 bucket information for file storage.

6. Build and run the server using the following command:

   ```
   go run main.go
   ```

7. The server should now be running at `http://localhost:8000`.

## Deployment

The Chess Course Platform backend can be deployed on an EC2 instance using the following steps:

1. Set up an EC2 instance on AWS.
2. Install Go and configure your environment on the EC2 instance.
3. Clone the repository onto the EC2 instance.
4. Install the necessary dependencies as mentioned in the "Getting Started" section.
5. Configure the necessary environment variables on the EC2 instance.
6. Build and run the server on the EC2 instance.

Ensure that you have configured the appropriate security groups, firewalls, and other necessary settings to allow incoming traffic to the EC2 instance.
