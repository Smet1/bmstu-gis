package db

import (
	"github.com/Smet1/bmstu-gis/internal/config"
	"net/url"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func EnsureDBConn(config *config.Config) (*sqlx.DB, error) {
	v := url.Values{}
	v.Add("sslmode", config.DB.SSLMode)

	p := url.URL{
		Scheme:     config.DB.Database,
		Opaque:     "",
		User:       url.UserPassword(config.DB.Username, config.DB.Password),
		Host:       config.DB.Host,
		Path:       config.DB.Name,
		RawPath:    "",
		ForceQuery: false,
		RawQuery:   v.Encode(),
		Fragment:   "",
	}

	connectURL, err := pq.ParseURL(p.String())
	if err != nil {
		return nil, errors.Wrap(err, "can't create url for db connection")
	}

	instance, err := sqlx.Connect(config.DB.Database, connectURL)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect db")
	}

	return instance, nil
}
