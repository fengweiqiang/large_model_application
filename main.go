package main

import (
	"gitee.com/fengweiqiang/largeModel/config"
	"gitee.com/fengweiqiang/largeModel/router"
)

func main() {
	config.InitLlm()
	router.RunServer("8080")
}
