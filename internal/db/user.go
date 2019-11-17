package db

import (
	"database/sql"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type User struct {
	ID         int64     `db:"id"`
	Login      string    `db:"login"`
	Password   string    `db:"password"`
	Registered time.Time `db:"registered"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(&u,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&u.Login, validation.Required, validation.Length(5, 50)),
	)
}

func (u *User) Insert(db *sqlx.DB) error {
	u.Registered = time.Now()
	query := `
INSERT INTO users (login, password, avatar, karma, registered) 
VALUES (:login, :password, :avatar, :karma, :registered)
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

func (u *User) Select() {

}

func (u *User) Update() {

}

func (u *User) Delete() {

}
