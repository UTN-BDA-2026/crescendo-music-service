package services

import (
	"crescendo-api/mapping"
	"crescendo-api/models"
	"crescendo-api/repositories"
	"crescendo-api/security"
	"database/sql"
	"errors"
	"log"
	"net/mail"
	"time"
	"unicode"
)

type UserService interface {
	Register(mapping.UserRegisterDTO) (models.User, error)
	Login(mapping.UserLoginDTO) (string, error)
}

type userService struct {
	repository repositories.UserRepository
}

func NewUserService(repository repositories.UserRepository) UserService {
	service := userService{
		repository: repository,
	}

	return service
}

func isValidDateOfBirth(date time.Time) bool {
	now := time.Now().UTC()
	min := now.AddDate(-130, 0, 0)
	max := now.AddDate(-13, 0, 0)

	return !date.Before(min) && !date.After(max)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasLetter bool
	var hasNumber bool

	for _, c := range password {
		switch {
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsNumber(c):
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
}

func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	for _, c := range username {
		if !(unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_' || c == '.') {
			return false
		}
	}

	return true
}

func (s userService) Register(requestDTO mapping.UserRegisterDTO) (models.User, error) {

	if !isValidDateOfBirth(requestDTO.DateOfBirth) {
		return models.User{}, errors.New("invalid date of birth")
	}

	if !isValidEmail(requestDTO.Email) {
		return models.User{}, errors.New("invalid email")
	}

	if !isValidPassword(requestDTO.Password) {
		return models.User{}, errors.New("invalid password")
	}

	if !isValidUsername(requestDTO.Username) {
		return models.User{}, errors.New("invalid username")
	}

	passwordHash, err := security.HashPassword(requestDTO.Password)
	if err != nil {
		log.Printf("hashing failed: %v", err)
		return models.User{}, errors.New("something went wrong")
	}

	user := models.User{
		Username:     requestDTO.Username,
		Email:        requestDTO.Email,
		DateOfBirth:  requestDTO.DateOfBirth,
		PasswordHash: passwordHash,
		RegisterDate: time.Now().UTC(),
	}

	user.Id, err = s.repository.Create(user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s userService) Login(loginDTO mapping.UserLoginDTO) (string, error) {

	if loginDTO.Username == "" && loginDTO.Email == "" {
		return "", errors.New("invalid credentials")
	}

	if loginDTO.Username != "" && !isValidUsername(loginDTO.Username) {
		return "", errors.New("invalid credentials")
	}

	if loginDTO.Email != "" && !isValidEmail(loginDTO.Email) {
		return "", errors.New("invalid credentials")
	}

	if loginDTO.Password == "" {
		return "", errors.New("invalid credentials")
	}

	user, err := s.repository.GetByUsernameOrEmail(loginDTO.Username, loginDTO.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid credentials")
		}

		log.Printf("fetching user on login failed: %v", err)
		return "", errors.New("something went wrong")
	}

	token, err := security.GenerateLoginToken(user.Id, user.Username)

	if err != nil {
		log.Printf("token generating on login failed: %v", err)
		return "", errors.New("something went wrong")
	}

	return token, nil
}
