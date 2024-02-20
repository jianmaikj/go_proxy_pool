package main

import (
	"go_proxy_pool/cfg"
)

func Init() {
	//获取配置文件
	cfg.InitConfigData()
	cfg.InitReqClient()
	cfg.InitSources()
	InitTask()
}