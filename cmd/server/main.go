package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"starterpack-golang-cleanarch/internal/app/employee"
	"starterpack-golang-cleanarch/internal/platform/http/middleware"
	"starterpack-golang-cleanarch/internal/repository"
	"starterpack-golang-cleanarch/internal/utils"
	"starterpack-golang-cleanarch/internal/utils/log"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Initialize Logger
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	log.InitLogger(env)
	defer log.Sync()
	log.Info(context.Background(), fmt.Sprintf("Starting %s in %s environment...", os.Getenv("APP_NAME"), env))

	// 2. Database Connection
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "starteruser"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "supersecret_db_password"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "starterdb"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf(context.Background(), "Failed to open database connection: %v", err)
	}

	db := sqlx.NewDb(sqlDB, "postgres")

	if err = db.Ping(); err != nil {
		log.Fatalf(context.Background(), "Failed to ping database: %v", err)
	}
	log.Info(context.Background(), "Successfully connected to database.")

	defer func() {
		log.Info(context.Background(), "Closing database connection...")
		if err := db.Close(); err != nil {
			log.Errorf(context.Background(), "Error closing database connection: %v", err)
		}
	}()

	// 3. Initialize Validator
	appValidator := validator.New()

	// 4. Setup Main Router (using Gorilla Mux)
	r := mux.NewRouter()

	// 5. Register General Endpoints (NO AUTHENTICATION)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			http.Error(w, "Database not reachable", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}).Methods("GET")

	r.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		type ServerInfo struct {
			AppName     string `json:"appName"`
			AppVersion  string `json:"appVersion"`
			GoVersion   string `json:"goVersion"`
			Environment string `json:"environment"`
			CurrentTime string `json:"currentTime"`
		}
		info := ServerInfo{
			AppName:     os.Getenv("APP_NAME"),
			AppVersion:  "1.0.0",
			GoVersion:   "1.22.x",
			Environment: env,
			CurrentTime: time.Now().Format(time.RFC3339),
		}
		utils.RespondJSON(w, http.StatusOK, info)
	}).Methods("GET")

	// 6. Create a Sub-Router for Authenticated Routes
	authenticatedRouter := r.PathPrefix("/").Subrouter()
	authenticatedRouter.Use(middleware.RecoveryMiddleware)
	authenticatedRouter.Use(middleware.LoggingMiddleware)
	authenticatedRouter.Use(middleware.AuthMiddleware)

	// --- Dependency Injection (DI) & Feature Module Registration ---
	// Employee Module Wiring - ACTUAL IMPLEMENTATION FOR DEMO
	employeeRepo := repository.NewPostgreSQLEmployeeRepository(db)
	employeeService := employee.NewEmployeeService(employeeRepo)
	employeeHandler := employee.NewEmployeeHandler(employeeService, appValidator)
	employeeHandler.RegisterRoutes(authenticatedRouter)

	// Other modules (e.g., Auth, Reporting) will be wired here too.

	// 7. Start HTTP Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Infof(context.Background(), "Server listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(context.Background(), "HTTP server ListenAndServe: %v", err)
		}
	}()

	// 8. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info(context.Background(), "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf(context.Background(), "Server forced to shutdown: %v", err)
	}

	log.Info(context.Background(), "Server exited gracefully.")
}
