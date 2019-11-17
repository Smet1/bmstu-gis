package handlers

import (
	"encoding/json"
	"github.com/Smet1/bmstu-gis/internal/db"
	"github.com/Smet1/bmstu-gis/internal/logger"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetCreateHistoryHandler(conn *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())
		h := &db.History{}
		body, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		_ = json.Unmarshal(body, h)

		err := h.Insert(conn)
		if err != nil {
			log.WithError(err).Error("can't create history")

			Response(res, http.StatusBadRequest, Error{Error: "can't create history"})
			return
		}

		Response(res, http.StatusCreated, h)
	}
}

func GetGetHistoryHandler(conn *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())

		userID := chi.URLParam(req, "userID") // from a route like /users/{userID}
		id, err := strconv.Atoi(userID)
		if err != nil {
			log.WithError(err).Error("can't convert url param to int")

			Response(res, http.StatusBadRequest, Error{Error: "invalid id"})
			return
		}

		l := &db.ListHistory{}
		err = l.Get(conn, int64(id))
		if err != nil {
			log.WithError(err).Error("can't get history")

			Response(res, http.StatusBadRequest, Error{Error: "can't get history"})
			return
		}

		Response(res, http.StatusOK, l)
	}
}

func GetClearHistoryHandler(conn *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())

		userID := chi.URLParam(req, "userID") // from a route like /users/{userID}
		id, err := strconv.Atoi(userID)
		if err != nil {
			log.WithError(err).Error("can't convert url param to int")

			Response(res, http.StatusBadRequest, Error{Error: "invalid id"})
			return
		}

		err = db.ClearHistory(conn, int64(id))
		if err != nil {
			log.WithError(err).Error("can't clear history")

			Response(res, http.StatusBadRequest, Error{Error: "can't clear history"})
			return
		}

		Response(res, http.StatusOK, nil)
	}
}
