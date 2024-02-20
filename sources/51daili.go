package sources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"go_proxy_pool/module"
	"go_proxy_pool/utils"
)

// 51daili
type WuyiDaili struct {
	Api          string
	Type         string //http/socks5
	AreaCodeList string //省份
	//Verify   bool
}

func (s *WuyiDaili) GetProvince(areaCode string) string {
	switch areaCode[:3] {
	case "440":
		return "guangdong"
	case "360":
		return "jiangxi"
	case "430":
		return "hunan"
	case "350":
		return "fujian"
	case "330":
		return "zhejiang"
	case "310":
		return "shanghai"
	}
	return ""
}
func (s *WuyiDaili) Parse(resBody []byte) (proxies []*module.ProxyIp) {
	result := gjson.ParseBytes(resBody)
	//fmt.Println("result:", result)
	if result.Get("code").Int() != 0 {
		return
	}
	for _, proxy := range result.Get("data").Array() {
		expiredTimeStr, err := utils.ParseTimeStr("2006-01-02 15:04:05", proxy.Get("expireTime").String())
		expiredTimestamp := expiredTimeStr.Unix()
		if err != nil {
			continue
		}
		proxies = append(proxies, &module.ProxyIp{
			Addr:        fmt.Sprintf("%s:%s", proxy.Get("ip").String(), proxy.Get("port").String()),
			Type:        s.Type,
			Country:     "cn",
			Source:      "51daili",
			Province:    s.GetProvince(proxy.Get("ipaddress").String()),
			City:        proxy.Get("IpAddressName").String(),
			CreatedAt:   utils.GetUTCTimeNowStr(),
			ExpiredAt:   utils.Convert2UTCTimeStr(expiredTimeStr),
			ExpiredTime: expiredTimestamp,
		})
	}

	return
}
func (s *WuyiDaili) GetApi() string {
	return s.Api
}