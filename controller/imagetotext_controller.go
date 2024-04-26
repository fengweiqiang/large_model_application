package controller

import (
	"gitee.com/fengweiqiang/largeModel/config"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ImageToTextController struct {
}

func (i ImageToTextController) ImageToText(ctx *gin.Context) {
	des, _ := ctx.GetPostForm("des")
	var err error
	parts := []llms.ContentPart{}
	var savePath = "./tmp/"
	fileMap := ctx.Request.MultipartForm.File
	for k, v := range fileMap {
		for i, vf := range v {
			trimExt := strings.TrimLeft(filepath.Ext(vf.Filename), ".")
			tmpflie := savePath + vf.Filename
			err = ctx.SaveUploadedFile(vf, tmpflie)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
				return
			}
			defer os.Remove(tmpflie)
			datas, err := os.ReadFile(tmpflie)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
				return
			}
			parts = append(parts, llms.BinaryPart("image/"+trimExt, datas))
			log.Println(k, i, vf.Filename)
		}
	}
	parts = append(parts, llms.TextPart(des))
	content := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: parts,
		},
	}
	generateContent, err := config.Llm.GenerateContent(ctx, content, llms.WithModel("llava:7b"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	ctx.JSON(http.StatusOK, generateContent)

}
