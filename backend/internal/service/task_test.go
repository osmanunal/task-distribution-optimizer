package service

import (
	"context"
	"task-distribution-optimizer/internal/model"
	"task-distribution-optimizer/internal/port"
	pkgmodel "task-distribution-optimizer/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TaskProviderMock, TaskProvider arayüzünün mock implementasyonu
type TaskProviderMock struct {
	mock.Mock
}

func (m *TaskProviderMock) GetTasks(ctx context.Context) ([]model.TaskProviderResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.TaskProviderResponse), args.Error(1)
}

// TaskRepositoryMock, TaskRepository arayüzünün mock implementasyonu
type TaskRepositoryMock struct {
	mock.Mock
}

func (m *TaskRepositoryMock) UpsertTasks(ctx context.Context, tasks []model.Task) error {
	args := m.Called(ctx, tasks)
	return args.Error(0)
}

func (m *TaskRepositoryMock) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Task), args.Error(1)
}

func (m *TaskRepositoryMock) MarkTasksAsProcessed(ctx context.Context, taskIDs []int64) error {
	args := m.Called(ctx, taskIDs)
	return args.Error(0)
}

// EmployeeRepositoryMock, EmployeeRepository arayüzünün mock implementasyonu
type EmployeeRepositoryMock struct {
	mock.Mock
}

func (m *EmployeeRepositoryMock) GetAllEmployees(ctx context.Context) ([]model.Employee, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Employee), args.Error(1)
}

func (m *EmployeeRepositoryMock) CreateEmployee(ctx context.Context, employee model.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

// TestSetup, test için gerekli tüm bileşenleri içeren bir yapı
type TestSetup struct {
	TaskProvider *TaskProviderMock
	TaskRepo     *TaskRepositoryMock
	EmployeeRepo *EmployeeRepositoryMock
	TaskService  port.TaskService
	Ctx          context.Context
}

// setupTest, test için gerekli tüm yapıları oluşturan yardımcı fonksiyon
func setupTest() TestSetup {
	mockTaskProvider := new(TaskProviderMock)
	mockTaskRepo := new(TaskRepositoryMock)
	mockEmployeeRepo := new(EmployeeRepositoryMock)

	taskService := NewTaskService(mockTaskProvider, mockTaskRepo, mockEmployeeRepo)
	ctx := context.Background()

	return TestSetup{
		TaskProvider: mockTaskProvider,
		TaskRepo:     mockTaskRepo,
		EmployeeRepo: mockEmployeeRepo,
		TaskService:  taskService,
		Ctx:          ctx,
	}
}

func TestTaskService_SyncTasks(t *testing.T) {
	setup := setupTest()

	// örnek veri
	sampleTasks := []model.TaskProviderResponse{
		{
			ExternalID: 1,
			Name:       "Test Task 1",
			Difficulty: 3,
			Duration:   5,
		},
		{
			ExternalID: 2,
			Name:       "Test Task 2",
			Difficulty: 2,
			Duration:   8,
		},
	}

	expectedModelTasks := []model.Task{
		{
			ExternalID: 1,
			Name:       "Test Task 1",
			Difficulty: 3,
			Duration:   5,
		},
		{
			ExternalID: 2,
			Name:       "Test Task 2",
			Difficulty: 2,
			Duration:   8,
		},
	}

	// Mock davranışlarını tanımla
	setup.TaskProvider.On("GetTasks", setup.Ctx).Return(sampleTasks, nil)
	setup.TaskRepo.On("UpsertTasks", setup.Ctx, mock.MatchedBy(func(tasks []model.Task) bool {
		// Tam eşleşme kontrolü yerine, alanların ayrı ayrı kontrolü
		if len(tasks) != len(expectedModelTasks) {
			return false
		}

		for i, task := range tasks {
			if task.ExternalID != expectedModelTasks[i].ExternalID ||
				task.Name != expectedModelTasks[i].Name ||
				task.Difficulty != expectedModelTasks[i].Difficulty ||
				task.Duration != expectedModelTasks[i].Duration {
				return false
			}
		}
		return true
	})).Return(nil)

	err := setup.TaskService.SyncTasks(setup.Ctx)

	assert.NoError(t, err)
	setup.TaskProvider.AssertExpectations(t)
	setup.TaskRepo.AssertExpectations(t)
}

func TestTaskService_SyncTasks_ProviderError(t *testing.T) {
	setup := setupTest()

	// Mock davranışlarını tanımla - provider hata döndürecek
	providerErr := assert.AnError
	setup.TaskProvider.On("GetTasks", setup.Ctx).Return([]model.TaskProviderResponse{}, providerErr)

	err := setup.TaskService.SyncTasks(setup.Ctx)

	assert.Error(t, err)
	assert.Equal(t, providerErr, err)
	setup.TaskProvider.AssertExpectations(t)
	// UpsertTasks çağrılmamalı, çünkü provider hatası var
	setup.TaskRepo.AssertNotCalled(t, "UpsertTasks")
}

func TestTaskService_SyncTasks_RepositoryError(t *testing.T) {
	setup := setupTest()

	sampleTasks := []model.TaskProviderResponse{
		{
			ExternalID: 1,
			Name:       "Test Task 1",
			Difficulty: 3,
			Duration:   5,
		},
	}

	// Mock davranışlarını tanımla - repo hata döndürecek
	repoErr := assert.AnError
	setup.TaskProvider.On("GetTasks", setup.Ctx).Return(sampleTasks, nil)
	setup.TaskRepo.On("UpsertTasks", setup.Ctx, mock.Anything).Return(repoErr)

	err := setup.TaskService.SyncTasks(setup.Ctx)

	// Sonuçları doğrula - hata dönmeli
	assert.Error(t, err)
	assert.Equal(t, repoErr, err)
	setup.TaskProvider.AssertExpectations(t)
	setup.TaskRepo.AssertExpectations(t)
}

func TestTaskService_TaskPlanner(t *testing.T) {
	setup := setupTest()

	tasks := []model.Task{
		{
			BaseModel:  pkgmodel.BaseModel{ID: 1},
			ExternalID: 101,
			Name:       "provider1",
			Difficulty: 5,
			Duration:   10,
		},
		{
			BaseModel:  pkgmodel.BaseModel{ID: 2},
			ExternalID: 102,
			Name:       "provider1",
			Difficulty: 2,
			Duration:   8,
		},
	}

	employees := []model.Employee{
		{
			BaseModel:  pkgmodel.BaseModel{ID: 1},
			Name:       "provider1",
			Difficulty: 5,
		},
		{
			BaseModel:  pkgmodel.BaseModel{ID: 2},
			Name:       "provider1",
			Difficulty: 2,
		},
	}

	// Mock davranışlarını tanımla
	setup.TaskRepo.On("GetAllTasks", setup.Ctx).Return(tasks, nil)
	setup.EmployeeRepo.On("GetAllEmployees", setup.Ctx).Return(employees, nil)

	result, err := setup.TaskService.TaskPlanner(setup.Ctx)

	// Sonuçları doğrula
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// Toplam hafta sayısı kontrolü
	assert.Greater(t, result.TotalWeeks, 0)

	// Çalışan iş yükü kontrolleri
	assert.Equal(t, len(employees), len(result.Workloads))

	// İş yüklerinin doğru dağıtıldığını kontrol et
	totalAssignments := 0
	for _, workload := range result.Workloads {
		// Her çalışana atanan görevleri say
		totalAssignments += len(workload.Assignments)

		// Her çalışanın haftalık planının doğru hesaplandığını kontrol et
		if len(workload.Assignments) > 0 {
			totalHoursInPlan := 0
			for _, week := range workload.WeeklyPlan {
				totalHoursInPlan += week.Hours
			}
			assert.Equal(t, workload.TotalHours, totalHoursInPlan,
				"Haftalık plan saatleri toplamı, toplam iş saatine eşit olmalı")
		}
	}

	// Tüm görevlerin atandığını kontrol et
	assert.Equal(t, len(tasks), totalAssignments, "Tüm görevler çalışanlara atanmalı")

	// Mock beklentilerini doğrula
	setup.TaskRepo.AssertExpectations(t)
	setup.EmployeeRepo.AssertExpectations(t)
}

func TestTaskService_TaskPlanner_NoTasks(t *testing.T) {
	setup := setupTest()

	tasks := []model.Task{}
	employees := []model.Employee{
		{
			BaseModel:  pkgmodel.BaseModel{ID: 1},
			Name:       "Uzman Çalışan",
			Difficulty: 5,
		},
	}

	// Mock davranışlarını tanımla
	setup.TaskRepo.On("GetAllTasks", setup.Ctx).Return(tasks, nil)
	setup.EmployeeRepo.On("GetAllEmployees", setup.Ctx).Return(employees, nil)

	result, err := setup.TaskService.TaskPlanner(setup.Ctx)

	// Sonuçları doğrula - tasks boş olduğu için bir uyarı değil ama sonuç dönmeli
	assert.NoError(t, err) // Hata beklenmiyorsa NoError kullanılır
	assert.Empty(t, result.Workloads, "Boş task listesi için workloads boş olmalı")
	assert.Equal(t, 0, result.TotalWeeks, "Boş task listesi için toplam hafta sayısı 0 olmalı")
	setup.TaskRepo.AssertExpectations(t)
	setup.EmployeeRepo.AssertExpectations(t)
}

func TestTaskService_TaskPlanner_NoEmployees(t *testing.T) {
	setup := setupTest()

	tasks := []model.Task{
		{
			BaseModel:  pkgmodel.BaseModel{ID: 1},
			Name:       "Test Task",
			Difficulty: 3,
			Duration:   5,
		},
	}
	employees := []model.Employee{}

	// Mock davranışlarını tanımla
	setup.TaskRepo.On("GetAllTasks", setup.Ctx).Return(tasks, nil)
	setup.EmployeeRepo.On("GetAllEmployees", setup.Ctx).Return(employees, nil)

	result, err := setup.TaskService.TaskPlanner(setup.Ctx)

	// Sonuçları doğrula - employees boş olduğu için bir uyarı değil ama sonuç dönmeli
	assert.NoError(t, err) // Hata beklenmiyorsa NoError kullanılır
	assert.Empty(t, result.Workloads, "Boş çalışan listesi için workloads boş olmalı")
	assert.Equal(t, 0, result.TotalWeeks, "Boş çalışan listesi için toplam hafta sayısı 0 olmalı")
	setup.TaskRepo.AssertExpectations(t)
	setup.EmployeeRepo.AssertExpectations(t)
}

func TestTaskService_TaskPlanner_RepositoryErrors(t *testing.T) {
	// 1. Task repository hatası
	t.Run("TaskRepository error", func(t *testing.T) {
		setup := setupTest()

		repoErr := assert.AnError
		setup.TaskRepo.On("GetAllTasks", setup.Ctx).Return([]model.Task{}, repoErr)

		_, err := setup.TaskService.TaskPlanner(setup.Ctx)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		setup.TaskRepo.AssertExpectations(t)
		// Employee repo çağrılmamalı çünkü task repo hatası var
		setup.EmployeeRepo.AssertNotCalled(t, "GetAllEmployees")
	})

	// 2. Employee repository hatası
	t.Run("EmployeeRepository error", func(t *testing.T) {
		setup := setupTest()

		tasks := []model.Task{
			{
				BaseModel:  pkgmodel.BaseModel{ID: 1},
				Name:       "Test Task",
				Difficulty: 3,
				Duration:   5,
			},
		}

		repoErr := assert.AnError
		setup.TaskRepo.On("GetAllTasks", setup.Ctx).Return(tasks, nil)
		setup.EmployeeRepo.On("GetAllEmployees", setup.Ctx).Return([]model.Employee{}, repoErr)

		_, err := setup.TaskService.TaskPlanner(setup.Ctx)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		setup.TaskRepo.AssertExpectations(t)
		setup.EmployeeRepo.AssertExpectations(t)
	})
}

func TestCalculateWeeklyPlan(t *testing.T) {
	tests := []struct {
		name             string
		totalHours       int
		weeklyWorkHours  int
		expectedWeeks    int
		expectedLastWeek int
	}{
		{
			name:             "İş yükü bir haftadan az",
			totalHours:       30,
			weeklyWorkHours:  45,
			expectedWeeks:    1,
			expectedLastWeek: 30,
		},
		{
			name:             "İş yükü tam bir hafta",
			totalHours:       45,
			weeklyWorkHours:  45,
			expectedWeeks:    1,
			expectedLastWeek: 45,
		},
		{
			name:             "İş yükü bir haftadan fazla",
			totalHours:       100,
			weeklyWorkHours:  45,
			expectedWeeks:    3, // 45 + 45 + 10
			expectedLastWeek: 10,
		},
		{
			name:             "İş yükü tam iki hafta",
			totalHours:       90,
			weeklyWorkHours:  45,
			expectedWeeks:    2,
			expectedLastWeek: 45,
		},
		{
			name:             "İş yükü sıfır",
			totalHours:       0,
			weeklyWorkHours:  45,
			expectedWeeks:    0,
			expectedLastWeek: 0,
		},
		{
			name:             "Farklı haftalık limit",
			totalHours:       62,
			weeklyWorkHours:  20,
			expectedWeeks:    4, // 20 + 20 + 20 + 2
			expectedLastWeek: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateWeeklyPlan(tc.totalHours, tc.weeklyWorkHours)

			// Hafta sayısını kontrol et
			assert.Equal(t, tc.expectedWeeks, len(result), "Hafta sayısı beklenen değere eşit olmalı")

			// Toplam saatleri kontrol et
			totalHoursInPlan := 0
			for _, week := range result {
				totalHoursInPlan += week.Hours
			}
			assert.Equal(t, tc.totalHours, totalHoursInPlan, "Plandaki toplam saat, girdi olan toplam saate eşit olmalı")

			// Son haftadaki saati kontrol et (varsa)
			if len(result) > 0 {
				assert.Equal(t, tc.expectedLastWeek, result[len(result)-1].Hours, "Son haftanın saati beklenen değere eşit olmalı")
			}

			// Hafta numaralarının doğru sıralandığını kontrol et
			for i, week := range result {
				assert.Equal(t, i+1, week.WeekNumber, "Hafta numarası doğru olmalı")
			}

			// Tam haftalarda maksimum limiti aşmamalı
			for i, week := range result {
				if i < len(result)-1 { // Son hafta hariç diğer haftalar
					assert.Equal(t, tc.weeklyWorkHours, week.Hours, "Ara haftalardaki saat, haftalık limit kadar olmalı")
				}
			}
		})
	}
}
