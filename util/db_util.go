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

type DBCollection string

func (c DBCollection) ToString() string {
	return string(c)
}

const (
	DB_COLLECTION_LOCAL_FILE DBCollection = "local_file"
	DB_COLLECTION_URL        DBCollection = "url"
)

func SaveMilvusEmbedding(ctx context.Context, model config.MODEL, collection DBCollection, docs []schema.Document) error {
	//文档放入不同的集合里
	store, err := loadStore(ctx, &model, collection.ToString())
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

func QueryMilvusEmbedding(ctx context.Context, model config.MODEL, collection DBCollection, prompt string) (docs []schema.Document, err error) {
	store, err := loadStore(ctx, &model, collection.ToString())
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

func loadStore(ctx context.Context, model *config.MODEL, collection string) (milvus.Store, error) {
	llm := config.Llm
	var err error
	if model != nil {
		llm, err = config.GetLoadLLm(*model)
		if err != nil {
			log.Println(err)
			return milvus.Store{}, err
		}
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
	store, err := milvus.New(ctx, clientConfig, milvus.WithEmbedder(embedder), milvus.WithIndex(autoindex), milvus.WithCollectionName(collection))
	if err != nil {
		log.Println(err)
		return milvus.Store{}, err
	}
	return store, nil
}
