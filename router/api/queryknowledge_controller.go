package api

import (
	"gitee.com/fengweiqiang/largeModel/config"
	"gitee.com/fengweiqiang/largeModel/util"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/memory"
	"log"
	"net/http"
)

// 让大模型根据知识库内容回答提问问题
type QueryKnowledgeController struct {
}

type queryKnowledgeRequest struct {
	//Title    string `json:"title" binding:"required"`
	Question string `json:"question" binding:"required"`
}

func (QueryKnowledgeController) QueryKnowledge(ctx *gin.Context) {
	var req queryKnowledgeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	llm, err := config.GetLoadLLm(config.MODEL_QWEN4B)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	documents, err := util.QueryMilvusEmbedding(ctx, config.MODEL_QWEN4B, req.Question)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(documents) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error documents not foud"})
		return
	}
	history := memory.NewChatMessageHistory()
	for _, doc := range documents {
		history.AddAIMessage(ctx, doc.PageContent)
	}

	conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))
	executor := agents.NewExecutor(
		agents.NewConversationalAgent(llm, nil),
		nil,
		agents.WithMemory(conversation),
	)
	options := []chains.ChainCallOption{
		chains.WithTemperature(0.8),
	}
	res, err := chains.Run(ctx, executor, req.Question, options...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": res})
}
