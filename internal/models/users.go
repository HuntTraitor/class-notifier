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

// authenticate the user and make sure credentials match using bcrypt
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var userid int
	var hashedPassword []byte

	stmt := `SELECT userid, hashed_password FROM users WHERE email = $1`

	err := m.DB.QueryRow(stmt, email).Scan(&userid, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return userid, nil
}

// probably a better way to do this - will implement later
func (m *UserModel) GetEmail(id int) (string, error) {
	var email string
	stmt := `SELECT email FROM users WHERE userid = $1`

	err := m.DB.QueryRow(stmt, id).Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoRecord
		} else {
			return "", err
		}
	}
	return email, nil
}

func (m *UserModel) Exsits(userid int) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT true FROM users WHERE userid = $1)`

	err := m.DB.QueryRow(stmt, userid).Scan(&exists)
	return exists, err
}
