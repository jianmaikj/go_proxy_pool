package main

import (
	"bufio"
	"encoding/json"
	"go_proxy_pool/cfg"
	"go_proxy_pool/utils"
	"log"
	"os"
	"sync"
	"time"
)

var wg2 sync.WaitGroup

// var mux sync.Mutex
var ch2 = make(chan int, 50)

// 是否抓取中
var run = false

func taskRun() {

	run = true
	defer func() {
		run = false
	}()

	count = 0
	log.Println("开始抓取代理...")
	for _, proxy := range cfg.ProxySources {
		wg2.Add(1)
		go fetchSource(proxy)
	}
	wg2.Wait()
	log.Printf("\r%s 代理抓取结束           \n", time.Now().Format(utils.StandardTimeFormat))

	//导出代理到文件
	export()

}

func fetchSource(sp cfg.ProxySource) {
	defer func() {
		wg2.Done()
		//log.Printf("%s 结束...",sp.Name)
	}()
	//log.Printf("%s 开始...", sp.Name)
	//urls := strings.Split(sp.GetApi(), ",")
	//var pis []module.ProxyIp

	res, err := cfg.ReqClient.GET(sp.GetApi()).Do()
	if err != nil {
		log.Println("代理抓取错误：" + err.Error())
		return
	}
	proxies := sp.Parse(res.Body)

	// 添加进代理池
	cfg.ProxyPool = append(cfg.ProxyPool, proxies...)
	//fmt.Println("res proxies>>", cfg.ProxyPool)
	//// 代理去重
	cfg.ProxyPool = utils.UniquePI(cfg.ProxyPool)
	// 过滤过期代理
	cfg.FilterIPs(cfg.ProxyPool)
	countAdd(len(proxies))
}

func export() {
	mux1.Lock()
	defer mux1.Unlock()
	//导出代理到文件
	err := os.Truncate("data.json", 0)
	if len(cfg.ProxyPool) == 0 {
		return
	}
	if err != nil {
		log.Printf("data.json清理失败：%s", err)
		return
	}
	file, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("data.json打开失败：%s", err)
		return
	}
	defer file.Close()

	data, err := json.Marshal(cfg.ProxyPool)
	if err != nil {
		log.Printf("代理json化失败：%s", err)
		return
	}
	buf := bufio.NewWriter(file)
	// 字节写入
	buf.Write(data)
	// 将缓冲中的数据写入
	err = buf.Flush()
	if err != nil {
		log.Println("代理json保存失败:", err)
	}
}

func InitTask() {
	//是否需要立即获取代理
	// 清理过期代理
	validCount := cfg.FilterIPs(cfg.ProxyPool)
	if validCount < cfg.Conf.Task.ProxyNum {
		//抓取代理
		log.Printf("当前代理数量 %d，不足 %d\n，开始抓取", validCount, cfg.Conf.Task.ProxyNum)
		taskRun()
	}

	//定时获取代理
	go func() {
		// 每 60 秒钟时执行一次
		getProxiesTimed := time.Duration(cfg.Conf.Task.GetProxiesTimed)
		ticker := time.NewTicker(getProxiesTimed * time.Second)
		for range ticker.C {
			// 清理过期代理
			validCount = cfg.FilterIPs(cfg.ProxyPool)
			if validCount < cfg.Conf.Task.ProxyNum {
				// 判断是否正在抓取或验证代理
				if !run && !verifyIS {
					log.Printf("定时抓取中>>>有效代理数量：%d，不足 %d\n，开始抓取", validCount, cfg.Conf.Task.ProxyNum)
					//抓取代理
					taskRun()
				}
			} else {
				//保存代理到本地
				export()
			}
		}
	}()

	////定时更换隧道IP
	//go func() {
	//	tunnelTime := time.Duration(conf.Config.TunnelTime)
	//	ticker := time.NewTicker(tunnelTime * time.Second)
	//	for range ticker.C {
	//		if len(ProxyPool) != 0 {
	//			httpsIp = getHttpsIp()
	//			httpIp = gethttpIp()
	//			socket5Ip = getSocket5Ip()
	//		}
	//	}
	//}()

	//// 验证代理存活情况
	//go func() {
	//	verifyTime := time.Duration(conf.Config.VerifyTime)
	//	ticker := time.NewTicker(verifyTime * time.Second)
	//	for range ticker.C {
	//		if !verifyIS && !run {
	//			VerifyProxy()
	//		}
	//	}
	//}()
}