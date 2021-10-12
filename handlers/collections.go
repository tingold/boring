package handlers

import (
	"context"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/tingold/boring/catalog"
)

//Collections implements the OGC Features API conformance endpoint
func Collections(ctx *gin.Context){

	cat, err := catalog.GetCatalog()
	if err != nil {
		ctx.AbortWithError(500, err)
	}
	url := location.Get(ctx)
	base := url.Scheme+"://"+url.Host+"/"
	returns, err := cat.GetCollections(context.WithValue(ctx, "boring_base", base))
	if err != nil{
		ctx.JSON(500,err.Error())
	}
	ctx.JSON(200, returns)

}
