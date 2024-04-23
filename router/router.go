package router

import (
	"gitee.com/fengweiqiang/largeModel/controller"
	"github.com/gin-gonic/gin"
)

func RunServer(port string) {
	r := gin.Default()
	group := r.Group("/app/v1")
	group.POST("/question", controller.QuestionController{}.Question)
	group.GET("/questionStream", controller.QuestionController{}.QuestionStream)
	group.GET("/questionStream2", controller.QuestionController{}.QuestionStream2)
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
