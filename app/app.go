package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)
func StartApp(){
	fmt.Println("the app is started")
	routes()
	router.Run(":8080")
}