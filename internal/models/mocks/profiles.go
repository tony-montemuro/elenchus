package mocks

type ProfileModel struct{}

func (m *ProfileModel) Insert(firstName, lastName, email, password string) error {
	return nil
}

func (m *ProfileModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *ProfileModel) Exists(id int) (bool, error) {
	return false, nil
}
