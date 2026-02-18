package repository

import (
	"database/sql"
	"go-fiber-crud/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(username, passwordhash string) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err := r.DB.Exec(query, username, passwordhash)
	return err
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, password FROM users WHERE username = $1"
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
