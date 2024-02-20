package cfg

import (
	"encoding/json"
	"github.com/jianmaikj/req"
	"go_proxy_pool/module"
	"go_proxy_pool/sources"
	"go_proxy_pool/utils"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type ProxySource interface {
	Parse(resBody []byte) (proxies []*module.ProxyIp)
	GetApi() string
}

var ProxyPool []*module.ProxyIp
var ProxySources []ProxySource

var Conf *Config

var ReqClient *req.Client

type Config struct {
	SystemProxy SystemProxy `yaml:"systemProxy" json:"systemProxy"` //系统代理,用于抓取代理
	App         app         `yaml:"app" json:"app"`
	Task        task        `yaml:"task" json:"task"`
	Rdb         rdb         `yaml:"rdb" json:"rdb"`
}

type rdb struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
}

type task struct {
	ProxyNum        int    `yaml:"proxyNum" json:"proxyNum"`
	VerifyTime      int    `yaml:"verifyTime" json:"verifyTime"`
	VerifyUrl       string `yaml:"verifyUrl" json:"verifyUrl"`
	ThreadNum       int    `yaml:"threadNum" json:"threadNum"`
	GetProxiesTimed int    `yaml:"getProxiesTimed" json:"getProxiesTimed"`
}

type app struct {
	Ip   string `yaml:"ip" json:"ip"`
	Port string `yaml:"port" json:"port"`
}

type SystemProxy struct {
	Host string `yaml:"host" json:"host"`
	Port string `yaml:"port" json:"port"`
}

func FilterIPs(ips []*module.ProxyIp) int {
	var result []*module.ProxyIp
	validCount := 0
	for _, px := range ips {
		if !utils.IsExpiredProxy(px) {
			result = append(result, px)
		}
		if utils.IsValidProxy(px) {
			validCount++
		}
	}
	ProxyPool = result
	return validCount
}

// 读取配置文件
func InitConfigData() {
	//导入配置文件
	yamlFile, err := os.ReadFile("cfg/config.yml")
	if err != nil {
		log.Println("配置文件打开错误：" + err.Error())
		err.Error()
		return
	}
	//将配置文件读取到结构体中
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Println("配置文件解析错误：" + err.Error())
		err.Error()
		return
	}
	//导入代理缓存
	file, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("代理json文件打开错误：" + err.Error())
		err.Error()
		return
	}
	defer file.Close()
	all, err := io.ReadAll(file)
	if err != nil {
		log.Println("代理json解析错误：" + err.Error())
		return
	}
	if len(all) == 0 {
		return
	}
	err = json.Unmarshal(all, &ProxyPool)
	if err != nil {
		log.Println("代理json解析错误：" + err.Error())
		return
	}

}

func InitReqClient() {
	ReqClient = req.NewClient(
		&req.ClientConfig{
			Timeout: 5,
		})
}

func InitSources() {
	ProxySources = []ProxySource{
		//	&DMDaili{Api: "http://get.dmdaili.com/dmgetip.asp?apikey=ee350ccc&pwd=d3024c8aff338167787da2d2275c93b9&getnum=200&httptype=1&geshi=2&fenge=1&fengefu=&Contenttype=1&operate=all&setcity=all&provin=fujian", Type: "http", Province: "fujian"}, //福建地区

		//&sources.DMDaili{Api: "http://get.dmdaili.com/dmgetip.asp?apikey=ee350ccc&pwd=d3024c8aff338167787da2d2275c93b9&getnum=200&httptype=1&geshi=2&fenge=1&fengefu=&Contenttype=1&operate=all&setcity=all&provin=zhejiang", Type: "http", Province: "zhejiang"}, //浙江地区
		//
		//&sources.DMDaili{Api: "http://get.dmdaili.com/dmgetip.asp?apikey=ee350ccc&pwd=d3024c8aff338167787da2d2275c93b9&getnum=200&httptype=1&geshi=2&fenge=1&fengefu=&Contenttype=1&operate=all&setcity=all&provin=guangdong", Type: "http", Province: "guangdong"}, //广东地区
		//
		//&sources.WuyiDaili{Api: "http://bapi.51daili.com/unlimitedip/getip?linePoolIndex=1&packid=17&time=5&qty=200&port=2&format=json&field=regioncode,expiretime&pid=fc38e6a977585bfd31d5980b64b3ed09&usertype=17&uid=41920&regionCodeList=360000,440000,331200,330200&skey=autoaddwhiteip&dt=2&notRegionCode=440400,361100,360800", Type: "socks5", AreaCodeList: "360000,440000"}, //广东/江西

		//&sources.WuyiDaili{Api: "http://bapi.51daili.com/traffic/getip?linePoolIndex=1&packid=12&time=1&qty=100&port=2&format=json&field=expiretime,regioncode,isptype&dt=1&providerIds=1,5&regionCode=440000&usertype=17&uid=41920&notRegionCode=440400", Type: "socks5", AreaCodeList: "360000,440000"}, //广东
		&sources.WuyiDaili{Api: "http://bapi.51daili.com/traffic/getip?linePoolIndex=1&packid=12&time=1&port=2&format=json&field=expiretime,regioncode,isptype&providerIds=1,5&usertype=17&uid=41920&qty=200&dt=3&notRegionCode=110000,120000,130000,140000,150000,210000,220000,310000,370000,410000,420000,500000,510000,610000,620000,630000,640000,650000,230000,540000", Type: "socks5", AreaCodeList: "360000,440000"}, //广东

		//&sources.PinYi{Api: "http://zltiqu.pyhttp.taolop.com/getip?count=100&neek=13938&type=2&yys=0&port=11&sb=&mr=2&sep=0&ts=1&ys=1&cs=1&regions=440000", Type: "socks5"}, //广东/江西/湖南
	}
}