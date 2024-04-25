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

func (q QuestionController) Question(ctx *gin.Context) {
	var request QuestionRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	call, err := config.Llm.Call(ctx, request.Question, llms.WithTemperature(0.9)) //llms.WithStopWords([]string{"帮助"}),

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"response": call})
}

func (q QuestionController) QuestionStream(ctx *gin.Context) {
	request, b := ctx.GetQuery("question")
	if !b {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "request param not found"})
		return
	}
	c := make(chan []byte, 256)
	go func() {
		defer close(c)
		_, err := config.Llm.Call(ctx, request, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			c <- chunk
			fmt.Println(string(chunk))
			//ctx.SSEvent("message", string(chunk))
			return nil
		}))
		if err != nil {
			//todo 错误处理
			return
		}
	}()
	ctx.Stream(func(w io.Writer) bool {
		time.Sleep(time.Second)
		if data, ok := <-c; ok {
			ctx.SSEvent("message", string(data))
			return true
		} else {
			return false
		}
	})
}
