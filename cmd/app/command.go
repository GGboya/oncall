package app

import (
	"fmt"
	"oncall/config"
	"oncall/internal/cronjob"
)

type Command interface {
	Execute(deps *AppDependencies) error
}

type DayCommand struct{}

func (c *DayCommand) Execute(deps *AppDependencies) error {
	fmt.Println("Executing day command...")
	return cronjob.ExecuteTask(deps.APIService, "day", config.Fenxi)
}
