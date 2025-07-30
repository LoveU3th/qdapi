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
	// 优先尝试从环境变量加载配置
	configs, err := qdapi.LoadConfigFromEnv()
	if err != nil {
		// 如果环境变量加载失败，尝试从文件加载
		fmt.Printf("从环境变量加载配置失败，尝试从config.json加载: %v\n", err)
		configs, err = qdapi.LoadConfigFromJSON("./config.json")
		if err != nil {
			fmt.Printf("%s\n", "读取配置文件失败,请检查抓包数据或环境变量配置")
			return
		}
		fmt.Printf("成功从config.json加载了%d个账号配置\n", len(configs))
	} else {
		fmt.Printf("成功从环境变量加载了%d个账号配置\n", len(configs))
	}
	
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
