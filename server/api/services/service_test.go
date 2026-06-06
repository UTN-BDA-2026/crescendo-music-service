package services_test

import "crescendo-api/models"

type mockUserRepository struct {
	createFunc           func(user models.User) (int, error)
	getByIdFunc          func(id int) (models.User, error)
	updateFunc           func(user models.User) (models.User, error)
	deleteFunc           func(id int) error
	getByUsernameOrEmail func(username string, email string) (models.User, error)
}

func (m mockUserRepository) Create(user models.User) (int, error) {
	return m.createFunc(user)
}

func (m mockUserRepository) GetById(id int) (models.User, error) {
	return m.getByIdFunc(id)
}

func (m mockUserRepository) Update(user models.User) (models.User, error) {
	return m.updateFunc(user)
}

func (m mockUserRepository) Delete(id int) error {
	return m.deleteFunc(id)
}

func (m mockUserRepository) GetByUsernameOrEmail(username string, email string) (models.User, error) {
	return m.getByUsernameOrEmail(username, email)
}
