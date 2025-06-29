# 🚀 Golang Clean Architecture Starterpack

This is a Golang backend starterpack built with **Clean Architecture** principles (also known as Hexagonal Architecture or Port and Adapters). This starterpack is designed to provide a solid, modular, and easy-to-maintain foundation for Go web applications, focusing on *testability*, *scalability*, and *developer experience*.

It is suitable for projects that require a well-organized structure from the outset and are flexible enough to evolve, including multi-tenancy implementation and potential microservices adaptation in the future.

## ✨ Key Features

* **Clean Architecture:** Clear project structure separating `domain`, `service` (usecase), `handler` (transport), and `repository` (persistence) layers for robust, testable, and maintainable code.
* **Go Modules:** Modern Go dependency management for reliable builds.
* **Structured Logging:** Uses `go.uber.org/zap` for consistent, high-performance, and structured logging, crucial for debugging and monitoring.
* **Custom Error Handling:** A structured custom error system (`internal/utils/errors`) with consistent mapping to HTTP status codes, ensuring uniform API error responses.
* **Pagination Helper:** Generic utility (`internal/utils/response`) and DTOs (Data Transfer Objects) for consistent pagination implementation across various list endpoints.
* **Database Ready:** Initial configuration for PostgreSQL with SQLX, equipped with robust schema setup via `psql` directly in Docker for enhanced reliability during local development.
    * **Note on Schema:** Uses `id SERIAL PRIMARY KEY` (auto-incrementing integer) and `tenant_id VARCHAR(36)` (for UUID strings) in the initial schema for broader compatibility and to avoid `uuid-ossp` extension setup complexity in basic setups.
* **Docker Support:** `Dockerfile` for application containerization and Docker (via Makefile) for easy local PostgreSQL database setup and management.
* **Makefile:** Comprehensive automation scripts for common development tasks (build, test, lint, Docker commands, database management including full resets).
* **GCP Ready:** Architecture is designed to be highly compatible for seamless deployment to Google Cloud Platform services (e.g., Cloud Run, Cloud SQL, Kubernetes).
* **OpenAPI (Swagger) Documentation:** Includes the `api/` directory for API specifications (`openapi.yaml`), essential for generating and visualizing comprehensive API documentation.
* **Demo API Endpoints:** Features basic `/health` and `/info` endpoints for quick application health verification, plus a **full CRUD demo for the "Employee" module (Create, Get All with Pagination, Get By ID)** to illustrate the consistent code structure and interaction patterns.

## 🛠️ Technology Stack

* **Programming Language:** Go (Golang)
* **Web Framework:** [Gorilla Mux](https://github.com/gorilla/mux) (a lightweight and flexible HTTP router, ideal for Clean Architecture. Easily swappable with other frameworks like Fiber or Gin if preferred later.)
* **Database:** PostgreSQL (a powerful, open-source relational database)
* **ORM/Query Builder:** [SQLX](https://github.com/jmoiron/sqlx) (an extension to Go's `database/sql` package, offering a balance between raw SQL power and easy struct mapping).
* **Logging:** [Zap](https://github.com/uber-go/zap) (a blazing fast, structured, leveled logging library).
* **Validation:** [go-playground/validator](https://go-playground.github.io/validator) (a powerful and extensible Go Struct validation package).
* **Database Migrations (Advanced/Future):** [golang-migrate/migrate](https://golang-migrate.run/) (a database migration tool; while schema setup for this starterpack's demo is done via direct `psql` execution for simplicity, `golang-migrate` is included in the Makefile for future, more complex migration management).
* **Containerization:** Docker (industry-standard for packaging and running applications).
* **Linter:** [golangci-lint](https://golangci-lint.run/) (a fast Go linters runner, ensuring code quality and consistency).

## 🚀 Getting Started (Local Setup)

Follow the steps below to run the starterpack on your local machine.

### Prerequisites

Before starting, ensure you have the following tools installed and properly configured:

* **Git:** [Download & Install Git](https://git-scm.com/downloads) (for version control).
* **Visual Studio Code (VS Code):** [Download & Install VS Code](https://code.visualstudio.com/) (recommended IDE).
    * **Essential VS Code Extensions:**
        * [Go](https://marketplace.visualstudio.com/items?itemName=golang.go) (by Go Team at Google): Provides rich language support for Go.
        * [Remote - WSL](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-wsl) (by Microsoft): **Crucial for Windows developers**, allowing you to develop in a Linux environment directly from VS Code.
        * [Docker](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-docker) (by Microsoft): For managing Docker containers and images from VS Code.
* **Docker Desktop:** [Download & Install Docker Desktop](https://www.docker.com/products/docker-desktop/) (for running containers).
    * **Configuration:** After installation, open Docker Desktop Settings -> Resources -> WSL Integration. Ensure **"Enable integration with my default WSL distro"** is checked, and the specific **toggle for your Linux distribution (e.g., Ubuntu) is enabled**.
* **WSL 2 (Windows Subsystem for Linux) & Linux Distribution:**
    * Open **PowerShell/CMD as Administrator** and run: `wsl --install` (This typically installs WSL 2 and Ubuntu by default).
    * Follow the on-screen instructions to create your Linux username and password.

### Setup Steps (Follow Exactly)

To avoid common setup issues with `Makefile` and environment variables, follow these steps precisely. We will create critical files directly in the terminal to guarantee their integrity.

1.  **Perform a Clean Start (Crucial):**
    * **Close ALL WSL terminals.**
    * Open **PowerShell (NOT WSL)** as Administrator.
    * Stop and remove any lingering Docker containers/volumes:
        ```powershell
        docker stop starterpack-golang-cleanarch-db || Write-Host "Container not running"
        docker rm starterpack-golang-cleanarch-db || Write-Host "Container not found"
        docker volume prune -f || Write-Host "No volumes to prune"
        ```
    * **Delete your existing project folder:** Navigate to `D:\Playground\` (or wherever your project is) in **Windows File Explorer** and **delete the `starterpack-golang-cleanarch` folder completely**.

2.  **Initialize New Project Folder in WSL:**
    * Open a **new WSL terminal** (e.g., Ubuntu).
    * Navigate to your desired directory (e.g., `cd /mnt/d/Playground/` or `cd ~` for Linux filesystem, recommended for performance).
    * ```bash
        mkdir starterpack-golang-cleanarch
        cd starterpack-golang-cleanarch
        ```

3.  **Initialize Go Modules:**
    ```bash
    go mod init starterpack-golang-cleanarch
    ```

4.  **Create Project Directories:**
    ```bash
    mkdir -p cmd/server
    mkdir -p internal/{app,domain,repository,utils/errors,utils/log,platform/http/middleware}
    mkdir -p migrations
    mkdir -p api
    mkdir -p scripts
    ```

5.  **Install Go SDK in WSL (if not already done correctly):**
    * `sudo apt update && sudo apt upgrade -y`
    * `wget https://go.dev/dl/go1.22.4.linux-amd64.tar.gz` (Check [go.dev/dl/](https://go.dev/dl/) for latest version)
    * `sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.4.linux-amd64.tar.gz`
    * `echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile && source ~/.profile`
    * Verify: `go version`

6.  **Install `make`, `curl`, `tar` in WSL (if not already done):**
    * `sudo apt install make curl tar -y`
    * Verify: `make --version`

7.  **Create `.env` File (GUARANTEED No Trailing Spaces):**
    * **Copy the entire block below (including `cat << EOF > .env` and `EOF`) and paste it directly into your WSL terminal.** This method ensures no hidden characters or trailing spaces.
    ```bash
    cat << EOF > .env
    # Application Environment
    APP_ENV=development # Options: production, development, testing
    APP_NAME="starterpack-golang-cleanarch"

    # Server Port
    PORT=8080

    # Database Configuration (PostgreSQL)
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=starteruser
    DB_PASSWORD=supersecret_db_password
    DB_NAME=starterdb

    # JWT Configuration (Example, adjust as needed for actual authentication)
    JWT_SECRET=this-is-a-super-secret-jwt-key-please-change-me-in-production
    JWT_EXPIRES_IN_MINUTES=60
    REFRESH_TOKEN_EXPIRES_IN_HOURS=720 # 30 days
    EOF
    ```

8.  **Create `Makefile` (GUARANTEED Correct Tabs and No Trailing Spaces):**
    * **Copy the entire content of the `Makefile` block provided below (including `cat << 'EOF' > Makefile` and `EOF`).** Use the "Copy" button for the code block.
    * **Paste it directly into your WSL terminal.** This method ensures correct tab indentation and no trailing spaces.
    ```bash
    cat << 'EOF' > Makefile
    # Makefile
    # This line MUST be the very first non-comment line in your Makefile
    # It loads environment variables from .env file directly into Makefile's scope.
    include .env
    export $(shell sed 's/=.*//' .env) # Exports all variables loaded from .env to the shell for commands

    .PHONY: build run test clean lint migrate-up migrate-down docker-build docker-run-db docker-stop-db help docker-reset-db

    APP_NAME := starterpack-golang-cleanarch
    BINARY_NAME := $(APP_NAME)
    BUILD_DIR := bin

    # Go related variables
    GO_VERSION := 1.22
    GO_LINTER := golangci-lint
    GO_MIGRATE := migrate # Still keep it for other commands if needed, but not for primary table creation now

    # Docker variables for local PostgreSQL (these values will be overridden by .env if set)
    DOCKER_IMAGE := $(APP_NAME)
    DOCKER_CONTAINER_NAME := $(APP_NAME)-db
    DOCKER_DB_PORT := 5432:5432
    DOCKER_DB_PASSWORD := $(DB_PASSWORD)
    DOCKER_DB_USER := $(DB_USER)
    DOCKER_DB_NAME := $(DB_NAME)

    # Default target
    all: build

    # Build the application binary
    build:
    	@echo "Building $(APP_NAME)..."
    	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

    # Run the application
    run: build
    	@echo "Running $(APP_NAME)..."
    	./$(BUILD_DIR)/$(BINARY_NAME)

    # Run tests
    test:
    	@echo "Running tests..."
    	go test -v ./...

    # Clean build artifacts
    clean:
    	@echo "Cleaning build artifacts..."
    	rm -rf $(BUILD_DIR)

    # Run linters (installs if not present)
    lint:
    	@echo "Running linters..."
    	@if ! command -v $(GO_LINTER) &> /dev/null; then \
    		echo "$(GO_LINTER) not found. Installing..."; \
    		go install [github.com/golangci/golangci-lint/cmd/golangci-lint@latest](https://github.com/golangci/golangci-lint/cmd/golangci-lint@latest); \
    	fi
    	$(GO_LINTER) run ./...

    # Docker commands
    docker-build:
    	@echo "Building Docker image..."
    	docker build -t $(DOCKER_IMAGE) .

    docker-run-db: docker-stop-db
    	@echo "Starting PostgreSQL Docker container..."
    	docker run --name $(DOCKER_CONTAINER_NAME) -e POSTGRES_PASSWORD=$(DOCKER_DB_PASSWORD) -e POSTGRES_USER=$(DOCKER_DB_USER) -e POSTGRES_DB=$(DOCKER_DB_NAME) -p $(DOCKER_DB_PORT) -d postgres:15-alpine
    	@echo "Waiting for database to be responsive..."
    	# Loop until the database is ready to accept connections (initial check)
    	until docker exec $(DOCKER_CONTAINER_NAME) pg_isready -U $(DOCKER_DB_USER) -d $(DOCKER_DB_NAME); do \
    		sleep 1; \
    	done;
    	@echo "Database is initially responsive. Waiting for full startup sequence..."
    	sleep 15 # Wait for 15 seconds to let the DB fully settle
    	@echo "Database is ready. Running initial schema setup..."
    	cat migrations/000001_create_users_table.up.sql | docker exec -i $(DOCKER_CONTAINER_NAME) psql -U $(DOCKER_DB_USER) -d $(DOCKER_DB_NAME)
    	@echo "Initial schema setup complete."

    docker-stop-db:
    	@echo "Stopping PostgreSQL Docker container..."
    	docker stop $(DOCKER_CONTAINER_NAME) || true
    	docker rm $(DOCKER_CONTAINER_NAME) || true

    docker-reset-db: docker-stop-db
    	@echo "Pruning unused Docker volumes (this might remove volumes from other projects if not careful!)..."
    	docker volume prune -f || true
    	@echo "Database has been completely reset."
    	make docker-run-db # Start a fresh database container after pruning

    # Database Migrations (migrate CLI - mostly for future, now using direct psql)
    # These targets might not be used directly for users table anymore, but are kept for structure.
    migrate-install:
    	@echo "Installing migrate CLI tool..."
    	@if ! command -v $(GO_MIGRATE) &> /dev/null; then \
    		echo "$(GO_MIGRATE) not found. Installing..."; \
    		curl -L [https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz](https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz) | tar xvz; \
    		sudo mv migrate /usr/local/bin/; \
    	fi

    # migrate-up: migrate-install
    # 	@echo "Running database migrations up..."
    # 	migrate -path migrations -database "postgresql://$(DOCKER_DB_USER):$(DOCKER_DB_PASSWORD)@localhost:$(shell echo $(DOCKER_DB_PORT) | cut -d':' -f1)/$(DOCKER_DB_NAME)?sslmode=disable" up

    # migrate-down: migrate-install
    # 	@echo "Running database migrations down..."
    # 	migrate -path migrations -database "postgresql://$(DOCKER_DB_USER):$(DOCKER_DB_PASSWORD)@localhost:$(shell echo $(DOCKER_DB_PORT) | cut -d':' -f1)/$(DOCKER_DB_NAME)?sslmode=disable" down 1

    # Help message
    help:
    	@echo "Usage:"
    	@echo "  make build          Builds the application binary"
    	@echo "  make run            Runs the application"
    	@echo "  make test           Runs all tests"
    	@echo "  make clean          Removes build artifacts"
    	@echo "  make lint           Runs linters"
    	@echo "  make docker-build   Builds the Docker image"
    	@echo "  make docker-run-db  Starts a local PostgreSQL database in Docker (stops existing one first and creates schema)"
    	@echo "  make docker-stop-db Stops and removes the local PostgreSQL database"
    	@echo "  make docker-reset-db Resets the database completely and recreates schema"
    	@echo "  make migrate-up     (CLI migrate tool) Applies database migrations"
    	@echo "  make migrate-down   (CLI migrate tool) Reverts the last database migration"
    	@echo "  make help           Displays this help message"
    EOF
    ```

9.  **Create `Dockerfile`:**
    ```bash
    cat << 'EOF' > Dockerfile
    FROM golang:1.22-alpine AS builder
    WORKDIR /app
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    RUN CGO_ENABLED=0 go build -o /main ./cmd/server
    FROM alpine:latest
    RUN apk --no-cache add ca-certificates
    WORKDIR /root/
    COPY --from=builder /main .
    EXPOSE 8080
    CMD ["./main"]
    EOF
    ```

10. **Create Migration File (`000001_create_users_table.up.sql`):**
    * Only the `.up.sql` file is needed as we'll run it directly.
    ```bash
    cat << 'EOF' > migrations/000001_create_users_table.up.sql
    -- migrations/000001_create_users_table.up.sql (SIMPLIFIED FOR TESTING - DIRECT EXECUTION)
    -- This version uses INTEGER ID and VARCHAR for tenant_id.
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        tenant_id VARCHAR(36) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        phone_number VARCHAR(50) NOT NULL,
        role VARCHAR(50) NOT NULL DEFAULT 'user',
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

    -- Create indexes
    CREATE INDEX idx_users_tenant_id ON users (tenant_id);
    CREATE UNIQUE INDEX idx_users_email_tenant_id ON users (email, tenant_id);
    EOF
    ```

11. **Create Utility Files (`internal/utils/errors/errors.go`, `internal/utils/log/log.go`, `internal/utils/response.go`):**
    * Open VS Code (`code .`) and create/update these files with the content I provided previously.

12. **Create Middleware Files (`internal/platform/http/middleware/auth_middleware.go`, `logging_middleware.go`, `recovery_middleware.go`):**
    * Open VS Code (`code .`) and create/update these files with the content I provided previously.

13. **Create Employee Module Files (`internal/app/employee/model.go`, `errors.go`, `service.go`, `handler.go`):**
    * Open VS Code (`code .`) and create/update these files with the content I provided previously.

14. **Update Main File (`cmd/server/main.go`):**
    * Open VS Code (`code .`) and update this file with the content I provided previously.

15. **Update OpenAPI File (`api/openapi.yaml`):**
    * Open VS Code (`code .`) and update this file with the content I provided previously.

### **Phase 3: Execution (Final Confirmation!)**

1.  **Close ALL WSL terminals.** (Crucial for a clean environment refresh).
2.  **Open Docker Desktop in Windows.** Ensure it is fully running.
3.  **Open a NEW WSL terminal.**
4.  **Navigate to your project directory:**
    ```bash
    cd /mnt/d/Playground/starterpack-golang-cleanarch
5.  **Perform a Full Database Reset (This is the crucial step that *should* create the `users` table):**
    ```bash
    make docker-reset-db
    ```
    * Observe the output carefully. You should see `CREATE TABLE`, `CREATE INDEX`, `CREATE INDEX`, and `Initial schema setup complete.`.

6.  **Verify `users` Table in Database (The most important confirmation):**
    ```bash
    docker exec -it starterpack-golang-cleanarch-db psql -U starteruser -d starterdb -c "\dt"
    ```
    * **The output MUST now show `public | users | table | starteruser`**. If it only shows `schema_migrations`, then there's still a fundamental issue with `psql` running the SQL (which is unlikely given previous successes, but worth checking).

7.  **Run the Go Application:**
    ```bash
    make run
* Observe the application logs.

8.  **Try POST /employees in Postman.**
    * This should now successfully create a new employee!

9. **Run Swagger Documentation.**
    ```bash
    docker run -p 8081:8080 -e SWAGGER_JSON=/app/openapi.yaml -v $(pwd)/api:/app swaggerapi/swagger-ui