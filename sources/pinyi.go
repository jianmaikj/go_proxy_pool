package sources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"go_proxy_pool/module"
	"go_proxy_pool/utils"
	"time"
)

// 51daili
type PinYi struct {
	Api  string
	Type string //http/socks5
	//AreaCodeList string //省份
	//Verify   bool
}

//	func (s *PinYi) GetProvince(areaCode string) string {
//		switch areaCode[:3] {
//		case "440":
//			return "guangdong"
//		case "360":
//			return "jiangxi"
//		case "430":
//			return "hunan"
//		case "350":
//			return "fujian"
//		case "330":
//			return "zhejiang"
//		case "310":
//			return "shanghai"
//		}
//		return ""
//	}
func (s *PinYi) Parse(resBody []byte) (proxies []*module.ProxyIp) {
	result := gjson.ParseBytes(resBody)
	if result.Get("code").Int() != 0 {
		return
	}
	for _, proxy := range result.Get("data").Array() {
		expiredTimeStr, err := utils.ParseTimeStr("2006-01-02 15:04:05", proxy.Get("expire_time").String())
		//expiredTimestamp := expiredTimeStr.Unix()
		if err != nil {
			continue
		}
		proxies = append(proxies, &module.ProxyIp{
			Addr:        fmt.Sprintf("%s:%s", proxy.Get("ip").String(), proxy.Get("port").String()),
			Type:        s.Type,
			Country:     "cn",
			Source:      "pinyi",
			Province:    "guangdong",
			Isp:         proxy.Get("isp").String(),
			City:        proxy.Get("city").String(),
			CreatedAt:   utils.GetUTCTimeNowStr(),
			ExpiredAt:   utils.Convert2UTCTimeStr(expiredTimeStr),
			ExpiredTime: time.Now().Unix() + 30,
		})
	}

	return
}
func (s *PinYi) GetApi() string {
	return s.Api
}