package utils

import (
	"go_proxy_pool/module"
	"sort"
	"time"
)

const StandardTimeFormat = "2006-01-02 15:04:05"
const CheckValidTolerance = int64(6)

//const GetValidTolerance = int64(6)

func GetUTCTimeNowStr() string {
	currentTime := time.Now().UTC()
	// 将 UTC 时间格式化为字符串
	utcTimeString := currentTime.Format(StandardTimeFormat)
	return utcTimeString
}

func Convert2UTCTimeStr(t time.Time) string {
	return t.UTC().Format(StandardTimeFormat)
}

func ParseTimeStr(format string, timeStr string) (time.Time, error) {
	// 使用模板在字符串得到时间实例
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tmp, err := time.ParseInLocation(format, timeStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return tmp, nil

}

func IsValidProxy(px *module.ProxyIp) bool {
	if px.Addr == "" {
		return false
	}
	//exTime, err := ParseTimeStr(StandardTimeFormat, px.ExpiredAt)
	//if err != nil {
	//	return false
	//}

	if px.ExpiredTime-time.Now().Unix() > CheckValidTolerance {
		return true
	}
	return false

}

func IsExpiredProxy(px *module.ProxyIp) bool {
	if px.Addr == "" {
		return true
	}
	//exTime, err := ParseTimeStr(StandardTimeFormat, px.ExpiredAt)
	//if err != nil {
	//	return false
	//}
	//expiredTolerance := int64(4)
	if px.ExpiredTime-time.Now().Unix() > 0 {
		return false
	}
	return true

}

func SortedIPsByReqNum(ips []*module.ProxyIp) []*module.ProxyIp {
	sort.Slice(ips, func(i, j int) bool {

		if ips[i].RequestNum == ips[j].RequestNum {
			//ti, _ := ParseTimeStr(StandardTimeFormat, ips[i].ExpiredAt)
			//tj, _ := ParseTimeStr(StandardTimeFormat, ips[j].ExpiredAt)
			return ips[i].ExpiredTime > ips[j].ExpiredTime
		}
		return ips[i].RequestNum < ips[j].RequestNum
	})

	return ips
}

// UniquePI 去除结构体数组中指定键值相同的元素
func UniquePI(proxies []*module.ProxyIp) []*module.ProxyIp {
	uniqueMap := make(map[string]bool)
	var result []*module.ProxyIp

	for _, p := range proxies {
		k := p.Addr
		if !uniqueMap[k] {
			uniqueMap[k] = true
			result = append(result, p)
		}
	}

	return result
}