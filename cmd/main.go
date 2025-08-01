package main

import (
	"fmt"
	"github.com/pzx521521/qdapi"
	"github.com/pzx521521/qdapi/sign"
	"log"
	"net/http"
	"runtime"
	"sync"
)

func main() {
	// 优先尝试从 Docker 挂载路径读取配置
	configPaths := []string{
		"/app/config/config.json", // Docker 挂载路径
		"./config.json",           // 本地路径
	}
	
	var configs []qdapi.QiDianApiConfig
	var err error
	var configPath string
	
	for _, path := range configPaths {
		configs, err = qdapi.LoadConfigFromJSON(path)
		if err == nil {
			configPath = path
			break
		}
	}
	
	if err != nil {
		fmt.Printf("读取配置文件失败，请检查以下路径是否存在配置文件:\n")
		for _, path := range configPaths {
			fmt.Printf("- %s\n", path)
		}
		fmt.Printf("错误信息: %v\n", err)
		return
	}
	
	fmt.Printf("成功从 %s 加载了%d个账号配置\n", configPath, len(configs))
	
	var cli *http.Client
	if runtime.GOOS == "darwin" {
		//for charles
		cli = qdapi.GetProxyClient()
	} else {
		//for github action
		cli = qdapi.GetInsecureClient()
	}
	CheckInAndDoTaskMulti(cli, configs...)
}
func CheckInAndDoTaskMulti(cli *http.Client, configs ...qdapi.QiDianApiConfig) {
	var wg sync.WaitGroup
	for i, config := range configs {
		wg.Add(1)
		go func(index int, config qdapi.QiDianApiConfig) {
			CheckInAndDoTask(cli, config)
			wg.Done()
		}(i, config)
	}
	wg.Wait()
}

func CheckInAndDoTask(client *http.Client, config qdapi.QiDianApiConfig) {
	meta, err := sign.NewMeta(config.QdInfo, config.SdkSign)
	if err != nil {
		log.Printf("QdInfo或SdkSign解析错误:%v\n", err)
		return
	}

	log.Printf("%v\n", meta)
	api := qdapi.NewQiDianApi(meta, config.YwKey, config.YwGuid)
	api.Cli = client
	resp, err := api.CheckIn()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	log.Printf("%s:%v\n", api.TipName(), resp)
	err = qdapi.DoTask(api, config.TaskType...)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
