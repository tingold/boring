package model

import "time"

type Extent struct {

	SpatialExtent *SpatialExtent `json:"spatial"`
	TemporalExtent *TemporalExtent `json:"temporal"`

}

func NewSpatialExtent() SpatialExtent {

	return SpatialExtent{
		CRS: "http://www.opengis.net/def/crs/OGC/1.3/CRS84",
		BBox: nil,
	}
}

type SpatialExtent struct {

	CRS string `json:"crs"`
	BBox *[]float64 `json:"bbox"`
}

func NewTemporalExtent() TemporalExtent {

	return TemporalExtent{
		Interval: nil,
		TRS:      &[]string{"http://www.opengis.net/def/uom/ISO-8601/0/Gregorian"},
	}

}

type TemporalExtent struct {

	Interval *[]time.Time
	TRS *[]string `json:"trs"`
}