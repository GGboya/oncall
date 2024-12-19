package app

import (
	"fmt"
	"oncall/internal/cronjob"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCommand creates a *cobra.Command object with default parameters
func NewCommand(deps *AppDependencies) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "oncall",
		Short: "Launch an oncall server",
		Long:  "Launch an oncall server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(deps)
		},
	}

	subCommands := map[string]Command{
		"day": &DayCommand{},
	}

	for name, cmdLogic := range subCommands {
		subCmd := &cobra.Command{
			Use:   name,
			Short: fmt.Sprintf("Execute %s logic", name),
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmdLogic.Execute(deps)
			},
		}
		rootCmd.AddCommand(subCmd)
	}
	return rootCmd
}

func Run(deps *AppDependencies) error {
	cronScheduler := cronjob.NewCronScheduler()
	taskFactory := &DefaultTaskFactory{}

	// 动态注册任务
	tasks := taskFactory.CreateTasks(deps.APIService)
	for _, task := range tasks {
		if err := cronScheduler.RegisterTask(task); err != nil {
			logrus.WithError(err).Errorf("Failed to register task: %s", task.Name())
			return err
		}
	}

	// 启动调度器
	return cronScheduler.Start()
}
