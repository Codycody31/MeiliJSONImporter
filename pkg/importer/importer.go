package importer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/meilisearch/meilisearch-go"
	"github.com/urfave/cli/v2"
	"go.codycody31.dev/MeiliJSONImporter/pkg/utils"
)

func SetupFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Value: "http://127.0.0.1:7700",
			Usage: "The MeiliSearch host URL. Default is http://127.0.0.1:7700.",
		},
		&cli.StringFlag{
			Name:    "master-key",
			Aliases: []string{"m"},
			Value:   "",
			Usage:   "The master key for MeiliSearch, required if your instance uses authentication.",
		},
		&cli.StringFlag{
			Name:     "index",
			Aliases:  []string{"i"},
			Value:    "",
			Usage:    "The name of the MeiliSearch index to which the JSON data will be pushed.",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "json",
			Aliases:  []string{"j"},
			Value:    "",
			Usage:    "Path to the JSON file containing the data to be imported.",
			Required: true,
		},
		&cli.IntFlag{
			Name:  "batch-size",
			Value: 10 * 1024 * 1024, // 10 MiB
			Usage: "Maximum batch size in bytes for pushing documents.",
		},
	}
}

func HandleImport(c *cli.Context) error {
	log.Println("Starting MeiliJSONImporter")

	host := c.String("host")
	masterKey := c.String("master-key")
	index := c.String("index")
	jsonPath := c.String("json")
	batchSize := c.Int("batch-size")

	// Ensure the batch size is at least 1 MiB and less than 95 MiB
	if batchSize < 1024*1024 {
		return fmt.Errorf("batch size must be at least 1 MiB")
	}
	if batchSize > 95*1024*1024 {
		return fmt.Errorf("batch size must be less than 95 MiB")
	}

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   host,
		APIKey: masterKey,
	})
	indexObj := client.Index(index)

	file, err := os.Open(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	log.Printf("Pushing data to index: %s", index)

	// Total size of the JSON file in mib
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	log.Printf("Total file size: %.2f MiB", float64(fileInfo.Size())/1024/1024)

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	// Expect the start of an array
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("failed to read starting token of JSON: %v", err)
	}

	// Calculate bases on 10 random objects sizes the average size of an object and the expected number of objects per batch
	var objects []map[string]interface{}
	expectedNumOfObjectsPerBatch := 50 // placeholder used to calculate the average object size
	for i := 0; i < expectedNumOfObjectsPerBatch; i++ {
		var obj map[string]interface{}
		if err := decoder.Decode(&obj); err != nil {
			return fmt.Errorf("failed to decode JSON object: %v", err)
		}
		objects = append(objects, obj)
	}
	averageObjSize := utils.CalculateBatchSize(objects) / expectedNumOfObjectsPerBatch
	expectedNumOfObjectsPerBatch = batchSize / averageObjSize
	log.Printf("Average object size: %.2f KiB", float64(averageObjSize)/1024)
	log.Printf("Expected number of objects per batch: %d", expectedNumOfObjectsPerBatch)

	objects = make([]map[string]interface{}, 0, expectedNumOfObjectsPerBatch)
	currentBatchSize := 0

	// Processing loop
	for decoder.More() {
		var obj map[string]interface{}
		if err := decoder.Decode(&obj); err != nil {
			return fmt.Errorf("failed to decode JSON object: %v", err)
		}

		encodedObj, err := json.Marshal(obj)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON object: %v", err)
		}
		objSize := len(encodedObj)

		if currentBatchSize+objSize > batchSize {
			// Batch is full, send it
			if err := utils.SendBatch(indexObj, objects); err != nil {
				return err
			}
			objects = make([]map[string]interface{}, 0, expectedNumOfObjectsPerBatch)
			currentBatchSize = 0
		}

		// Add object to batch
		objects = append(objects, obj)
		currentBatchSize += objSize
	}

	// Send any remaining objects
	if len(objects) > 0 {
		if err := utils.SendBatch(indexObj, objects); err != nil {
			return err
		}
	}

	log.Printf("Data pushed successfully to index: %s", index)
	return nil
}
