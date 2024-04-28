package api

import (
	"fmt"
	"gitee.com/fengweiqiang/largeModel/config"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ImageToTextController struct {
}

func (i ImageToTextController) ImageToText(ctx *gin.Context) {
	des, _ := ctx.GetPostForm("des")
	var err error
	parts := []llms.ContentPart{}
	var savePath = "./tmp/"
	if ctx.Request.MultipartForm != nil {
		fileMap := ctx.Request.MultipartForm.File
		for k, v := range fileMap {
			for i, vf := range v {
				trimExt := strings.TrimLeft(filepath.Ext(vf.Filename), ".")
				tmpflie := fmt.Sprintf("%s%d%s", savePath, time.Now().Unix(), vf.Filename)
				err = ctx.SaveUploadedFile(vf, tmpflie)
				defer os.Remove(tmpflie)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
					return
				}
				datas, err := os.ReadFile(tmpflie)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
					return
				}
				parts = append(parts, llms.BinaryPart("image/"+trimExt, datas))
				log.Println(k, i, vf.Filename)
			}
		}
	}
	if strings.TrimSpace(des) == "" {
		parts = append(parts, llms.TextPart("介绍一下自己，并举例说明能干哪些事情"))
	} else {
		parts = append(parts, llms.TextPart(des))
	}
	content := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: parts,
		},
	}
	generateContent, err := config.Llm.GenerateContent(ctx, content, llms.WithModel("llava:7b"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, generateContent)

}
