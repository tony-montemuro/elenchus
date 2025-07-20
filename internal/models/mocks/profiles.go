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

func (m *ProfileModel) GetProfileNames(id int) (string, string, error) {
	return "", "", nil
}

func (m *ProfileModel) UpdateProfileNames(firstName, lastName string, profileID int) error {
	return nil
}
