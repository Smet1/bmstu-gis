package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Smet1/bmstu-gis/internal/db"
	"github.com/Smet1/bmstu-gis/internal/logger"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
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

			Response(res, http.StatusBadRequest, Error{Error: "can't create user"})
			return
		}

		Response(res, http.StatusCreated, u)
	}
}

func GetGetUserHandler(conn *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())

		userID := chi.URLParam(req, "userID") // from a route like /users/{userID}
		id, err := strconv.Atoi(userID)
		if err != nil {
			log.WithError(err).Error("can't convert url param to int")

			Response(res, http.StatusBadRequest, Error{Error: "invalid id"})
			return
		}

		u := &db.User{}
		err = u.Select(conn, int64(id))
		if err != nil {
			if errors.Cause(err) == sql.ErrNoRows {
				Response(res, http.StatusNotFound, Error{Error: "user not found"})
				return
			}

			log.WithError(err).Error("can't get user")

			Response(res, http.StatusBadRequest, Error{Error: "can't create user"})
			return
		}

		Response(res, http.StatusOK, u)
	}
}
