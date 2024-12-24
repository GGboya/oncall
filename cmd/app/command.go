package app

import (
	"oncall/config"
	"oncall/internal/cronjob"

	"github.com/sirupsen/logrus"
)

type Command interface {
	Execute(deps *AppDependencies, q int) error
}

type DayCommand struct{}

func (c *DayCommand) Execute(deps *AppDependencies, q int) error {
	logrus.Info("Executing day command...")
	return cronjob.ExecuteTask(deps.APIService, "day", config.Fenxi, q)
}

type WeekCommand struct{}

func (c *WeekCommand) Execute(deps *AppDependencies, q int) error {
	logrus.Info("Executing week command...")
	return cronjob.ExecuteTask(deps.APIService, "week", config.Fenlei, q)
}

type QueryCommand struct{}

func (c *QueryCommand) Execute(deps *AppDependencies, q int) error {
	logrus.Info("Executing query command...")
	return cronjob.ExecuteTask(deps.APIService, "week", config.Fenxi, q)
}
