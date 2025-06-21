# Kumparan Article Go

This is a RESTful API project built with Go, designed with layered architecture. The project includes full CRUD functionality for users and articles, with a JWT authentication system (Access & Refresh Tokens) and caching using Redis.

## Key Features

* **User Management**: Register, Login, Update, Delete, and Get (all/specific).
* **Article Management**: Full CRUD (Create, Read, Update, Delete).
* **JWT Authentication**: Utilizes short-lived Access Tokens and long-lived Refresh Tokens for security.
* **Authorization**: Users can only modify or delete their own articles and profiles.
* **Caching**: Uses Redis to cache frequently accessed endpoints (like article details) to improve performance.
* **Pagination**: The article list endpoint supports pagination (`page` & `limit`).
* **Full-Text Search**: Ability to search for articles by keywords in the title and body.
* **Development Ready**: Comes with `docker-compose` for easy environment setup and live-reloading using **Air**.

## Technology Stack

* **Language**: Go (Golang) 1.24
* **Database**: PostgreSQL
* **Cache**: Redis
* **Web Framework/Router**: Gorilla Mux
* **Database Driver**: `pgx/v5`
* **Authentication**: `golang-jwt/jwt/v5`
* **Password Hashing**: `golang.org/x/crypto/bcrypt`
* **Containerization**: Docker & Docker Compose
* **Live Reload**: Air

## Project Structure

```

/
├── cmd/app/main.go             \# Main application entry point
├── internal/
│   ├── handlers/               \# HTTP Controllers
│   ├── models/                 \# Data structs (entities & DTOs)
│   ├── repositories/           \# Data access logic (SQL queries)
│   ├── router/                 \# Route definitions separated by domain
│   └── services/               \# Core business logic
├── pkg/
│   ├── middleware/             \# Middleware (JWT)
│   └── utils/                  \# Helper functions (response, password, token)
├── db/init.sql                 \# Database initialization schema
├── .env                        \# Configuration file (NOT committed to git)
├── .air.toml                   \# Configuration for live-reload
├── docker-compose.yml          \# Docker services orchestration
├── Dockerfile                  \# Docker image build instructions
└── go.mod                      \# Project dependencies

````

## How to Run

### 1. Prerequisites

* [Docker](https://www.docker.com/products/docker-desktop/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [Go](https://go.dev/doc/install) (version 1.24)

### 2. Environment Setup

1.  **Clone the Repository**
    ```bash
    git clone https://github.com/dhifanrazaqa/kumparan-article.git
    cd kumparan-article
    ```

2.  **Create `.env` File**
    Create a file named `.env` in the project's root directory. Copy the contents below into it. This file holds all the necessary application configurations.
    ```env
    # Application Configuration
    APP_PORT=8080
    
    # Database Configuration
    DATABASE_URL=postgres://user:password@postgres:5432/articledb?sslmode=disable
    TEST_DATABASE_URL="postgres://user:password@postgres:5432/articledb_test?sslmode=disable"

    # Postgres Container Configuration
    POSTGRES_USER="user"
    POSTGRES_PASSWORD="password"
    POSTGRES_DB="articledb"
    
    # Redis Configuration
    REDIS_URL=redis:6379
    
    # JWT Secret Keys (Replace with strong, random values)
    JWT_SECRET_KEY=a-very-secret-key-for-your-access-tokens
    REFRESH_TOKEN_SECRET=another-very-secret-key-for-refresh-tokens
    ```

### 3. Run the Application

Execute the following command from the project's root directory.
```bash
docker-compose up --build
````

  * This command builds the Go image, then starts the API, PostgreSQL, and Redis containers.
  * The database is automatically initialized using `db/init.sql`.
  * The API will be running at `http://localhost:8080`.
  * Any changes to `.go` files will automatically restart the server (live-reload).

-----

## API Endpoint Documentation

### Authentication (`/auth`)

| Method | Endpoint         | Description                                        | Request Body                                     |
| :----- | :--------------- | :------------------------------------------------- | :----------------------------------------------- |
| `POST` | `/auth/login`    | Logs in to get an Access and Refresh Token.        | `{"username": "...", "password": "..."}`         |
| `POST` | `/auth/refresh`  | Gets a new Access Token using a Refresh Token.     | `{"refreshToken": "..."}`                        |

### Users (`/users`)

| Method   | Endpoint          | Description                                         | Authorization Header | Request Body                                                  |
| :------- | :---------------- | :-------------------------------------------------- | :------------------- | :------------------------------------------------------------ |
| `POST`   | `/users`          | Registers a new user.                               |                      | `{"name": "Full Name", "username": "...", "password": "..."}` |
| `GET`    | `/users`          | Gets a list of all users.                           | -                    | -                                                             |
| `GET`    | `/users/{id}`     | Gets details for a single user by ID.               | -                    | -                                                             |
| `PUT`    | `/users/{id}`     | Updates a user's profile (only owner can perform).  | `Bearer <token>`     | `{"username": "(optional)", "name": "(optional)", "password": "(optional)"}`            |
| `DELETE` | `/users/{id}`     | Deletes a user's account (only owner can perform).  | `Bearer <token>`     | -                                                             |

### Articles (`/articles`)

| Method   | Endpoint           | Description                                       | Authorization Header | Request Body                                    | Optional Query Params          |
| :------- | :----------------- | :------------------------------------------------ | :------------------- | :---------------------------------------------- | :----------------------------- |
| `POST`   | `/articles`        | Creates a new article.                            | `Bearer <token>`     | `{"title": "...", "body": "..."}`                | -                              |
| `GET`    | `/articles`        | Gets a paginated list of articles.                | -                    | -                                               | `page`, `limit`, `author`, `query` |
| `GET`    | `/articles/{id}`   | Gets details for a single article by ID.          | -                    | -                                               | -                              |
| `PUT`    | `/articles/{id}`   | Updates an article (only original author can perform). | `Bearer <token>`     | `{"title": "(optional)", "body": "(optional)"}` | -                              |
| `DELETE` | `/articles/{id}`   | Deletes an article (only original author can perform). | `Bearer <token>`     | -                                               | -                              |
