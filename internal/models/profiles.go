package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type ProfileModelInterface interface {
	Insert(firstName, lastName, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

type Profile struct {
	ID             int
	FirstName      string
	LastName       string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Updated        time.Time
	Deleted        *time.Time
}

type ProfileModel struct {
	DB *sql.DB
}

func (m *ProfileModel) Insert(firstName, lastName, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO profile (first_name, last_name, email, hashed_password, created, updated)
	VALUES (?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, firstName, lastName, email, hashedPassword)
	if err != nil {
		// check if error is specifically because email already in use
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "profile_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *ProfileModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *ProfileModel) Exists(id int) (bool, error) {
	return false, nil
}
