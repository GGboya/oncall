package cronjob

import (
	"oncall/internal/apiservice"
)

type Task interface {
	Name() string     // 任务名称，便于调试和日志
	Schedule() string // Cron 表达式，定义任务的调度规则
	Execute() error   // 任务执行逻辑
}

type WeeklyTask struct {
	apiService *apiservice.APIService
}

func NewWeeklyTask(apiService *apiservice.APIService) *WeeklyTask {
	return &WeeklyTask{apiService: apiService}
}

func (t *WeeklyTask) Name() string {
	return "week"
}

func (t *WeeklyTask) Schedule() string {
	return "0 0 * * 0" // 每周日 0 点
}

func (t *WeeklyTask) Execute() error {
	return executeTask(t.apiService, t.Name())
}

type DailyTask struct {
	apiService *apiservice.APIService
}

func NewDailyTask(apiService *apiservice.APIService) *DailyTask {
	return &DailyTask{apiService: apiService}
}

func (t *DailyTask) Name() string {
	return "day"
}

func (t *DailyTask) Schedule() string {
	return "0 0 * * *" // 每天 0 点
}

func (t *DailyTask) Execute() error {
	return executeTask(t.apiService, t.Name())
}
