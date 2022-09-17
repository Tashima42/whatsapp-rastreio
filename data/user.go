package data

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Number    string `json:"number"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) CreateUser(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO user_accounts(number, created_at, updated_at) VALUES($1, $2, $3) RETURNING id;",
		u.Number,
		time.Now(),
		time.Now(),
	).Scan(&u.ID)
}

func (u *User) AddObject(db *sql.DB, object Object) error {
	return db.QueryRow(
		"INSERT INTO object_user_account(object_id, user_account_id, created_at, updated_at) VALUES($1, $2, $3, $4);",
		object.ID,
		u.ID,
		time.Now(),
		time.Now(),
	).Err()
}

func (u *User) GetById(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, number, created_at, updated_at FROM user_accounts WHERE id=$1 LIMIT 1;",
		u.ID,
	).Scan(&u.ID, &u.Number, &u.CreatedAt, &u.UpdatedAt)
}

func (u *User) GetByNumber(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, number, created_at, updated_at FROM user_accounts WHERE number=$1 LIMIT 1;",
		u.Number,
	).Scan(&u.ID, &u.Number, &u.CreatedAt, &u.UpdatedAt)
}

func (u *User) GetUserObjects(db *sql.DB) ([]Object, error) {
	var objects []Object
	rows, err := db.Query(
		`SELECT obj.id, obj.code, obj.name, obj.delivered, obj.last_sent_hash, obj.created_at, obj.updated_at
	FROM object_user_account as oua
	JOIN objects as obj
	ON oua.object_id = obj.id
	WHERE oua.user_account_id = $1;`,
		u.ID,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var object Object
		rows.Scan(&object.ID, &object.Code, &object.Name, &object.Delivered, &object.LastSentHash, &object.CreatedAt, &object.UpdatedAt)
		objects = append(objects, object)
	}
	return objects, nil
}

func GetUsers(db *sql.DB) ([]User, error) {
	var users []User
	rows, err := db.Query(
		"SELECT id, number, created_at, updated_at FROM user_accounts;",
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Number, &user.CreatedAt, &user.UpdatedAt)
		users = append(users, user)
	}
	return users, nil
}
