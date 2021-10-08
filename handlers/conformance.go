package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tingold/boring/model"
)

//Conformance implements the OGC Features API confromance endpoint
func Conformance(ctx *gin.Context){

	conformance := model.Conformance{ConformsTo: &[]string{
		"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/core",
		"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/oas30",
		"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/geojson",

	}}
	ctx.JSON(200, conformance)
}