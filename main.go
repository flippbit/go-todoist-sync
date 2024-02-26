package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
	"todoist/syncer/cmd/files"
	"todoist/syncer/cmd/todoist"

	"github.com/joho/godotenv"
)

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

func SaveStructToJSON(data interface{}, path string) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config := envVar{}
	config.getEnvVar()

	todostClient := todoist.NewClient(config.apiToken)
	items, err := todostClient.GetAllCompletedTasks()
	if err != nil {
		log.Fatalf("Error getting completed tasks: %v", err)
	}

	fmt.Printf("Fetched %d completed tasks\n", len(items))

	groupedItems := todoist.GroupItemsByDate(items)

	fmt.Printf("Total completed tasks fetched: %d\n", len(items))
	fmt.Printf("Total dates: %d\n", len(groupedItems))

	var allDates []string
	for date := range groupedItems {
		allDates = append(allDates, date)
	}

	sort.Slice(allDates, func(i, j int) bool {
		ti, _ := time.Parse("2006-01-02", allDates[i])
		tj, _ := time.Parse("2006-01-02", allDates[j])
		return ti.Before(tj)
	})

	for _, date := range allDates {
		itemsOnDate := groupedItems[date]
		fmt.Printf("%s :::\n", date)

		var itemsAsText string
		for _, item := range itemsOnDate {
			itemsAsText += fmt.Sprintf("- [x] %s #%s\n", item.Content, item.Project)
		}

		files.CreateFile(fmt.Sprintf("out/%s.md", date), itemsAsText)
	}

}
