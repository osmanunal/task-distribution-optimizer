package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"taskmanager/internal/model"
	"taskmanager/internal/port"
)

type Provider2 struct {
	URL  string
	Name string
}

type Provider2Task struct {
	ID     int64 `json:"id"`
	Zorluk int   `json:"zorluk"`
	Sure   int   `json:"sure"`
}

func NewProvider2(url string) port.TaskProvider {
	return &Provider2{
		URL:  url,
		Name: "provider2",
	}
}

func (p *Provider2) GetTasks(ctx context.Context) ([]model.TaskProviderResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("provider2 için istek oluştururken hata: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("provider2 API'sine erişirken hata: %v", err)
	}
	defer resp.Body.Close()

	var rawTasks []Provider2Task
	if err := json.NewDecoder(resp.Body).Decode(&rawTasks); err != nil {
		return nil, fmt.Errorf("provider2 verilerini parse ederken hata: %v", err)
	}

	var tasks []model.TaskProviderResponse
	for _, rawTask := range rawTasks {
		tasks = append(tasks, model.TaskProviderResponse{
			ExternalID: rawTask.ID,
			Difficulty: rawTask.Zorluk,
			Duration:   rawTask.Sure,
			Name:       p.Name,
		})
	}

	return tasks, nil
}
