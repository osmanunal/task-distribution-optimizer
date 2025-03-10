package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"task-distribution-optimizer/internal/model"
	"task-distribution-optimizer/internal/port"
)

const task1URL = "https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/mock-one"

type TaskProvider1 struct {
	Name string
}

// Provider1'den gelen verilerin yapısı
type TaskProvider1Response struct {
	ID                int64 `json:"id"`
	Value             int   `json:"value"`
	EstimatedDuration int   `json:"estimated_duration"`
}

func NewProvider1() port.TaskProvider {
	return &TaskProvider1{
		Name: "provider1",
	}
}

func (p *TaskProvider1) GetTasks(ctx context.Context) ([]model.TaskProviderResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, task1URL, nil)
	if err != nil {
		return nil, fmt.Errorf("provider1 için istek oluştururken hata: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("provider1 API'sine erişirken hata: %v", err)
	}
	defer resp.Body.Close()

	var rawTasks []TaskProvider1Response
	if err := json.NewDecoder(resp.Body).Decode(&rawTasks); err != nil {
		return nil, fmt.Errorf("provider1 verilerini parse ederken hata: %v", err)
	}

	var tasks []model.TaskProviderResponse
	for _, rawTask := range rawTasks {
		tasks = append(tasks, model.TaskProviderResponse{
			ExternalID: rawTask.ID,
			Difficulty: rawTask.Value,
			Duration:   rawTask.EstimatedDuration,
			Name:       p.Name,
		})
	}

	return tasks, nil
}
