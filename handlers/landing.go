package handlers

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/tingold/boring/model"
)
//Landing implements the OGC Feature Service Landing page
func Landing(ctx *gin.Context){

	url := location.Get(ctx)

	base := url.Scheme+"://"+url.Host+"/"

	links := []model.Link{
		{
			Href:  base,
			Rel:   "self",
			Title: "this document",
			Type: "application/json",
		},
		{
			Href:  base+"api",
			Rel:   "service-desc",
			Title: "API Definition",
			Type: "application/vnd.oai.openapi+json;version=3.0",
		},
		{
			Href:  base+"conformance",
			Rel:   "conformance",
			Title: "Conformance classes implemented by this service",
			Type: "application/json",
		},
		{
			Href:  base+"collections",
			Rel:   "data",
			Title: "Feature Collections",
			Type: "application/json",
		},
	}

	landing := model.Landing{
		Title:       "Boring Server",
		Description: "A boring implementation of the OGC Feature API",
		Links:    &links,
	}

	ctx.JSON(200, landing)

}