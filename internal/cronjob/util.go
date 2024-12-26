package cronjob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"oncall/internal/apiservice"
	"os"
	"strings"
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

func ExecuteTask(apiService *apiservice.APIService, interval string, taskName string, q int) error {
	logrus.Info("Executing task")
	data, err := apiService.Oncall.FetchAllData(interval, taskName, q)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch data")
		return err
	}

	logrus.Info("Fetched data, preparing to send DeepSeekRequest")
	response, err := apiService.DeepSeek.SendDeepSeekRequest(data, taskName)
	if err != nil {
		logrus.WithError(err).Error("Failed to send request")
		return err
	}

	// 将内容中的数据进行分类，存入到 es 中
	SendToEs(apiService, response.Choices[0].Message.Content)

	logrus.Infof("Response: %s\n", response.Choices[0].Message.Content)
	// 写入文件前检查目录
	outputDir := fmt.Sprintf("%s/output/%s", os.Getenv("OUTPUT_DIR"), interval)
	if err := ensureDir(outputDir); err != nil {
		logrus.WithError(err).Errorf("Failed to create directory %s", outputDir)
		return err
	}

	outputFile := fmt.Sprintf("%s/output/%s/output_%s_%d.txt", os.Getenv("OUTPUT_DIR"), interval, time.Now().Format("20060102"), time.Now().UnixNano()/int64(time.Millisecond))
	if err := WriteToFile(outputFile, response.Choices[0].Message.Content); err != nil {
		logrus.WithError(err).Error("Failed to write to file")
		return err
	}
	logrus.Infof("Response successfully written to %s\n", outputFile)
	return nil
}

func SendToEs(apiService *apiservice.APIService, content string) error {
	// 先整理出所有的数据
	lines := strings.Split(content, "\n")
	var questions []UserQuestion
	var currentQuestion UserQuestion

	for _, line := range lines {
		line = strings.TrimSpace(line) // 去掉空白字符
		if strings.Contains(line, "**名字**") {
			currentQuestion.Name = strings.TrimPrefix(strings.Split(line, "：")[1], " ")
		} else if strings.Contains(line, "**项目(部门)**") {
			currentQuestion.Project = strings.TrimPrefix(strings.Split(line, "：")[1], " ")
		} else if strings.Contains(line, "**配额组(组)**") {
			currentQuestion.QuotaGroup = strings.TrimPrefix(strings.Split(line, "：")[1], " ")
		} else if strings.Contains(line, "**问题类型**") {
			currentQuestion.IssueType = strings.TrimPrefix(strings.Split(line, "：")[1], " ")
		} else if strings.Contains(line, "**帖子链接**") {
			currentQuestion.PostLink = strings.TrimPrefix(strings.Split(line, "：")[1], " ")
		} else if strings.Contains(line, "**问题归因**") {
			currentQuestion.Cause = strings.TrimPrefix(strings.Split(line, "：")[1], " ")
		} else if strings.Contains(line, "**问题描述**") {
			currentQuestion.Description = strings.TrimPrefix(strings.Split(line, "：")[1], " ")
		} else if strings.Contains(line, "**分析**") {
			currentQuestion.Analysis = strings.TrimPrefix(strings.Split(line, "：")[1], " ")
		} else if strings.Contains(line, "创建时间") {
			tempStr := strings.TrimPrefix(strings.Split(line, "：")[1], " ")
			createdAt, err := time.Parse("2006-01-02 15:04:05.000 -0700 MST", tempStr)
			if err != nil {
				logrus.WithError(err).Errorf("Failed to parse time: %s", tempStr)
				currentQuestion.Timestamp = time.Now()
			}
			if err == nil {
				currentQuestion.Timestamp = createdAt
			}
		}

		// 判断是否到达问题结尾
		if line == "" && currentQuestion.Name != "" {
			questions = append(questions, currentQuestion)
			currentQuestion = UserQuestion{} // 重置当前问题
		}
	}

	// 如果最后一个问题没有结尾行，也加入结果
	if currentQuestion.Name != "" {
		questions = append(questions, currentQuestion)
	}

	logrus.Info("Questions parsed successfully.")

	for _, doc := range questions {
		// 将文档转为 JSON 格式
		docJSON, err := json.Marshal(doc)
		if err != nil {
			return err
		}

		// 将文档插入到索引中
		insertReq := bytes.NewReader(docJSON)
		insertRes, err := apiService.Oncall.Es.Index("brainpp-oncall-questions", insertReq)
		if err != nil {
			return err
		}
		defer insertRes.Body.Close()

		if insertRes.IsError() {
			return err
		}
	}
	logrus.Info("Questions inserted successfully.")
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

type UserQuestion struct {
	Name        string    `json:"name"`
	Project     string    `json:"project"`
	QuotaGroup  string    `json:"quota_group"`
	PostLink    string    `json:"post_link"`
	IssueType   string    `json:"issue_type"`
	Cause       string    `json:"cause"`
	Description string    `json:"description"`
	Analysis    string    `json:"analysis"`
	Timestamp   time.Time `json:"@timestamp"`
}
