package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"todoist/syncer/cmd/todoist"

	"github.com/TylerBrock/colorjson"
	"github.com/joho/godotenv"
)

func printJSON(data []byte) {
	var obj map[string]interface{}
	json.Unmarshal(data, &obj)

	f := colorjson.NewFormatter()
	f.Indent = 4

	s, _ := f.Marshal(obj)
	fmt.Println(string(s))
}

func printStruct(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(string(b))
}

type ResourceType string

const (
	ResourceTypeProjects ResourceType = "projects"
	ResourceTypeUser     ResourceType = "user"
	ResourceTypeItems    ResourceType = "items"
)

func PerformSync(apiToken string, resourceTypes []ResourceType) (todoist.SyncResponse, error) {
	url := "https://api.todoist.com/sync/v9/sync"

	resourceTypesJSON, err := json.Marshal(resourceTypes)
	if err != nil {
		return todoist.SyncResponse{}, fmt.Errorf("error marshalling resource types: %w", err)
	}

	requestBodyStr := fmt.Sprintf(`sync_token=*&resource_types=%s`, string(resourceTypesJSON))
	requestBody := []byte(requestBodyStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return todoist.SyncResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return todoist.SyncResponse{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return todoist.SyncResponse{}, fmt.Errorf("error reading response body: %w", readErr)
	}

	var syncResponse todoist.SyncResponse
	if err = json.Unmarshal(body, &syncResponse); err != nil {
		return todoist.SyncResponse{}, fmt.Errorf("error parsing response: %w", err)
	}

	return syncResponse, nil
}

type envVar struct {
	apiToken string
}

func (e *envVar) getEnvVar() {
	e.apiToken = os.Getenv("API_TOKEN")

	if e.apiToken == "" {
		log.Fatal("API_TOKEN environment variable is not set")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config := envVar{}
	config.getEnvVar()

	// resourceTypes := []ResourceType{ResourceTypeItems}
	// syncResponse, err := PerformSync(config.apiToken, resourceTypes)
	// if err != nil {
	// 	log.Fatalf("Error performing sync: %v", err)
	// }

	// if syncResponse.User != nil {
	// 	printStruct(syncResponse.User)
	// }

	// if syncResponse.Items != nil {
	// 	printStruct(syncResponse.Items)
	// }

	todostClient := todoist.NewClient(config.apiToken)
	todostClient.GetAllCompletedTasks()
}
