package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"taskmanager/internal/model"
	"taskmanager/internal/port"
)

type Provider1 struct {
	URL  string
	Name string
}

// Provider1'den gelen verilerin yapısı
type TaskProvider1Response struct {
	ID                int64 `json:"id"`
	Value             int   `json:"value"`
	EstimatedDuration int   `json:"estimated_duration"`
}

func NewProvider1(url string) port.TaskProvider {
	return &Provider1{
		URL:  url,
		Name: "provider1",
	}
}

func (p *Provider1) GetTasks(ctx context.Context) ([]model.TaskProviderResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.URL, nil)
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
