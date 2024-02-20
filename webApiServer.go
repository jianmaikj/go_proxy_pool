package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_proxy_pool/cfg"
	"go_proxy_pool/module"
	"go_proxy_pool/utils"
	"log"
	"strconv"
)

type Home struct {
	TunnelProxy map[string]string `yaml:"tunnelProxy" json:"tunnelProxy"`
	Sum         int               `yaml:"sum" json:"sum"`
	Type        map[string]int    `yaml:"type" json:"type"`
	//Anonymity   map[string]int    `yaml:"anonymity" json:"anonymity"`
	Country map[string]int `yaml:"country" json:"country"`
	Source  map[string]int `yaml:"source" json:"source"`
}

var lastIP *module.ProxyIp

func Run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	//首页
	r.GET("/", index)

	//查询
	r.GET("/get", get)

	//删除
	r.GET("/delete", delete)

	//抓取代理
	r.GET("/spider", spiderUp)

	////更换隧道代理IP
	//r.GET("/tunnelUpdate", tunnelUpdate)

	log.Printf("webApi启动 - 监听IP端口 -> %s\n", cfg.Conf.App.Ip+":"+cfg.Conf.App.Port)

	r.Run(cfg.Conf.App.Ip + ":" + cfg.Conf.App.Port)

}
func index(c *gin.Context) {
	home := Home{
		Sum:  len(cfg.ProxyPool),
		Type: make(map[string]int),
		//Anonymity: make(map[string]int),
		Country: make(map[string]int),
		Source:  make(map[string]int),
		//TunnelProxy: make(map[string]string)
	}
	for i := range cfg.ProxyPool {
		home.Type[cfg.ProxyPool[i].Type] += 1
		//home.Anonymity[cfg.ProxyPool[i].Anonymity] += 1
		home.Country[cfg.ProxyPool[i].Country] += 1
		home.Source[cfg.ProxyPool[i].Source] += 1
	}
	//home.TunnelProxy["HTTP"] = httpIp
	//home.TunnelProxy["HTTPS"] = httpsIp
	//home.TunnelProxy["SOCKET5"] = socket5Ip
	jsonByte, _ := json.Marshal(&home)
	jsonStr := string(jsonByte)
	c.String(200, jsonStr)
}

func get(c *gin.Context) {
	if len(cfg.ProxyPool) == 0 {
		c.String(200, fmt.Sprintf("{\"code\": 200, \"msg\": \"代理池是空的\"}"))
		return
	}
	var prs []*module.ProxyIp
	var jsonByte []byte
	ty := c.DefaultQuery("type", "all")
	//an := c.DefaultQuery("anonymity", "all")
	country := c.DefaultQuery("country", "all")
	province := c.DefaultQuery("province", "all")
	so := c.DefaultQuery("source", "all")
	co := c.DefaultQuery("count", "1")
	//target := c.DefaultQuery("count", "1")
	//proxyPool := deepCopy(cfg.ProxyPool)
	for _, v := range cfg.ProxyPool {
		if (v.Type == ty || ty == "all") && (v.Country == country || country == "all") && (v.Source == so || so == "all") && (v.Province == province || province == "all") {
			prs = append(prs, v)
		}
	}
	if co == "all" {
		jsonByte, _ = json.Marshal(prs)
	} else if co == "1" {
		fmt.Println("lastIP>>", lastIP)
		prs = utils.SortedIPsByReqNum(prs)
		choice := prs[0]
		for _, v := range prs {
			if lastIP == nil || v.Addr != lastIP.Addr {
				if utils.IsValidProxy(choice) {
					choice = v
					break
				}

			}
		}
		choice.RequestNum += 1
		lastIP = choice
		jsonByte, _ = json.Marshal([]module.ProxyIp{*choice})

	} else {
		coInt, err := strconv.Atoi(co)
		if err != nil {
			c.String(200, fmt.Sprintf("{\"code\": 200, \"msg\": \"count参数错误\"}"))
			return
		}
		prs = utils.SortedIPsByReqNum(prs)
		jsonByte, _ = json.Marshal(prs[:coInt])

	}
	jsonStr := string(jsonByte)
	c.String(200, jsonStr)
}

func delete(c *gin.Context) {
	if len(cfg.ProxyPool) == 0 {
		c.String(200, fmt.Sprintf("{\"code\": 200, \"msg\": \"代理池是空的\"}"))
		return
	}
	ip := c.Query("ip")
	port := c.Query("port")
	i := delIp(ip + ":" + port)
	c.String(200, fmt.Sprintf("{\"code\": 200, \"count\": %d}", i))
}

//func verify(c *gin.Context) {
//	if verifyIS {
//		c.String(200, fmt.Sprintf("{\"code\": 200, \"msg\": \"验证中\"}"))
//	} else if run {
//		c.String(200, fmt.Sprintf("{\"code\": 200, \"msg\": \"代理抓取中，请稍后再来验证\"}"))
//	} else {
//		go VerifyProxy()
//		c.String(200, fmt.Sprintf("{\"code\": 200, \"msg\": \"开始验证代理\"}"))
//	}
//}

func spiderUp(c *gin.Context) {
	if run {
		c.String(200, fmt.Sprintf("{\"code\": 200, \"msg\": \"抓取中\"}"))
	} else if verifyIS {
		c.String(200, fmt.Sprintf("{\"code\": 200, \"msg\": \"代理验证中，请稍后再来抓取\"}"))
	} else {
		go taskRun()
		c.String(200, fmt.Sprintf("{\"code\": 200, \"msg\": \"开始抓取代理IP\"}"))
	}
}

//func tunnelUpdate(c *gin.Context) {
//	if len(ProxyPool) == 0 {
//		c.String(200, fmt.Sprintf("{\"code\": 200, \"msg\": \"代理池是空的\"}"))
//	}
//	httpsIp = getHttpsIp()
//	httpIp = gethttpIp()
//	socket5Ip = getSocket5Ip()
//	c.String(200, fmt.Sprintf("{\"code\": 200, \"HTTP\": \"%s\",\"HTTPS\": \"%s\",\"SOCKET5\": \"%s\" }", httpIp, httpsIp, socket5Ip))
//}

func delIp(addr string) int {
	lock.Lock()
	defer lock.Unlock()
	var in int
	for i, v := range cfg.ProxyPool {
		if v.Addr == addr {
			in++
			if i+1 < len(cfg.ProxyPool) {
				cfg.ProxyPool = append(cfg.ProxyPool[:i], cfg.ProxyPool[i+1:]...)
			} else {
				cfg.ProxyPool = cfg.ProxyPool[:i]
			}
		}
	}
	return in
}