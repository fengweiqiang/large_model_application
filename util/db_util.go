package util

import (
	"context"
	"gitee.com/fengweiqiang/largeModel/config"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/milvus"
	"log"
)

func SaveMilvusEmbedding(ctx context.Context, model config.MODEL, docs []schema.Document) error {
	store, err := loadStore(ctx, model)
	if err != nil {
		log.Println(err)
		return err
	}

	documents, err := store.AddDocuments(ctx, docs)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(documents)
	return nil
}

func QueryMilvusEmbedding(ctx context.Context, model config.MODEL, prompt string) (docs []schema.Document, err error) {
	store, err := loadStore(ctx, model)
	if err != nil {
		log.Println(err)
		return []schema.Document{}, err
	}
	docRetrieved, err := store.SimilaritySearch(ctx, prompt, 3)
	if err != nil {
		return nil, err
	}

	return docRetrieved, nil
}

func loadStore(ctx context.Context, model config.MODEL) (milvus.Store, error) {
	llm, err := config.GetLoadLLm(model)
	if err != nil {
		log.Println(err)
		return milvus.Store{}, err
	}

	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Println(err)
		return milvus.Store{}, err
	}
	clientConfig := client.Config{
		Address: config.MilieusDBAddress,
	}
	autoindex, _ := entity.NewIndexAUTOINDEX(entity.L2)
	store, err := milvus.New(ctx, clientConfig, milvus.WithEmbedder(embedder), milvus.WithIndex(autoindex))
	if err != nil {
		log.Println(err)
		return milvus.Store{}, err
	}
	return store, nil
}
