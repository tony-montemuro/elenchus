package mocks

import "github.com/tony-montemuro/elenchus/internal/models"

type ProfileModel struct{}

func (m *ProfileModel) Insert(firstName, lastName, email, password string) error {
	return nil
}

func (m *ProfileModel) Authenticate(email, password string) (models.Profile, error) {
	return models.Profile{}, nil
}

func (m *ProfileModel) Exists(id int) (bool, error) {
	return false, nil
}
