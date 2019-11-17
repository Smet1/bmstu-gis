package handlers

import (
	"encoding/json"
	"github.com/Smet1/bmstu-gis/internal/db"
	"github.com/Smet1/bmstu-gis/internal/logger"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
)

func GetCreateUserHandler(conn *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())
		u := &db.User{}
		body, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		_ = json.Unmarshal(body, u)

		err := u.Insert(conn)
		if err != nil {
			log.WithError(err).Error("can't create user")
		}
	}
}
