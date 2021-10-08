package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/location"
	gzap "github.com/gin-contrib/zap"
	"github.com/spf13/viper"
	"github.com/tingold/boring/config"
	"github.com/tingold/boring/handlers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)




func main(){

	preflight()
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gzap.Ginzap(zap.L(),time.RFC3339,true))
	router.Use(gzap.RecoveryWithZap(zap.L(), true))

	if !viper.GetBool(config.SERVER_DEBUG){
		gin.SetMode(gin.ReleaseMode)
	}

	if viper.GetBool(config.SERVER_OVERRIDE_REQUESTED_HOST){
		router.Use(location.New(location.Config{
			Scheme:  viper.GetString(config.SERVER_PROTO),
			Host:    viper.GetString(config.SERVER_HOSTNAME),
			Base:    viper.GetString("/"),
		}))
	} else {
		router.Use(location.Default())
	}


	//landing page
	router.GET("/",handlers.Landing)
	//conformance
	router.GET("/conformance", handlers.Conformance)
	//collections
	router.GET("/collections", handlers.Collections)

	address := ":"+strconv.Itoa(viper.GetInt(config.SERVER_PORT))
	if viper.GetBool(config.SSL_ENABLE) {
		router.RunTLS(address, viper.GetString(config.SSL_CERTFILE), viper.GetString(config.SSL_KEYFILE))
	} else {
		router.Run(address)
	}

}

func preflight(){

	viper.New()

	//server settings
	viper.SetDefault(config.SERVER_DEBUG, false)
	viper.SetDefault(config.SERVER_PORT, 8000)
	viper.SetDefault(config.SERVER_HEALTHCHECK, true)
	viper.SetDefault(config.SERVER_HEALTHCHECK_PATH, "/health")
	//todo: revist this when not in hardcore dev
	viper.SetDefault(config.LOG_LEVEL,"INFO")


	//server hostname settings (for proxies)
	viper.SetDefault(config.SERVER_OVERRIDE_REQUESTED_HOST, false)
	viper.SetDefault(config.SERVER_ADVERTISE_PORT, 80)
	viper.SetDefault(config.SERVER_HOSTNAME, "localhost")
	viper.SetDefault(config.SERVER_PROTO, "http")

	//ssl settings
	//todo: implement SSL handling
	viper.SetDefault(config.SSL_ENABLE, false)

	//database
	viper.SetDefault(config.POSTGRES_HOST,"localhost")
	viper.SetDefault(config.POSTGRES_USER,"postgis")
	viper.SetDefault(config.POSTGRES_PASSWORD,"password")
	viper.SetDefault(config.POSTGRES_PORT, 5432)
	viper.SetDefault(config.POSTGRES_DATABASE, "postgis")
	viper.SetDefault(config.POSTGRES_READONLY, true)
	viper.SetDefault(config.POSTGRES_FAILFAST, true)

	viper.AutomaticEnv()

	loggerConfig := zap.NewProductionConfig()
	//sampling will reduce the amount of logs written when the server is under load
	//disable if we are in debug
	if viper.GetBool(config.SERVER_DEBUG) {
		loggerConfig.Sampling = nil
	}

	loggerConfig.Level.UnmarshalText([]byte(viper.GetString(config.LOG_LEVEL)))
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	loggerConfig.EncoderConfig.TimeKey = "ts"
	loggerConfig.EncoderConfig.LevelKey = "l"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.ErrorOutputPaths = []string{"stderr"}
	logger, _ := loggerConfig.Build()
	zap.ReplaceGlobals(logger)
}