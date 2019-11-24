package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Point struct {
	ID      int64         `db:"id" json:"id,omitempty"`
	Name    string        `db:"name" json:"name,omitempty"`
	Cabinet bool          `db:"cabinet" json:"cabinet"`
	Stair   bool          `db:"stair" json:"stair"`
	X       int64         `db:"x" json:"x,omitempty"`
	Y       int64         `db:"y" json:"y,omitempty"`
	Level   int64         `db:"level" json:"level,omitempty"`
	NodeID  sql.NullInt64 `db:"node_id" json:"node_id,omitempty"`
}

func (p *Point) Insert(db *sqlx.DB) error {
	query := `
INSERT INTO map_points (name, cabinet, stair, x, y, level) 
VALUES (:name, :cabinet, :stair, :x, :y, :level)
RETURNING id
`
	row, err := db.NamedQuery(query, p)
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

	p.ID = res.Int64

	return nil
}

type ListPoints struct {
	Points []*Point `json:"points"`
}

func (l *ListPoints) Get(db *sqlx.DB) error {
	query := `
SELECT id, name, cabinet, stair, x, y, level, node_id
FROM map_points
`
	err := db.Select(&l.Points, query)
	if err != nil {
		return errors.Wrap(err, "can't do query")
	}

	return nil
}
