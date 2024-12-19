package app

import (
	"oncall/internal/apiservice"
	"oncall/internal/cronjob"
)

type TaskFactory interface {
	CreateTasks(apiService *apiservice.APIService) []cronjob.Task
}

type DefaultTaskFactory struct{}

func (f *DefaultTaskFactory) CreateTasks(apiService *apiservice.APIService) []cronjob.Task {
	return []cronjob.Task{
		cronjob.NewWeeklyTask(apiService),
		cronjob.NewDailyTask(apiService),
	}
}
