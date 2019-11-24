package db

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type News struct {
	ID      int64     `db:"id" json:"id,omitempty"`
	Title   string    `db:"title" json:"title,omitempty"`
	Payload string    `db:"payload" json:"payload,omitempty"`
	Created time.Time `db:"created" json:"created,omitempty"`
}

func (n *News) Insert(db *sqlx.DB) error {
	n.Created = time.Now()
	query := `
INSERT INTO news (title, payload, created) 
VALUES (:title, :payload, :created)
RETURNING id
`
	row, err := db.NamedQuery(query, n)
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

	n.ID = res.Int64

	return nil
}

type ListNews struct {
	News []*News `json:"news"`
}

func (l *ListNews) Get(db *sqlx.DB) error {
	query := `
SELECT title, payload, created
FROM news
`
	err := db.Select(&l.News, query)
	if err != nil {
		return errors.Wrap(err, "can't do query")
	}

	return nil
}
