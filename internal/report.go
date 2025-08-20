package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"git.mmeiblog.cn/mei/FatAgent/configs"
	"git.mmeiblog.cn/mei/FatAgent/pkg"
)

const (
	reportPath  = "/receive/airConditioner"
	httpTimeout = 30 * time.Second
)

type airConditionerBodyMessage struct {
	AgentID   int          `json:"agentID"`
	Timestamp int64        `json:"timestamp"`
	Message   pkg.ACStatus `json:"message"`
}

var httpClient = &http.Client{
	Timeout: httpTimeout,
}

func Report() {
	now := time.Now().Unix()

	ac, err := pkg.NewACController(configs.AC_PORT)
	if err != nil {
		log.Printf("错误：创建控制器失败: %v", err)
		return
	}

	ACStatus, err := ac.GetAllStatus()
	if err != nil {
		log.Printf("错误：查询状态失败: %v", err)
		return
	}

	payload := airConditionerBodyMessage{
		AgentID:   configs.AGENT_ID,
		Timestamp: now,
		Message:   *ACStatus,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("错误：JSON序列化失败: %v", err)
		return
	}

	url := configs.REPORT_URL + reportPath
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("错误：创建HTTP请求失败: %v", err)
		return
	}

	req.Header.Add("x-api-key", configs.X_API_KEY)
	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("错误：发送HTTP请求失败: %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("警告：HTTP响应状态码非200，状态码: %d", res.StatusCode)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("错误：读取HTTP响应失败: %v", err)
		return
	}

	log.Printf("上报成功，响应内容:%s", body)
}
