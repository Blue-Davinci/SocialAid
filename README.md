<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=190px src="https://i.ibb.co/LhnG5QWC/socialaid-high-resolution-logo.png" alt="Project logo"></a>
</p>


# SocialAid Tracker

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Contributing](../CONTRIBUTING.md)

## About <a name = "about"></a>

SocialAid Tracker is a Management Information System (MIS) designed to track beneficiaries of social protection programs. It enables efficient storage and retrieval of program details, household geolocation data, and household member information. The system ensures data security by encrypting sensitive data such as phone numbers while providing seamless access through RESTful API endpoints.

## Getting Started <a name = "getting_started"></a>

These instructions will help you set up SocialAid Tracker on your local machine for development and testing. See [deployment](#deployment) for details on deploying the system in a production environment.

### Prerequisites

Ensure you have the following installed before proceeding:

```
- Go (1.18 or later)
- PostgreSQL (latest version)
- Docker (optional for containerized setup)
- Git
```

### Installing

Follow these steps to set up the development environment:

1. **Clone the repository:**
   ```sh
   git clone https://github.com/Blue-Davinci/socialaid.git
   cd socialaid-tracker
   ```

2. **Set up the database**:
   ```sh
   psql -U postgres -c "CREATE DATABASE socialaid;"
   ```

3. **Install the Go dependencies:** The Go tools will automatically download and install the dependencies listed in the `go.mod` file when you build or run the project. To download the dependencies without building or running the project, you can use the `go mod download` command:

    ```bash
    go mod download
    ```

4. **Set up the database:** The project uses a PostgreSQL database. You'll need to create a new database and update the connection string in your configuration file or environment variables.
We use `GOOSE` for all the data migrations and `SQLC` as the abstraction layer for the DB. To proceed
with the migration, navigate to the `Schema` director:
```bash
cd internal\sql\schema
```
- Then proceed by using the `goose {connection string} up` to execute an <b>Up migration</b> as shown:
- <b>Note:</b> You can use your own environment variable or load it from the env file.

```bash
goose postgres postgres://socialaid:password@localhost/socialaid  up
```


5. **Configure environment variables (create a `.env` file):**
   ```bash
      SOCIALAID_DB_DSN= ADD YOUR POSTGRESQL DSN HERE
      SOCIALAID_DATA_ENCRYPTION_KEY= ADD YOUR HEX ENCODED ENCRYPTION KEY STRING HERE, you can use the generateEncryptioKey() func in `internal_helpers.go` to generate one.
   ```

5. **Build the project:** You can build the project using the makefile's command:

    ```bash
    make build/api
    ```
    This will create an executable file in the current directory.
    <b>Note: The generated executable is for the windows environment</b>

6. **Run the project:** You can run the project using the `go run ./cmd/api` or use <b>`MakeFile`</b> and do:

    ```bash
    make run/api

## Usage <a name = "usage"></a>

Once the server is running, you can interact with the API using Postman or Curl:

- Retrieve single household info:
  ```sh
  curl -X GET http://localhost:8080/v1/house_holds/-household_ID-
  ```

- Create a new household beneficiary:
  ```sh
  curl -X POST http://localhost:8080/v1/house_holds/member -d '{"name": "John Doe", "program_id": 1, "geolocation_id": 2}' -H "Content-Type: application/json"
  ```

- For authentication, check the `006 users.sql` migration file for the plain text token example. You will need to select the ApiKey method for authorization
using the key: `ApiKey` and the **Token** as the Value.
- This will allow you to access all `\house_holds` routes

For more details, refer to the API documentation.
