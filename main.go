package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	
	db, err := connectionToDB(ctx)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to DB")
}

func connectionToDB(ctx context.Context)(*pgxpool.Pool, error){
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), 
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("Parse config:", err)
	}

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("Create connection pool:", err)
	}

	for i := 0; i < 30; i++ {
		if err := db.Ping(ctx); err == nil {
			return db, nil
		}
		time.Sleep(1* time.Second)
	}

	return nil, fmt.Errorf("Database not availible")
}