package repositories

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"fmt"
)

type UserRepository interface {
	Create(user models.User) (int, error)
	GetById(id int) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id int) error
}

type userRepository struct {
	db database.DBTX
}

func NewUserRepository(db database.DBTX) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r userRepository) Create(user models.User) (int, error) {
	var id int

	err := r.db.QueryRow(`
		INSERT INTO users (
			username,
			email,
			password_hash,
			register_date,
			date_of_birth,
			profile_image_url
		)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id
	`,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.RegisterDate,
		user.DateOfBirth,
		user.ProfilePictureUrl,
	).Scan(&id)

	return id, err
}

func (r userRepository) GetById(id int) (models.User, error) {
	row := r.db.QueryRow(`
			SELECT id, username, email, password_hash,
				register_date, date_of_birth, profile_image_url
			FROM users
			WHERE id = $1
		`, id)
	var user models.User
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.RegisterDate,
		&user.DateOfBirth,
		&user.ProfilePictureUrl,
	)
	user.DateOfBirth = user.DateOfBirth.UTC() //Correccion de conversión de fechas de Postgres a Go
	user.RegisterDate = user.RegisterDate.UTC()
	return user, err
}

func (r userRepository) Update(user models.User) (models.User, error) {
	var updated models.User

	err := r.db.QueryRow(`
		UPDATE users
		SET username = $1,
			email = $2,
			password_hash = $3,
			register_date = $4,
			date_of_birth = $5,
			profile_image_url = $6
		WHERE id = $7
		RETURNING id, username, email, password_hash,
				  register_date, date_of_birth, profile_image_url
	`,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.RegisterDate,
		user.DateOfBirth,
		user.ProfilePictureUrl,
		user.Id,
	).Scan(
		&updated.Id,
		&updated.Username,
		&updated.Email,
		&updated.PasswordHash,
		&updated.RegisterDate,
		&updated.DateOfBirth,
		&updated.ProfilePictureUrl,
	)
	updated.DateOfBirth = updated.DateOfBirth.UTC() //Correccion de conversión de fechas de Postgres a Go
	updated.RegisterDate = updated.RegisterDate.UTC()
	return updated, err
}

func (r userRepository) Delete(id int) error {
	result, err := r.db.Exec(`
		DELETE FROM users
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}

	return nil
}
