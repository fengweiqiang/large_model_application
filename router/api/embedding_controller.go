package api

import (
	"gitee.com/fengweiqiang/largeModel/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmbeddingController struct {
}

type embeddingRequest struct {
	Text string `json:"text" binding:"required"`
}

func (e EmbeddingController) Embedding(ctx *gin.Context) {
	request := embeddingRequest{}
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	texts := []string{}
	texts = append(texts, request.Text)
	emb, err := config.Llm.CreateEmbedding(ctx, texts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	//log.Println("嵌入向量的数量：", len(emb))
	//for i, e := range emb {
	//	fmt.Printf("%d: %v...\n", i, e)
	//}
	ctx.JSON(http.StatusOK, gin.H{"message": len(emb)})
}
