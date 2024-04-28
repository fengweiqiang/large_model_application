package router

import (
	"gitee.com/fengweiqiang/largeModel/router/api"
	"github.com/gin-gonic/gin"
)

func RunServer(port string) {
	r := gin.Default()
	group := r.Group("/app/v1")
	group.POST("/question", api.QuestionController{}.Question)
	group.GET("/questionStream", api.QuestionController{}.QuestionStream)

	group.POST("/template", api.TemplateController{}.Template)
	group.POST("/imagetotext", api.ImageToTextController{}.ImageToText)
	group.POST("/embedding", api.EmbeddingController{}.Embedding)
	group.POST("/knowledge", api.KnowledgeController{}.Knowledge)
	group.POST("/queryKnowledge", api.QueryKnowledgeController{}.QueryKnowledge)
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
