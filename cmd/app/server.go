package app

import (
	"oncall/config"
	"oncall/internal/apiservice"
	"oncall/internal/cronjob"
	"oncall/internal/httpclient"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCommand creates a *cobra.Command object with default parameters
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Launch a oncall server",
		Long:  "Launch a oncall server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.NewConfig()
			if err != nil {
				return err
			}
			return Run(cfg)
		},
	}
	return cmd
}

func Run(cfg *config.Config) error {
	httpClient := httpclient.NewHTTPClient()
	apiService := apiservice.NewAPIService(httpClient, cfg)
	cronScheduler := cronjob.NewCronScheduler()

	// 动态注册任务
	tasks := []cronjob.Task{
		cronjob.NewWeeklyTask(apiService),
		cronjob.NewDailyTask(apiService),
	}

	for _, task := range tasks {
		if err := cronScheduler.RegisterTask(task); err != nil {
			logrus.WithError(err).Errorf("Failed to register task: %s", task.Name())
			return err
		}
	}

	// 启动调度器
	return cronScheduler.Start()
}
