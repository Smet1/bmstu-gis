package db

import (
	"net/url"

	"github.com/Smet1/bmstu-gis/internal/config"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func EnsureDBConn(cfg *config.Config) (*sqlx.DB, error) {
	v := url.Values{}
	v.Add("sslmode", cfg.DB.SSLMode)

	p := url.URL{
		Scheme:     cfg.DB.Database,
		Opaque:     "",
		User:       url.UserPassword(cfg.DB.Username, cfg.DB.Password),
		Host:       cfg.DB.Host,
		Path:       cfg.DB.Name,
		RawPath:    "",
		ForceQuery: false,
		RawQuery:   v.Encode(),
		Fragment:   "",
	}

	connectURL, err := pq.ParseURL(p.String())
	if err != nil {
		return nil, errors.Wrap(err, "can't create url for db connection")
	}

	instance, err := sqlx.Connect(cfg.DB.Database, connectURL)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect db")
	}

	return instance, nil
}
