package config

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
)

var Llm *ollama.LLM

type MODEL string

func (m MODEL) ToString() string {
	return string(m)
}

const (
	MilieusDBAddress string = "127.0.0.1:19530"
	MODEL_LLAVA7B    MODEL  = "llava:7b"
	MODEL_QWEN4B     MODEL  = "qwen:4b"
)

func InitLlm() {
	llm, err := ollama.New(ollama.WithModel(MODEL_LLAVA7B.ToString()))
	if err != nil {
		panic(err)
	}
	response, err := llm.Call(context.Background(), "who are you?", llms.WithTemperature(0.1))
	if err != nil {
		panic(err)
	}
	log.Printf("llm response: %s\n", response)
	Llm = llm
}
func GetLoadLLm(model MODEL) (*ollama.LLM, error) {
	llm, err := ollama.New(ollama.WithModel(model.ToString()))

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return llm, nil
}
