package pathfinding

import (
	"github.com/Smet1/bmstu-gis/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/karalabe/cookiejar.v2/graph"
	"gopkg.in/karalabe/cookiejar.v2/graph/bfs"
)

type Map struct {
	graph       *graph.Graph
	EndPoints   map[string]*db.Point
	StairPoints map[string]*db.Point
	allPoints   map[int64]*db.Point
	Points      db.ListPoints
	DbConn      *sqlx.DB
	// тк в графе начало с 0, а айдишнки могут быть любыми, то нужно хранить эту разницу для правильного перевода точек
	// костыль пезда
	dirtyHackIDs int64
}

func InitMap(conn *sqlx.DB) (*Map, error) {
	points := &db.ListPoints{}
	err := points.Get(conn)
	if err != nil {
		return nil, errors.Wrap(err, "can't get points")
	}

	endPoints := make(map[string]*db.Point)
	stairPoints := make(map[string]*db.Point)
	allPoints := make(map[int64]*db.Point)

	gr := graph.New(len(points.Points))
	minus := points.Points[0].ID
	for i := range points.Points {
		allPoints[points.Points[i].ID] = points.Points[i]

		switch {
		case points.Points[i].Cabinet:
			endPoints[points.Points[i].Name] = points.Points[i]
		case points.Points[i].Stair:
			stairPoints[points.Points[i].Name] = points.Points[i]
		}

		if !points.Points[i].NodeID.Valid {
			continue
		}

		gr.Connect(int(points.Points[i].ID-minus), int(points.Points[i].NodeID.Int64-minus))
	}

	return &Map{
		graph:        gr,
		EndPoints:    endPoints,
		StairPoints:  stairPoints,
		allPoints:    allPoints,
		Points:       *points,
		DbConn:       conn,
		dirtyHackIDs: minus,
	}, nil
}

func (m *Map) pathInt(src, dest int64) []int {
	b := bfs.New(m.graph, int(src-m.dirtyHackIDs))
	return b.Path(int(dest - m.dirtyHackIDs))
}

func (m *Map) Path(srcName, destName string) (*db.ListPoints, error) {
	src, ok := m.EndPoints[srcName]
	if !ok {
		return nil, errors.New("can't find source point")
	}

	dest, ok := m.EndPoints[destName]
	if !ok {
		return nil, errors.New("can't find end point")
	}

	path := m.pathInt(src.ID, dest.ID)
	if len(path) == 0 {
		return nil, errors.New("can't find path between src and dest points")
	}

	points := make([]*db.Point, 0, len(path))
	for i := range path {
		points = append(points, m.allPoints[int64(path[i])+m.dirtyHackIDs])
	}

	return &db.ListPoints{Points: points}, nil
}
