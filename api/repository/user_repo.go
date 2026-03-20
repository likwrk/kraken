package repository

import (
	"database/sql"
	"kraken/api/models"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetAll() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, name, age FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *userRepo) Create(user models.User) error {
	_, err := r.db.Exec("INSERT INTO users(name, age) VALUES(?, ?)", user.Name, user.Age)
	return err
}
