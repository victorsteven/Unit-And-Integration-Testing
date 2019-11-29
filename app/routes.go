package app

import "efficient-api/controllers"

func routes() {
	router.GET("/messages/:message_id", controllers.GetMessage)
	router.POST("/messages", controllers.CreateMessage)
	router.GET("/ping", controllers.Ping)

}
