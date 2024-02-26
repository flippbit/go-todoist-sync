package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	apiToken string
}

func NewClient(apiToken string) *Client {
	return &Client{
		apiToken: apiToken,
	}
}

func (c *Client) doRequest(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *Client) GetCompletedTasks(limit int, offset int) ([]ParsedItem, error) {
	url := "https://api.todoist.com/sync/v9/completed/get_all" + fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	responseBody, err := c.doRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error making API request: %v", err)
		return nil, err
	}

	var completedTaskResponse CompletedTasksResponse
	err = json.Unmarshal(responseBody, &completedTaskResponse)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return nil, err
	}

	parsedItems := []ParsedItem{}

	for _, item := range completedTaskResponse.Items {
		parsedItem := ParsedItem{
			Id:          item.Id,
			TaskId:      item.TaskId,
			Content:     item.Content,
			CompletedAt: item.CompletedAt,
			Project:     completedTaskResponse.Projects[item.ProjectId].Name,
		}
		parsedItems = append(parsedItems, parsedItem)
	}

	return parsedItems, nil
}

func (c *Client) GetAllCompletedTasks() ([]ParsedItem, error) {
	allTasks := []ParsedItem{}
	offset := 0
	limit := 200

	for {
		tasks, err := c.GetCompletedTasks(limit, offset)
		log.Printf("Fetched %d completed tasks", len(tasks))
		if err != nil {
			log.Printf("Error getting completed tasks: %v", err)
			break
		}

		if len(tasks) == 0 {
			break
		}

		allTasks = append(allTasks, tasks...)
		offset += limit
	}

	log.Printf("Total completed tasks fetched: %d", len(allTasks))

	return allTasks, nil
}

func GroupItemsByDate(items []ParsedItem) map[string][]ParsedItem {
	groupedItems := map[string][]ParsedItem{}

	for _, item := range items {
		date := item.CompletedAt.Format("2006-01-02")
		groupedItems[date] = append(groupedItems[date], item)
	}

	return groupedItems
}

func ConvertItemToText(item ParsedItem) string {
	return fmt.Sprintf("%s - %s", item.CompletedAt.Format("15:04"), item.Content)
}
