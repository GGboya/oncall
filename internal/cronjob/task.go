package cronjob

import (
	"oncall/config"
	"oncall/internal/apiservice"
)

type Task interface {
	Name() string     // 任务名称，如 Fenlei、Fenxi 等
	Schedule() string // Cron 表达式，定义任务的调度规则
	Execute() error   // 任务执行逻辑
	Kind() string     // 用于区分任务类型，例如：week、day等
}

type WeeklyTask struct {
	apiService *apiservice.APIService
}

func NewWeeklyTask(apiService *apiservice.APIService) *WeeklyTask {
	return &WeeklyTask{apiService: apiService}
}

func (t *WeeklyTask) Name() string {
	return config.Fenlei
}

func (t *WeeklyTask) Schedule() string {
	return "0 0 * * 0" // 每周日 0 点
}

func (t *WeeklyTask) Execute() error {
	return ExecuteTask(t.apiService, t.Kind(), t.Name())
}

func (t *WeeklyTask) Kind() string {
	return "week"
}

type DailyTask struct {
	apiService *apiservice.APIService
}

func NewDailyTask(apiService *apiservice.APIService) *DailyTask {
	return &DailyTask{apiService: apiService}
}

func (t *DailyTask) Name() string {
	return config.Fenxi
}

func (t *DailyTask) Schedule() string {
	return "0 0 * * *" // 每天 0 点
}

func (t *DailyTask) Execute() error {
	return ExecuteTask(t.apiService, t.Kind(), t.Name())
}

func (t *DailyTask) Kind() string {
	return "day"
}
