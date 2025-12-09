package usecase

import (
	domain "user-service/internal/domen"
	"user-service/pkg"
)

type UserUsecase interface {
	GetUserByTelegramID(telegramID int64) (*domain.User, error)
	CreateUser(input *domain.User) (*domain.User, error)
}


type userUsecase struct {
	repo domain.UserRepository
}


func NewUserUsecase(repo domain.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}


func (u *userUsecase) GetUserByTelegramID(telegramID int64) (*domain.User, error) {
	// Read-Through pattern
	return pkg.ReadThroughUser(telegramID, func() (*domain.User, error) {
		// loader → repository chaqiradi
		return u.repo.GetByTelegramID(telegramID)
	})
}


func (u *userUsecase) CreateUser(input *domain.User) (*domain.User, error) {
	// 1️⃣ Validation
	if input.TelegramID == 0 || input.FirstName == "" || input.Role == "" {
		return nil, ErrInvalidUserInput
	}

	// 2️⃣ Write-Through pattern
	return pkg.WriteThroughUser(input, func(user *domain.User) (*domain.User, error) {
		// loader → repository Create chaqiriladi
		return u.repo.Create(user)
	})
}

// Error definition
var ErrInvalidUserInput = &UserError{Message: "invalid user input"}

type UserError struct {
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}

