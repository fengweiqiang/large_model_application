package controller

import (
	"context"
	"fmt"
	"gitee.com/fengweiqiang/largeModel/config"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"io"
	"net/http"
	"time"
)

type QuestionController struct {
}

type QuestionRequest struct {
	Question string `json:"question" binding:"required"`
}

func (q QuestionController) Question(g *gin.Context) {
	var request QuestionRequest
	err := g.BindJSON(&request)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	call, err := config.Llm.Call(g, request.Question, llms.WithTemperature(0.9)) //llms.WithStopWords([]string{"帮助"}),

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"response": call})
}

func (q QuestionController) QuestionStream(g *gin.Context) {
	request, b := g.GetQuery("question")
	if !b {
		g.JSON(http.StatusBadRequest, gin.H{"error": "request param not found"})
		return
	}
	c := make(chan []byte, 256)
	go func() {
		_, err := config.Llm.Call(g, request, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			c <- chunk
			fmt.Println(string(chunk))
			//g.SSEvent("message", string(chunk))
			return nil
		}))
		if err != nil {
			//todo 错误处理
			return
		}
	}()
	g.Stream(func(w io.Writer) bool {
		time.Sleep(time.Second)
		g.SSEvent("message", string(<-c))
		//todo 错误处理
		return true
	})

}
