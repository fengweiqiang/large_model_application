package api

import (
	"fmt"
	"gitee.com/fengweiqiang/largeModel/config"
	"gitee.com/fengweiqiang/largeModel/util"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// 上传知识库让大模型学习
type KnowledgeController struct {
}

type knowledgeRequest struct {
	Title      string `json:"title"`
	ContentURL string `json:"contentURL"`
}

func (k KnowledgeController) Knowledge(ctx *gin.Context) {
	title, b := ctx.GetPostForm("title")
	if !b {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "请输入知识库标题"})
		return
	}
	contentURL, existURL := ctx.GetPostForm("contentURL")
	savePath := "./knowledge/"
	fileDatas := make(map[string]string)
	if existURL {
		//以外网地址 https://baike.baidu.com/item/2024%E5%B9%B4%E5%B7%B4%E9%BB%8E%E5%A5%A5%E8%BF%90%E4%BC%9A/17619118?fromModule=home_hotspot 为例
		response, err := http.Get(contentURL)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}
		defer response.Body.Close()
		datas, _ := io.ReadAll(response.Body)
		tmpflie := fmt.Sprintf("%s%d%s", savePath, time.Now().Unix(), title)
		_ = os.WriteFile(tmpflie, datas, 0666)
		fileDatas[tmpflie] = tmpflie
	} else if ctx.Request.MultipartForm != nil {
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
		var loader documentloaders.Loader
		// todo 根据后缀名判断加载对应文件格式
		if existURL {
			//todo 优化：去除一些不必要的html标签，加快速度
			loader = documentloaders.NewHTML(file)
		} else {
			loader = documentloaders.NewText(file)
			//fileInfo, _ := file.Stat()
			//loader = documentloaders.NewPDF(file, fileInfo.Size())
		}
		split := textsplitter.NewRecursiveCharacter()
		split.ChunkSize = 300
		split.ChunkOverlap = 30
		docs, err := loader.LoadAndSplit(ctx, split)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		log.Println(k, len(docs))
		//查询切割转向量
		if existURL {
			err = util.SaveMilvusEmbedding(ctx, config.MODEL_LLAVA7B, util.DB_COLLECTION_URL, docs)
		} else {
			err = util.SaveMilvusEmbedding(ctx, config.MODEL_LLAVA7B, util.DB_COLLECTION_LOCAL_FILE, docs)
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": title + " ok"})
}
