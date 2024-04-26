package controller

import (
	"gitee.com/fengweiqiang/largeModel/config"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"log"
	"net/http"
)

type TemplateController struct {
}
type templateRequest struct {
	Input      string `json:"input"  binding:"required"`
	OutputLang string `json:"outputLang"  binding:"required"`
}

func (t TemplateController) Template(ctx *gin.Context) {
	req := templateRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate("你是一名精通各国语言的翻译专家", nil),
		prompts.NewHumanMessagePromptTemplate(
			`请把这句话 {{.input}} 翻译成 {{.outputLang}} 语言`,
			[]string{"input", "outputLang"},
		),
	})

	msgs, err := prompt.FormatMessages(map[string]any{
		"outputLang": req.OutputLang,
		"input":      req.Input,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	messageContents := []llms.MessageContent{}
	for _, message := range msgs {
		messageContents = append(messageContents, llms.TextParts(message.GetType(), message.GetContent()))
	}
	generateContent, err := config.Llm.GenerateContent(ctx, messageContents)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"response": generateContent})
}
