package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Smet1/bmstu-gis/internal/db"
	"github.com/Smet1/bmstu-gis/internal/logger"
	"github.com/jmoiron/sqlx"
)

func GetGetNewsHandler(conn *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())

		l := &db.ListNews{}
		err := l.Get(conn)
		if err != nil {
			log.WithError(err).Error("can't get news")

			Response(res, http.StatusBadRequest, Error{Error: "can't get news"})
			return
		}

		Response(res, http.StatusOK, l)
	}
}

func GetCreateNewsHandler(conn *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())

		n := &db.News{}
		body, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		_ = json.Unmarshal(body, n)

		err := n.Insert(conn)
		if err != nil {
			log.WithError(err).Error("can't create news")

			Response(res, http.StatusBadRequest, Error{Error: "can't create news"})
			return
		}

		Response(res, http.StatusOK, n)
	}
}
