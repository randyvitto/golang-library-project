package repository

import (
	"belajar-golang-rest-api/lat/domain"
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

type UserRepository struct {
	db *goqu.Database
}


func NewUser(con *sql.DB) domain.UserRepository {
	return &UserRepository{
		db: goqu.New("default", con),
	}
}

// FindByEmail implements domain.UserRepository.
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (usr domain.User, err error) {
	dataset := u.db.From("users").Where(goqu.C("email").Eq(email))
	_, err = dataset.ScanStructContext(ctx, &usr)
	return
}

// Delete implements domain.UserRepository.
func (u *UserRepository) Delete(ctx context.Context, id string) error {
	executor := u.db.Delete("users").Where(goqu.C("id").Eq(id)).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// Save implements domain.UserRepository.
func (u *UserRepository) Save(ctx context.Context, c *domain.User) error {
	executor := u.db.Insert("users").Rows(c).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// Update implements domain.UserRepository.
func (u *UserRepository) Update(ctx context.Context, c *domain.User) error {
	executor := u.db.Update("users").Where(goqu.C("id").Eq(c.Id)).Set(c).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}