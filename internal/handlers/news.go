package handlers

import (
	"github.com/Smet1/bmstu-gis/internal/db"
	"github.com/Smet1/bmstu-gis/internal/logger"
	"github.com/jmoiron/sqlx"
	"net/http"
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
