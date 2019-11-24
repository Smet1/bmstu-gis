package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Smet1/bmstu-gis/internal/db"
	"github.com/Smet1/bmstu-gis/internal/logger"
	"github.com/jmoiron/sqlx"
)

func GetGetPointsHandler(conn *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())

		listPoints := &db.ListPoints{}
		err := listPoints.Get(conn)
		if err != nil {
			log.WithError(err).Error("can't get points")

			Response(res, http.StatusBadRequest, Error{Error: "can't get points"})
			return
		}

		Response(res, http.StatusOK, listPoints)
	}
}

func GetCreatePointHandler(conn *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())

		point := &db.Point{}
		body, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		_ = json.Unmarshal(body, point)

		err := point.Insert(conn)
		if err != nil {
			log.WithError(err).Error("can't create point")

			Response(res, http.StatusBadRequest, Error{Error: "can't create point"})
			return
		}

		Response(res, http.StatusOK, point)
	}
}
