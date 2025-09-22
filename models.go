package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type User struct {
	Id         int
	Name       string
	Email      string
	Created_at time.Time
}

func (u *User) PrintUser() {
	fmt.Printf("%d: %s, %s - %s\n", u.Id, u.Name, u.Email, u.Created_at.Format("2006-01-02"))
}

func isDuplicateEmailError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
