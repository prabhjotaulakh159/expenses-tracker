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

	"github.com/gin-gonic/gin"
	"github.com/prabhjotaulakh159/expenses-tracker/db"
	"gorm.io/gorm"
)

func makeDatabaseConnection() (*gorm.DB, *sql.DB, error) {
	db, err := db.NewDb()
	if err != nil {
		return nil, nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	return db, sqlDb, err
}

func closeDatabaseConnection(sqlDb *sql.DB) {
	if err := sqlDb.Close(); err != nil {
		log.Fatalf("error in closing database connection: %s", err.Error())
	}
}

func getServer(router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}
}

func startServer(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error in starting server: %s", err)
	}
}

func doGracefulShutdown(srv *http.Server, shutDownContext context.Context) {
	log.Println("attempting to close server...")
	if err := srv.Shutdown(shutDownContext); err != nil {
		log.Fatalf("error in shutting down server: %s", err.Error())
	}
	// closeDatabaseConnection is defered here
	log.Println("attempting to close database connection...")
	log.Println("server and database connection closed successfully")
}

func main() {
	_, sqlDb, err := makeDatabaseConnection()
	if err != nil {
		log.Fatalf("error during database connection: %s", err.Error())
	}
	defer closeDatabaseConnection(sqlDb)
	log.Println("successfully connected to database")

	router := gin.Default()
	srv := getServer(router)
	go startServer(srv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutDownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doGracefulShutdown(srv, shutDownContext)
}
