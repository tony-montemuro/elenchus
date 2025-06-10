package models

import (
	"database/sql"
	"time"
)

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

func (m *ProfileModel) Insert(firstName, lastName, email, password string) (int, error) {
	return 0, nil
}

func (m *ProfileModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *ProfileModel) Exists(id int) (bool, error) {
	return false, nil
}
