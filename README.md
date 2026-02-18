# CRUD App with Go Fiber & PostgreSQL

A high-performance Todo application built with Go 1.26, Fiber v3, and PostgreSQL. Features user authentication, session management, and CSRF protection.

## Features

- **User Authentication**: Secure Login & Registration with password hashing (bcrypt).
- **Todo Management**: Create, Read, and Delete tasks.
- **Security**: 
  - CSRF Protection
  - Secure Session Management
  - Helmet Middleware for security headers
- **Architecture**: Clean architecture with separation of concerns (Handlers, Repositories, Models).
- **Database**: PostgreSQL integration using `pgx/v5`.
- **Templates**: Server-side rendering with HTML templates.

## Tech Stack

- **Language**: Go 1.26.0
- **Framework**: [Fiber v3](https://github.com/gofiber/fiber)
- **Database**: PostgreSQL
- **Driver**: [pgx/v5](https://github.com/jackc/pgx)
- **Templates**: [html/v2](https://github.com/gofiber/template)

## Prerequisites

- [Go](https://go.dev/dl/) 1.26+
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

1. **Clone the repository**
   ```bash
   git clone git@github.com:bayazidsustami/golang-todo-fullstack.git
   cd go-fiber-crud
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure Environment Variables**
   Create a `.env` file in the root directory:
   ```env
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASS=password
   DB_NAME=fiber_app
   DB_PORT=5432
   SERVER_PORT=:8090
   ```

4. **Setup Database**
   Create a PostgreSQL database named `fiber_app` (or whatever you set in `.env`) and run the following SQL commands to create the schema:

   ```sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       username TEXT UNIQUE NOT NULL,
       password TEXT NOT NULL
   );

   CREATE TABLE todos (
       id SERIAL PRIMARY KEY,
       user_id INT REFERENCES users(id) ON DELETE CASCADE,
       title TEXT NOT NULL
   );
   ```

## Running the Application

Start the server:
```bash
go run cmd/web/main.go
```
The application will be available at `http://localhost:8090`.

## Project Structure

```
.
├── cmd/web/main.go        # Entry point
├── internal
│   ├── config             # Configuration logic
│   ├── database           # Database connection
│   ├── handlers           # HTTP Handlers
│   ├── models             # Data models
│   └── repository         # Data access layer
├── views                  # HTML Templates
├── go.mod                 # Go module file
└── .env                   # Environment variables
```
