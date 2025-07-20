package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// func (m *ProfileModel) GetProfile(id int) (ProfilePublic, error) {

type ProfileModelInterface interface {
	Insert(string, string, string, string) error
	Authenticate(string, string) (Profile, error)
	Exists(int) (bool, error)
	GetProfile(int) (ProfilePublic, error)
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

type ProfilePublic struct {
	ID        int
	FirstName string
	LastName  string
	Deleted   *time.Time
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

func (m *ProfileModel) Authenticate(email, password string) (Profile, error) {
	var p Profile

	stmt := `SELECT id, first_name, last_name, email, hashed_password, created, updated, deleted 
	FROM profile WHERE email = ? AND deleted IS NULL`
	err := m.DB.QueryRow(stmt, email).Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email, &p.HashedPassword, &p.Created, &p.Updated, &p.Deleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Profile{}, ErrInvalidCredentials
		} else {
			return Profile{}, err
		}
	}

	err = bcrypt.CompareHashAndPassword(p.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return Profile{}, ErrInvalidCredentials
		} else {
			return Profile{}, err
		}
	}

	return p, nil
}

func (m *ProfileModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT id FROM profile WHERE id = ?)`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}

func (m *ProfileModel) GetProfile(id int) (ProfilePublic, error) {
	var p ProfilePublic
	stmt := `SELECT p.id, p.first_name, p.last_name
	FROM profile p
	WHERE p.id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.FirstName, &p.LastName)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return p, ErrInvalidCredentials
	}

	return p, err
}
