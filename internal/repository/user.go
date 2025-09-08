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
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (usr domain.User,err error) {
	dataset :=u.db.From("users").Where(goqu.C("email").Eq(email))
	_, err = dataset.ScanStructContext(ctx, &usr)
	return
}