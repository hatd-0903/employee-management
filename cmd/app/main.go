package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"employee-management/internal/config"
	"employee-management/internal/handlers"
	"employee-management/internal/middleware"
	"employee-management/internal/repositories"
	"employee-management/internal/services"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	db, err := connectWithRetry(cfg.DSN(), 10, 2*time.Second)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	employeeRepo := repositories.NewMySQLEmployeeRepository(db)
	departmentRepo := repositories.NewMySQLDepartmentRepository(db)

	employeeService := services.NewEmployeeService(employeeRepo, departmentRepo)
	departmentService := services.NewDepartmentService(departmentRepo)
	exportService := services.NewExportService(employeeRepo, cfg.ExportDir)

	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	departmentHandler := handlers.NewDepartmentHandler(departmentService, employeeService)
	exportHandler := handlers.NewExportHandler(exportService)

	router := handlers.NewRouter(employeeHandler, departmentHandler, exportHandler)

	handler := middleware.Chain(router,
		middleware.Recover,
		middleware.Logging,
		middleware.BasicAuth(cfg.AuthUsers),
	)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("employee-management listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

// connectWithRetry waits for MySQL to accept connections, which matters when
// the app and database start together under docker-compose.
func connectWithRetry(dsn string, attempts int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 1; i <= attempts; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			err = db.PingContext(ctx)
			cancel()
			if err == nil {
				return db, nil
			}
			db.Close()
		}
		log.Printf("waiting for database (attempt %d/%d): %v", i, attempts, err)
		time.Sleep(delay)
	}
	return nil, err
}
