package sources

import (
	"fmt"
	"github.com/tidwall/gjson"
	"go_proxy_pool/module"
	"go_proxy_pool/utils"
)

type DMDaili struct {
	Api      string
	Type     string //http/socks5
	Province string //省份
	//Verify   bool
}

func (s *DMDaili) Parse(resBody []byte) (proxies []*module.ProxyIp) {
	result := gjson.ParseBytes(resBody)
	if result.Get("code").Int() != 0 {
		return
	}
	for _, proxy := range result.Get("data").Array() {
		expiredTimeStr, err := utils.ParseTimeStr("2006/1/2 15:04:05", proxy.Get("endtime").String())
		expiredTimestamp := expiredTimeStr.Unix()
		if err != nil {
			continue
		}
		proxies = append(proxies, &module.ProxyIp{
			Addr:        fmt.Sprintf("%s:%d", proxy.Get("ip").String(), proxy.Get("port").Int()),
			Type:        s.Type,
			Country:     "cn",
			Source:      "dmdaili",
			Province:    s.Province,
			City:        proxy.Get("city").String(),
			CreatedAt:   utils.GetUTCTimeNowStr(),
			ExpiredAt:   utils.Convert2UTCTimeStr(expiredTimeStr),
			ExpiredTime: expiredTimestamp,
		})
	}

	return
}
func (s *DMDaili) GetApi() string {
	return s.Api
}