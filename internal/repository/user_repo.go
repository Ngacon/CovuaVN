package repository

import (
	"database/sql"
	"github.com/Ngacon/covuavn/internal/domain"
)

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) Create(u domain.User) error {
	_, err := r.DB.Exec("INSERT INTO users VALUES ($1, $2)", u.Username, u.Password)
	return err
}

func (r *UserRepo) Get(username string) (string, error) {
	var pass string
	err := r.DB.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&pass)
	return pass, err
}