package domain

type UserRepository interface {
	GetByTelegramID(telegramID int64) (*User, error)
	Create(user *User) (*User, error)
}
