package pathfinding

import (
	"github.com/Smet1/bmstu-gis/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/karalabe/cookiejar.v2/graph"
	"gopkg.in/karalabe/cookiejar.v2/graph/bfs"
)

type Map struct {
	graph *graph.Graph
	//BFS    *bfs.Bfs
	Points db.ListPoints
	DbConn *sqlx.DB
}

func InitMap(conn *sqlx.DB) (*Map, error) {
	points := &db.ListPoints{}
	err := points.Get(conn)
	if err != nil {
		return nil, errors.Wrap(err, "can't get points")
	}

	mp := make(map[int64]*db.Point)
	for i := range points.Points {
		mp[points.Points[i].ID] = points.Points[i]
	}

	gr := graph.New(len(points.Points))
	//b := bfs.New(gr, int(points.Points[0].ID))
	minus := points.Points[0].ID
	for i := range points.Points {
		if !points.Points[i].NodeID.Valid {
			continue
		}

		gr.Connect(int(points.Points[i].ID-minus), int(points.Points[i].NodeID.Int64-minus))
	}

	return &Map{
		graph: gr,
		//BFS:    b,
		Points: *points,
		DbConn: conn,
	}, nil
}

func (m *Map) Path(src, dest int64) []int {
	b := bfs.New(m.graph, int(src))
	return b.Path(int(dest))
}
