package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"graphql-taxonomy/controller"
	_ "graphql-taxonomy/dao"
)

var (
	r = gin.Default()
)

func init() {

	v1 := r.Group("/v1")
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	{
		v1.GET("/ping", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"status": "OK",
			})
		})

		v1.POST("/getUsers", controller.GetUserHandler)
		v1.POST("/get-fb-audience", controller.GetAudience)
	}
	r.Run(":8080")
}
func main() {

}



