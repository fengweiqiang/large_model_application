# 大模型应用练习

## 准备环境

### 1. 安装ollama`https://ollama.com/`
### 2. 启动ollama `ollama serve`
### 3. 拉取模型 `https://ollama.com/library` `ollama pull [model]`
### 4.下载milvus向量数据库 `wget https://github.com/milvus-io/milvus/releases/download/v2.3.3/milvus-standalone-docker-compose.yml -O docker-compose.yml`
### 5.启动milvus `docker compose up -d`
### 6.拉取代码 `git clone https://github.com/fengweiqiang/large_model_application.git`
### 7.拉取依赖 `go mod tidy`
### 8.运行程序 `go run main.go`

## api
https://console-docs.apipost.cn/preview/4a6af016b75237d8/80e5a02e5d8c2d59

## 功能介绍

- [x] 大模型翻译器功能demo
- [x] 图片转文字demo `(ollama run llava:7b)`
- [x] 切分数据demo
- [x] 建立知识库，使用内嵌模型解析成向量数据存入milvus向量数据库
- [x] 让大模型具备知识库的能力( 本地文件 和 内外网 URL 地址)，对大模型提问知识库内容
- [x] 大模型上下文记忆能力


## 感谢B站大佬提供入门教学

kilmerfun <https://space.bilibili.com/341182735>
五里墩茶社 <https://space.bilibili.com/615957867>
