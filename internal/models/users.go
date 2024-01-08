package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	//use bcrypt to generate encrypted password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		//check if the error is something that violates the UNIQUE constraint
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" && strings.Contains(pqError.Message, "users_email_key") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exsits(id int) (bool, error) {
	return false, nil
}
