package models

import (
	"errors"

	"example.com/event-booking/db"
	"example.com/event-booking/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users (email, password) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPass, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPass)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var password string
	err := row.Scan(&u.ID, &password)

	if err != nil {
		return errors.New("Invalid email or password")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, password)

	if !passwordIsValid {
		return errors.New("Invalid email or password")
	}

	return nil
}
