package domain

import "context"

type User struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type UserRepository interface {
	Save(ctx context.Context, c *User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, c *User) error
	Delete(ctx context.Context, id string) error
}