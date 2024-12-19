package app

import (
	"oncall/config"
	"oncall/internal/cronjob"

	"github.com/sirupsen/logrus"
)

type Command interface {
	Execute(deps *AppDependencies) error
}

type DayCommand struct{}

func (c *DayCommand) Execute(deps *AppDependencies) error {
	logrus.Info("Executing day command...")
	return cronjob.ExecuteTask(deps.APIService, "day", config.Fenxi)
}
