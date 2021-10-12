package memory

import "github.com/paulmach/orb"

type Table struct {

	Name string
	Schema string
	GeometryColumnName string
	SRID int64
	GeometryType string
	Dimension int64
	Extent orb.Polygon
	//todo: detect temporal bounds?

}