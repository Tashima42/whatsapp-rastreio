package data

import (
	"database/sql"
	"time"
)

type Event struct {
	ID        int    `json:"id"`
	Hash      string `json:"hash"`
	Body      string `json:"body"`
	ObjectId  int    `json:"objectId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *Event) CreateEvent(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO events(hash, body,object_id, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id;",
		e.Hash,
		e.Body,
		e.ObjectId,
		time.Now(),
		time.Now(),
	).Scan(&e.ID)
}
func (e *Event) Update(db *sql.DB) error {
	e.UpdatedAt = time.Now()
	return db.QueryRow(
		"UPDATE events SET hash = $1, body = $2, updated_at = $3 WHERE id = $4;",
		e.Hash,
		e.Body,
		e.UpdatedAt,
		e.ID,
	).Err()
}
