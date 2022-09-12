package data

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID            int       `json:"id"`
	Attachment    string    `json:"attackment"`
	BucketId      int       `json:"bucketId"`
	Date          time.Time `json:"date"`
	UserAccountId int       `json:"userAccountId"`
}

func (t *Transaction) CreateTransaction(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO buckets(attachment, bucket_id, date, user_account_id) VALUES($1, $2, $3, $4) RETURNING id;",
		t.Attachment,
		t.BucketId,
		t.Date,
		t.UserAccountId,
	).Scan(&t.ID)
}
