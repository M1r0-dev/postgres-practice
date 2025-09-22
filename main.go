package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	//db connection
	ctx := context.Background()

	db, err := connectionToDB(ctx)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to DB")

	//example requests
	displayData(ctx, db)
	fmt.Println()

	user, err := getByUserId(ctx, db, 1) // just example arg
	if err != nil {
		log.Fatal("Failed to find user:", err)
	}
	user.PrintUser()

}

func connectionToDB(ctx context.Context) (*pgxpool.Pool, error) {
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
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("Database not availible")
}

func displayData(ctx context.Context, db *pgxpool.Pool) error {
	fmt.Println("===Current Data in DB===")
	rows, err := db.Query(ctx, "SELECT id, name, email, created_at FROM users ORDER BY id")
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("Users:")
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Created_at); err != nil {
			return err
		}
		u.PrintUser()
	}
	return nil
}

func getByUserId(ctx context.Context, db *pgxpool.Pool, id int) (*User, error) {
	var u User

	err := db.QueryRow(ctx,
		"SELECT id, name, email, created_at FROM users WHERE id = $1",
		id).
		Scan(&u.Id, &u.Name, &u.Email, &u.Created_at)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}
