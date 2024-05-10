package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func ServeAPI(db *sql.DB, devices []Device) {
	api := gin.Default()

	AddWebserverRoutes(api)

	api.GET("/api/v1/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Running hueTemperatureMonitor v.0.0.1 / api ver. 1",
		})
	})

	api.GET("/api/v1/sensors", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"sensors": devices})
	})

	api.GET("/api/v1/entries/:table/:min_date/:sensor_id", func(ctx *gin.Context) {
		var params RequestParmas
		if err := ctx.ShouldBindUri(&params); err != nil {
			ctx.JSON(400, gin.H{"message": err.Error()})
			return
		}

		entires, err := QuerySensorEntries(db, params.Table, params.MinDate, params.SensorID)

		if err != nil {
			ctx.JSON(400, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(200, gin.H{"items": entires})
		}
	})

	api.GET("/api/v1/entries/:table/:min_date/", func(ctx *gin.Context) {
		var params RequestParmas
		if err := ctx.ShouldBindUri(&params); err != nil {
			ctx.JSON(400, gin.H{"message": err.Error()})
			return
		}

		entires, err := QuerySensorEntries(db, params.Table, params.MinDate, "")

		if err != nil {
			ctx.JSON(400, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(200, gin.H{"items": entires})
		}
	})

	api.Run()
}
