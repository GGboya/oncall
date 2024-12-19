package cronjob

import (
	"fmt"
	"oncall/internal/apiservice"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// 创建文件夹辅助函数
func ensureDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		return os.MkdirAll(dirName, os.ModePerm)
	}
	return nil
}

func ExecuteTask(apiService *apiservice.APIService, kind string, taskName string) error {
	data, err := apiService.Oncall.FetchAllData(kind)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch data")
		return err
	}

	response, err := apiService.DeepSeek.SendDeepSeekRequest(data, taskName)
	if err != nil {
		logrus.WithError(err).Error("Failed to send request")
		return err
	}

	// 写入文件前检查目录
	outputDir := fmt.Sprintf("../output/%s", kind)
	if err := ensureDir(outputDir); err != nil {
		logrus.WithError(err).Errorf("Failed to create directory %s", outputDir)
		return err
	}

	outputFile := fmt.Sprintf("../output/%s/output_%s.txt", kind, time.Now().Format("20060102"))
	if err := WriteToFile(outputFile, response.Choices[0].Message.Content); err != nil {
		logrus.WithError(err).Error("Failed to write to file")
		return err
	}
	logrus.Infof("Response successfully written to %s\n", outputFile)
	return nil
}

// writeToFile 将内容写入指定的文件
func WriteToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
