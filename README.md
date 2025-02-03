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

1. Clone the repository:
   ```sh
   git clone https://github.com/Blue-Davinci/socialaid-tracker.git
   cd socialaid-tracker
   ```

2. Set up the database:
   ```sh
   psql -U postgres -c "CREATE DATABASE socialaid_tracker;"
   ```

3. Configure environment variables (create a `.env` file):
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_user
   DB_PASSWORD=your_password
   DB_NAME=socialaid_tracker
   ```

4. Run the database migrations:
   ```sh
   go run migrations/migrate.go
   ```

5. Start the server:
   ```sh
   go run main.go
   ```

## Usage <a name = "usage"></a>

Once the server is running, you can interact with the API using Postman or Curl:

- Retrieve all households:
  ```sh
  curl -X GET http://localhost:8080/api/households
  ```

- Create a new beneficiary:
  ```sh
  curl -X POST http://localhost:8080/api/beneficiaries -d '{"name": "John Doe", "program_id": 1, "geolocation_id": 2}' -H "Content-Type: application/json"
  ```

For more details, refer to the API documentation.
