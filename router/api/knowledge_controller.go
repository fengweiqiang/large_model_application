package api

import (
	"fmt"
	"gitee.com/fengweiqiang/largeModel/config"
	"gitee.com/fengweiqiang/largeModel/util"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
	"log"
	"net/http"
	"os"
	"time"
)

// 上传知识库让大模型学习
type KnowledgeController struct {
}

type knowledgeRequest struct {
	Title string `json:"title"`
}

func (k KnowledgeController) Knowledge(ctx *gin.Context) {
	title, b := ctx.GetPostForm("title")
	if !b {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "请输入知识库标题"})
		return
	}
	savePath := "./knowledge/"
	fileDatas := make(map[string]string)
	if ctx.Request.MultipartForm != nil {
		for _, v := range ctx.Request.MultipartForm.File {
			for _, f := range v {
				tmpflie := fmt.Sprintf("%s%d%s", savePath, time.Now().Unix(), f.Filename)
				err := ctx.SaveUploadedFile(f, tmpflie)
				defer os.Remove(tmpflie)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
					return
				}
				fileDatas[tmpflie] = tmpflie
			}
		}
	}
	for k, v := range fileDatas {
		file, _ := os.Open(v)
		p := documentloaders.NewText(file)
		split := textsplitter.NewRecursiveCharacter()
		split.ChunkSize = 300
		split.ChunkOverlap = 30
		docs, err := p.LoadAndSplit(ctx, split)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		log.Println(k, len(docs))
		//查询切割转向量
		err = util.SaveMilvusEmbedding(ctx, config.MODEL_LLAVA7B, docs)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": title + " ok"})
}
