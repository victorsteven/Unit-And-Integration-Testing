package app

import "efficient-api/controllers"

func routes() {
	router.GET("/messages/:message_id", controllers.GetMessage)
	router.GET("/messages", controllers.GetAllMessages)
	router.POST("/messages", controllers.CreateMessage)
	router.PUT("/messages/:message_id", controllers.UpdateMessage)
	router.DELETE("/messages/:message_id", controllers.DeleteMessage)


	//router.GET("/ping", controllers.Ping)
}
