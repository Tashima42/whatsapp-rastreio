package data

import (
	"database/sql"
	"time"
)

type Group struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	DueDate            time.Time `json:"dueDate"`
	MonthlyTotalValue  float64   `json:"monthlyTotalValue"`
	MonthlyMemberValue float64   `json:"monthlyMemberValue"`
}

func (g *Group) CreateGroup(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO groups(name, due_date, monthly_total_value, monthly_member_value) VALUES($1, $2, $3, $4) RETURNING id;",
		g.Name,
		g.DueDate,
		g.MonthlyTotalValue,
		g.MonthlyMemberValue,
	).Scan(&g.ID)
}

func (g *Group) GetById(db *sql.DB) error {
	return db.QueryRow(
		"SELECT id, name, due_date, monthly_total_value, monthly_member_value FROM groups WHERE id=$1 LIMIT 1;",
		g.ID,
	).Scan(&g.ID, &g.Name, &g.DueDate, &g.MonthlyTotalValue, &g.MonthlyMemberValue)
}

func (g *Group) AddUserToGroup(db *sql.DB, userId int, role string) error {
	return db.QueryRow(
		"INSERT INTO group_user_account(group_id, user_account_id, role) VALUES($1, $2, $3);",
		g.ID,
		userId,
		role,
	).Err()
}
