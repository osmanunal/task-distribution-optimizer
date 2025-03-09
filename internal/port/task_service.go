package port

import (
	"context"
)

// TaskSyncService, task senkronizasyon servisi için port
type TaskSyncService interface {
	// SyncTasks, görevleri senkronize eder
	SyncTasks(ctx context.Context) error

	// TaskPlanner, görevleri çalışanlara planlar ve en optimal dağılımı hesaplar
	TaskPlanner(ctx context.Context) error
}
