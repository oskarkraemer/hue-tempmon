package main

import "github.com/gin-gonic/gin"

func AddWebserverRoutes(api *gin.Engine) {
	api.LoadHTMLGlob("templates/*")

	api.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "dashboard.html", gin.H{
			"title": "Hello World!",
		})
	})
}
