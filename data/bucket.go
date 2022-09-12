package data

import (
	"database/sql"
	"time"
)

type Bucket struct {
	ID                  int       `json:"id"`
	Name                string    `json:"name"`
	DueDate             time.Time `json:"dueDate"`
	MonthlyTotalValue   float64   `json:"monthlyTotalValue"`
	MonthlyMemberlValue float64   `json:"monthlyMemberValue"`
}

func (b *Bucket) CreateBucket(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO buckets(name, due_date, monthly_total_value, monthly_member_value) VALUES($1, $2, $3, $4) RETURNING id;",
		b.Name,
		b.DueDate,
		b.MonthlyTotalValue,
		b.MonthlyMemberlValue,
	).Scan(&b.ID)
}
