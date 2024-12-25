package db

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb() (*gorm.DB, error) {
	port, err := strconv.Atoi(os.Getenv("pg_port"))
	if err != nil {
		return nil, errors.New("no valid port in environnement")
	}

	host := os.Getenv("pg_host")
	usr := os.Getenv("pg_username")
	pw := os.Getenv("pg_password")
	dbname := os.Getenv("pg_database_name")
	dsnf := "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
	dsn := fmt.Sprintf(dsnf, host, usr, pw, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error in connecting to database: %s", err.Error())
	}

	return db, nil
}
