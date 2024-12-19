package apiservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"oncall/config"
	"oncall/internal/httpclient"
)

const maxBackoff = 16

// doRequest 发送通用 HTTP 请求
func doRequest(client httpclient.HTTPClient, method, url string, body interface{}, config *config.Config) ([]byte, error) {
	var requestBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		requestBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	SetRequestHeaders(req, config)

	// 添加一个重试机制，指数回避
	var resp *http.Response
	backoff := 1

	for backoff <= maxBackoff {
		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", err)
		}
		if resp.StatusCode == 429 { // Too Many Requests
			waitDuration := time.Second * time.Duration(backoff) // Default wait time
			time.Sleep(waitDuration)                             // Wait before retrying
			backoff <<= 1
			continue
		}
		break
	}

	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func SetRequestHeaders(req *http.Request, config *config.Config) {
	// 设置 API Key 和用户名的 Header
	req.Header.Add("Api-Key", config.AccessKey)
	req.Header.Add("Api-Username", config.UserName)
	req.Header.Add("User-Agent", "Go-Discourse-Client/1.0")
}

// decodeUnicode 将Unicode编码的字符串转换为普通字符串
func decodeUnicode(str string) (string, error) {
	var result strings.Builder
	for i := 0; i < len(str); i++ {
		if str[i] == '\\' && i+1 < len(str) && str[i+1] == 'u' {
			i += 2 // 跳过 "\u"
			code := str[i : i+4]
			i += 4 // 跳过四位十六进制数
			runeValue, err := strconv.ParseUint(code, 16, 64)
			if err != nil {
				return "", err
			}
			result.WriteRune(rune(runeValue))
		} else {
			result.WriteByte(str[i])
		}
	}
	// 使用正则表达式去除HTML标签
	re := regexp.MustCompile(`<[^>]+>`)
	text := re.ReplaceAllString(result.String(), "")
	return text, nil
}
