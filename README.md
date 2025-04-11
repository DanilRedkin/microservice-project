# Library Service

This project implements a simple **Library Service** that allows managing books and authors.

## API

The service supports both REST and gRPC APIs. Validation rules are described below:

- **UUIDs:** Book and Author IDs must be in UUID format.
- **Book Name:** Cannot be empty (minimum length: 1 character).
- **Author Name:** Must match the regex `^[A-Za-z0-9]+( [A-Za-z0-9]+)*$` and have a length between 1 and 512 characters.
- **Authors per Book:** A book can have a maximum of 10 authors.

## REST to gRPC Mapping

The API supports REST-to-gRPC communication using gRPC Gateway. Below are the primary endpoints:

- **Add a book:** `POST /v1/library/book`
- **Update a book:** `PUT /v1/library/book`
- **Get book info:** `GET /v1/library/book/{book_id}`
- **Register an author:** `POST /v1/library/author`
- **Update author info:** `PUT /v1/library/author`
- **Get author info:** `GET /v1/library/author/{author_id}`
- **Get author's books:** `GET /v1/library/author_books/{author_id}`

## Environment Variables

The service requires the following environment variables:

- **GRPC_PORT** – Port for the gRPC server (default: `9090`).
- **GRPC_GATEWAY_PORT** – Port for the REST-to-gRPC API (default: `8080`).

## Database Configuration

**Note:** This service supports persistent storage via **PostgreSQL** (recommended for production) rather than using in-memory storage.

To enable PostgreSQL storage, set the following environment variables:

- **POSTGRES_HOST** – Hostname for the PostgreSQL server.
- **POSTGRES_PORT** – Port for the PostgreSQL server.
- **POSTGRES_DB** – Name of the database.
- **POSTGRES_USER** – PostgreSQL username.
- **POSTGRES_PASSWORD** – PostgreSQL password.
- **POSTGRES_MAX_CONN** – Maximum number of connections allowed to the PostgreSQL database.

When these PostgreSQL configuration variables are set, the service will use PostgreSQL as the storage backend for books and authors. If not set, an in-memory implementation may be used for testing or development purposes.

## Setting Up

For setting up the service on a Unix system and testing with Postman, use the following commands:

```sh
export GRPC_PORT=9090    
export GRPC_GATEWAY_PORT=8080

# Set the PostgreSQL environment variables as needed:
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_DB=librarydb
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=yourpassword
export POSTGRES_MAX_CONN=10

make generate  

go run cmd/library/main.go
```

This will start both the gRPC and REST services, allowing you to test API requests using Postman or other tools.

## Using Postman for API Testing

You can interact with the Library Service using Postman:

1. Open Postman and create a new request.
2. Set the request type to `POST`, `GET`, or `PUT` depending on the API you want to test.
3. Enter the request URL (e.g., for registering an author: `POST http://localhost:8080/v1/library/author`).
4. If the request requires a body (e.g., registering an author), navigate to the **Body** tab and choose **raw**.
5. Use JSON format. For example, to add an author, provide:

   ```json
   {
     "name": "Alexander Pushkin"
   }
   ```

6. Click **Send** to execute the request and check the response.
7. The expected response format is JSON. For example, after adding an author, you might receive:

   ```json
   {
     "id": "a645a5be-babc-4f35-80b2-073714d0be97"
   }
   ```

8. Use the returned Author ID to add a book. For example, for adding a book: `POST http://localhost:8080/v1/library/book`

   ```json
   {
     "name": "Dubrovsky",
     "authorIds": [
       "a645a5be-babc-4f35-80b2-073714d0be97"
     ]
   }
   ```

9. The response for adding a book will look like this:

   ```json
   {
     "id": "9a009f61-ea7d-4311-9594-88fbd81d23d3",
     "name": "Dubrovsky",
     "authorIds": [
       "a645a5be-babc-4f35-80b2-073714d0be97"
     ]
   }
   ```

10. Check the `library.proto` file to see the required arguments for each request.

## Makefile and Testing

For convenience in local development, a **Makefile** is included. The following commands are available:

- **Run linter and tests:**
  ```sh
  make all
  ```
- **Run only tests:**
  ```sh
  make test
  ```
- **Run the linter:**
  ```sh
  make lint
  ```
- **Pull new test updates:**
  ```sh
  make update
  ```
- **Generate new functionality:**
  ```sh
  make generate
  ```
  The `generate` command is required before running the service for the first time. It ensures that necessary code, including generated API and mock implementations for testing, is created.
```
