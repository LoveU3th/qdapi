package qdapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetProxyClient() *http.Client {
	address := "localhost:8888"
	conn, err := net.DialTimeout("tcp", address, time.Second*2)
	if err != nil {
		return http.DefaultClient
	}
	defer conn.Close()
	//for Charles
	proxyURL, err := url.Parse("http://" + address)
	if err != nil {
		log.Fatal("Invalid proxy URL:", err)
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
}
func GetInsecureClient() *http.Client {
	//for Charles
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 忽略证书验证
		},
	}
	return &http.Client{
		Transport: tr,
	}
}

func LoadConfigFromJSON(filename string) ([]QiDianApiConfig, error) {
	var config []QiDianApiConfig
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

// LoadConfigFromEnv 从环境变量加载配置
func LoadConfigFromEnv() ([]QiDianApiConfig, error) {
	var configs []QiDianApiConfig
	
	// 支持多账号，通过索引区分 (CONFIG_0_*, CONFIG_1_*, ...)
	for i := 0; ; i++ {
		prefix := ""
		if i == 0 {
			// 第一个账号可以不带索引，保持向后兼容
			prefix = "CONFIG_"
		} else {
			prefix = "CONFIG_" + strconv.Itoa(i) + "_"
		}
		
		qdInfo := os.Getenv(prefix + "QDINFO")
		if qdInfo == "" && i == 0 {
			// 如果第一个账号没有配置，尝试不带前缀的环境变量
			qdInfo = os.Getenv("QDINFO")
		}
		
		// 如果没有找到配置，停止查找
		if qdInfo == "" {
			break
		}
		
		sdkSign := os.Getenv(prefix + "SDKSIGN")
		if sdkSign == "" && i == 0 {
			sdkSign = os.Getenv("SDKSIGN")
		}
		
		ywKey := os.Getenv(prefix + "YWKEY")
		if ywKey == "" && i == 0 {
			ywKey = os.Getenv("YWKEY")
		}
		
		ywGuid := os.Getenv(prefix + "YWGUID")
		if ywGuid == "" && i == 0 {
			ywGuid = os.Getenv("YWGUID")
		}
		
		taskTypeStr := os.Getenv(prefix + "TASKTYPE")
		if taskTypeStr == "" && i == 0 {
			taskTypeStr = os.Getenv("TASKTYPE")
		}
		
		// 解析任务类型
		var taskTypes []TaskType
		if taskTypeStr != "" {
			taskStrs := strings.Split(taskTypeStr, ",")
			for _, taskStr := range taskStrs {
				taskStr = strings.TrimSpace(taskStr)
				if taskInt, err := strconv.Atoi(taskStr); err == nil {
					taskTypes = append(taskTypes, TaskType(taskInt))
				}
			}
		} else {
			// 默认任务类型
			taskTypes = []TaskType{1, 2, 3}
		}
		
		config := QiDianApiConfig{
			QdInfo:   qdInfo,
			SdkSign:  sdkSign,
			YwKey:    ywKey,
			YwGuid:   ywGuid,
			TaskType: taskTypes,
		}
		
		configs = append(configs, config)
	}
	
	if len(configs) == 0 {
		return nil, fmt.Errorf("no configuration found in environment variables")
	}
	
	return configs, nil
}

func SaveConfigToJSON(filename string, data interface{}) error {
	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建JSON编码器
	encoder := json.NewEncoder(file)

	// 设置格式缩进（可选）
	encoder.SetIndent("", "  ")

	// 执行编码
	return encoder.Encode(data)
}
