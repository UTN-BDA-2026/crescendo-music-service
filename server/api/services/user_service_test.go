package services_test

import (
	"crescendo-api/mapping"
	"crescendo-api/models"
	"crescendo-api/security"
	"crescendo-api/services"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCanRegisterUser_Success(t *testing.T) {
	check := require.New(t)

	userRegisterDTO := mapping.UserRegisterDTO{
		Username:    "User",
		Email:       "user@mail.com",
		Password:    "password77",
		DateOfBirth: time.Date(1999, 8, 15, 0, 0, 0, 0, time.UTC).UTC(),
	}
	repo := mockUserRepository{
		createFunc: func(user models.User) (int, error) {
			return 1, nil
		},
	}
	userService := services.NewUserService(repo)

	before := time.Now().UTC()

	user, err := userService.Register(userRegisterDTO)

	after := time.Now().UTC()

	check.NoError(err)
	check.Equal(1, user.Id)
	check.Equal(userRegisterDTO.Username, user.Username)
	check.Equal(userRegisterDTO.Email, user.Email)
	check.NotEmpty(user.PasswordHash)
	check.NotEqual(userRegisterDTO.Password, user.PasswordHash)
	check.WithinDuration(before, user.RegisterDate, after.Sub(before))
	check.Equal(userRegisterDTO.DateOfBirth, user.DateOfBirth)
}

func TestCanRegisterUser_InvalidPassword(t *testing.T) {
	check := require.New(t)

	repo := mockUserRepository{
		createFunc: func(user models.User) (int, error) {
			t.Fatal("repository should not be called")
			return 0, nil
		},
	}

	userService := services.NewUserService(repo)

	user, err := userService.Register(mapping.UserRegisterDTO{
		Username:    "Username",
		Email:       "user@mail.com",
		Password:    "123",
		DateOfBirth: time.Date(1999, 8, 15, 0, 0, 0, 0, time.UTC),
	})

	check.Error(err)
	check.Equal(models.User{}, user)
}

func TestCanRegisterUser_InvalidEmail(t *testing.T) {
	check := require.New(t)

	repo := mockUserRepository{
		createFunc: func(user models.User) (int, error) {
			t.Fatal("repository should not be called")
			return 0, nil
		},
	}

	userService := services.NewUserService(repo)

	user, err := userService.Register(mapping.UserRegisterDTO{
		Username:    "Username",
		Email:       "invalid-email",
		Password:    "password77",
		DateOfBirth: time.Date(1999, 8, 15, 0, 0, 0, 0, time.UTC),
	})

	check.Error(err)
	check.Equal(models.User{}, user)
}

func TestCanRegisterUser_InvalidUsername(t *testing.T) {
	check := require.New(t)

	repo := mockUserRepository{
		createFunc: func(user models.User) (int, error) {
			t.Fatal("repository should not be called")
			return 0, nil
		},
	}

	userService := services.NewUserService(repo)

	user, err := userService.Register(mapping.UserRegisterDTO{
		Username:    "Invalid Username",
		Email:       "user@mail.com",
		Password:    "password77",
		DateOfBirth: time.Date(1999, 8, 15, 0, 0, 0, 0, time.UTC),
	})

	check.Error(err)
	check.Equal(models.User{}, user)
}

func TestCanUserLogin_ValidUsernameAndPassword(t *testing.T) {
	check := require.New(t)

	dto := mapping.UserLoginDTO{
		Username: "Username",
		Password: "Password77",
	}
	repo := mockUserRepository{
		getByUsernameOrEmail: func(username, email string) (models.User, error) {
			passHash, err := security.HashPassword(dto.Password)
			check.NoError(err)
			return models.User{
				Id:           5,
				Username:     dto.Username,
				PasswordHash: passHash,
			}, nil
		},
	}

	service := services.NewUserService(repo)

	token, err := service.Login(dto)
	check.NoError(err)
	check.NotEmpty(token)
}

func TestCanUserLogin_ValidEmailAndPassword(t *testing.T) {
	check := require.New(t)

	dto := mapping.UserLoginDTO{
		Email:    "test@mail.com",
		Password: "Password77",
	}
	repo := mockUserRepository{
		getByUsernameOrEmail: func(username, email string) (models.User, error) {
			passHash, err := security.HashPassword(dto.Password)
			check.NoError(err)
			return models.User{
				Id:           5,
				Email:        dto.Email,
				PasswordHash: passHash,
			}, nil
		},
	}

	service := services.NewUserService(repo)

	token, err := service.Login(dto)
	check.NoError(err)
	check.NotEmpty(token)
}

func TestCanUserLogin_IncorrectPassword(t *testing.T) {
	check := require.New(t)

	dto := mapping.UserLoginDTO{
		Username: "Username",
		Password: "Password77",
	}
	repo := mockUserRepository{
		getByUsernameOrEmail: func(username, email string) (models.User, error) {
			passHash, err := security.HashPassword("InvalidPassword6")
			check.NoError(err)
			return models.User{
				Id:           5,
				Username:     dto.Username,
				PasswordHash: passHash,
			}, nil
		},
	}

	service := services.NewUserService(repo)

	token, err := service.Login(dto)
	check.Error(err)
	check.Empty(token)
}
