package data

import (
	"database/sql"
	"time"
)

type Object struct {
	ID           int    `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Delivered    bool   `json:"delivered"`
	LastSentHash string `json:"lastSentHash"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (o *Object) CreateObject(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO objects(code, name, created_at, updated_at) VALUES($1, $2, $3, $4) RETURNING id;",
		o.Code,
		o.Name,
		time.Now(),
		time.Now(),
	).Scan(&o.ID)
}

func (o *Object) UpdateLastSentHash(db *sql.DB) error {
	o.UpdatedAt = time.Now()
	return db.QueryRow(
		"UPDATE objects SET last_sent_hash = $1, updated_at = $2 WHERE id = $3;",
		o.LastSentHash,
		o.UpdatedAt,
		o.ID,
	).Err()
}

func (o *Object) GetByCode(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, code, name, created_at, updated_at FROM objects WHERE code = $1",
		o.Code,
	).Scan(&o.ID, &o.Code, &o.Name, &o.CreatedAt, &o.UpdatedAt)
}

func (o *Object) GetObjectEvent(db *sql.DB) (Event, error) {
	var event Event
	err := db.QueryRow(
		"SELECT id, hash, body, created_at, updated_at FROM events WHERE object_id = $1 LIMIT 1;",
		o.ID,
	).Scan(&event.ID, &event.Hash, &event.Body, &event.CreatedAt, &event.UpdatedAt)
	return event, err
}

func GetObjectsByLastUpdated(db *sql.DB, lastUpdated time.Time) ([]Event, []Object, error) {
	var events []Event
	var objects []Object
	rows, err := db.Query(
		"SELECT obj.id, obj.code, ev.id, ev.hash FROM objects as obj JOIN events as ev ON ev.object_id = obj.id WHERE obj.delivered = false AND ev.updated_at < $1",
		lastUpdated,
	)
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		var event Event
		var object Object
		err = rows.Scan(&object.ID, &object.Code, &event.ID, &event.Hash)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, event)
		objects = append(objects, object)
	}
	return events, objects, nil
}
