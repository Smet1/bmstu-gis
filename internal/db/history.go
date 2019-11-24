package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type History struct {
	ID        int64  `db:"id" json:"id,omitempty"`
	UserID    int64  `db:"user_id" json:"user_id,omitempty"`
	PointFrom string `db:"point_from" json:"point_from,omitempty"`
	PointTo   string `db:"point_to" json:"point_to,omitempty"`
}

func (h *History) Insert(db *sqlx.DB) error {
	query := `
INSERT INTO users_history (user_id, point_from, point_to) 
VALUES (:user_id, :point_from, :point_to)
RETURNING id
`
	row, err := db.NamedQuery(query, h)
	if err != nil {
		return errors.Wrap(err, "can't do query")
	}
	defer row.Close()

	res := &sql.NullInt64{}
	for row.Next() {
		err = row.Scan(res)
		if err != nil {
			return errors.Wrap(err, "can't get id")
		}
	}

	h.ID = res.Int64

	return nil
}

func ClearHistory(db *sqlx.DB, userID int64) error {
	query := `
DELETE FROM users_history
WHERE user_id = $1
`
	_, err := db.Exec(query, userID)
	if err != nil {
		return errors.Wrap(err, "can't do query")
	}

	return nil
}

type ListHistory struct {
	History []*History `json:"history"`
}

func (l *ListHistory) Get(db *sqlx.DB, userID int64) error {
	query := `
SELECT user_id, point_from, point_to
FROM users_history
WHERE user_id = $1
`
	err := db.Select(&l.History, query, userID)
	if err != nil {
		return errors.Wrap(err, "can't do query")
	}

	return nil
}
