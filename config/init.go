package config

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
)

var Llm *ollama.LLM

func InitLlm() {
	llm, err := ollama.New(ollama.WithModel("qwen:4b"))
	if err != nil {
		panic(err)
	}
	response, err := llm.Call(context.Background(), "hello", llms.WithTemperature(0.1))
	if err != nil {
		panic(err)
	}
	log.Printf("llm response: %s", response)
	Llm = llm
}
