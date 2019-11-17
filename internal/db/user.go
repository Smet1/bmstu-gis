package db

import (
	"database/sql"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type User struct {
	ID         int64     `db:"id" json:"id,omitempty"`
	Login      string    `db:"login" json:"login,omitempty"`
	Password   string    `db:"password" json:"password,omitempty"`
	Registered time.Time `db:"registered" json:"registered,omitempty"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(&u,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&u.Login, validation.Required, validation.Length(5, 50)),

		validation.Field(&u.Password, validation.Required, validation.Length(3, 20)),
	)
}

func (u *User) Insert(db *sqlx.DB) error {
	u.Registered = time.Now()
	query := `
INSERT INTO users (login, password, registered) 
VALUES (:login, :password, :registered)
RETURNING id
`
	row, err := db.NamedQuery(query, u)
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

	u.ID = res.Int64

	return nil
}

func (u *User) GetByID(db *sqlx.DB, id int64) error {
	query := `
SELECT id, login, registered 
FROM users 
WHERE id = $1
`
	err := db.Get(u, query, id)
	if err != nil {
		return errors.Wrap(err, "can't do query")
	}

	return nil
}

func (u *User) GetByLogin(db *sqlx.DB, login string) error {
	query := `
SELECT id, login, password, registered 
FROM users 
WHERE login = $1
`
	err := db.Get(u, query, login)
	if err != nil {
		return errors.Wrap(err, "can't do query")
	}

	return nil
}

func (u *User) Update() {

}

func (u *User) Delete() {

}
