package main

import (
	"fmt"
	"time"
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
