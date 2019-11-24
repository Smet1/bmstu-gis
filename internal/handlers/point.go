package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Smet1/bmstu-gis/internal/db"
	"github.com/Smet1/bmstu-gis/internal/logger"
	"github.com/jmoiron/sqlx"
)

const (
	TypeStairs = "stairs"
	TypeAll    = "all"
	TypeEnds   = "ends"
)

type Points struct {
	Data *db.ListPoints `json:"data"`
	Type string         `json:"type"`
}

func GetGetPointsHandler(conn *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())
		pointType := req.URL.Query().Get("type")

		listPoints := &db.ListPoints{}
		var err error
		typeRequest := ""
		switch pointType {
		case TypeStairs:
			typeRequest = TypeStairs
			err = listPoints.GetStairs(conn)
		case TypeEnds:
			typeRequest = TypeEnds
			err = listPoints.GetEnds(conn)
		default:
			typeRequest = TypeAll
			err = listPoints.Get(conn)
		}

		if err != nil {
			log.WithError(err).Error("can't get points")

			Response(res, http.StatusBadRequest, Error{Error: "can't get points"})
			return
		}

		points := Points{
			Data: listPoints,
			Type: typeRequest,
		}

		Response(res, http.StatusOK, points)
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
