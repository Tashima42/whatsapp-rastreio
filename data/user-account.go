package data

import (
	"database/sql"
)

type UserAccount struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	City     string `json:"city"`
	PixKey   string `json:"pixKey"`
}

func (u *UserAccount) CreateUserAccount(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO user_accounts(username, email, name, city, pix_key) VALUES($1, $2, $3, $4, $5) RETURNING id;",
		u.Username,
		u.Email,
		u.Name,
		u.City,
		u.PixKey,
	).Scan(&u.ID)
}

func (u *UserAccount) GetById(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, username, email, name, city, pix_key FROM user_accounts WHERE id=$1 LIMIT 1;",
		u.ID,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Name, &u.City, &u.PixKey)
}

func (u *UserAccount) GetByUsername(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, username, email, name, city, pix_key, role FROM user_accounts WHERE username=$1 LIMIT 1;",
		u.Username,
	).Scan(&u.ID, &u.Username, &u.Name, &u.City, u.PixKey)
}
