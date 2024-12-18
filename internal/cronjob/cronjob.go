package cronjob

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type CronScheduler struct {
	cron  *cron.Cron
	tasks []Task
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		cron:  cron.New(),
		tasks: []Task{},
	}
}

// 注册任务
func (s *CronScheduler) RegisterTask(task Task) error {
	_, err := s.cron.AddFunc(task.Schedule(), func() {
		logrus.Infof("Running task: %s", task.Name())
		if err := task.Execute(); err != nil {
			logrus.WithError(err).Errorf("Task %s failed", task.Name())
		} else {
			logrus.Infof("Task %s completed successfully", task.Name())
		}
	})
	if err != nil {
		return err
	}
	s.tasks = append(s.tasks, task) // 保存任务信息
	return nil
}

// 启动调度器
func (s *CronScheduler) Start() error {
	s.cron.Start()
	logrus.Info("Scheduler started, wating for tasks")
	select {} // 阻塞主线程
}
