package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/meilisearch/meilisearch-go"
)

func SendBatch(indexObj *meilisearch.Index, objects []map[string]interface{}) error {
	data, err := json.Marshal(objects)
	if err != nil {
		return fmt.Errorf("failed to marshal objects: %v", err)
	}
	if _, err := indexObj.UpdateDocuments(data); err != nil {
		return fmt.Errorf("failed to update documents in MeiliSearch: %v", err)
	}
	fmt.Printf("Pushing batch of size %d bytes to index %s\n", len(data), indexObj.UID)
	return nil
}

// Calcs the batch size in bytes
func CalculateBatchSize(objects []map[string]interface{}) int {
	var size int
	for _, obj := range objects {
		data, err := json.Marshal(obj)
		if err != nil {
			log.Fatalf("failed to marshal object: %v", err)
		}
		size += len(data)
	}
	return size
}
