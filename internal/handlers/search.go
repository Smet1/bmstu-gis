package handlers

import (
	"github.com/Smet1/bmstu-gis/internal/logger"
	"github.com/Smet1/bmstu-gis/internal/pathfinding"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func GetPathFindingHandler(conn *sqlx.DB, bmstuMap *pathfinding.Map) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log := logger.GetLogger(req.Context())
		from := req.URL.Query().Get("from")
		to := req.URL.Query().Get("to")
		path, err := bmstuMap.Path(from, to)
		if err != nil {
			log.WithError(err).Error("can't get points")

			Response(res, http.StatusBadRequest, Error{Error: err.Error()})
			return
		}

		Response(res, http.StatusOK, path)
	}
}
