package data

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID         int       `json:"id"`
	Attachment string    `json:"attachment"`
	GroupId    int       `json:"groupId"`
	Date       time.Time `json:"date"`
	UserId     int       `json:"userId"`
	Paid       bool      `json:"paid"`
}

func (t *Transaction) CreateTransaction(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO buckets(attachment, group_id, date, user_account_id) VALUES($1, $2, $3, $4) RETURNING id;",
		t.Attachment,
		t.GroupId,
		t.Date,
		t.UserId,
	).Scan(&t.ID)
}
