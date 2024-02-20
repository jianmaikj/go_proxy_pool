package module

type ProxyIp struct {
	Addr string //IP地址:端口
	//Port     string //代理端口
	Country  string //代理国家
	Province string //代理省份
	City     string //代理城市
	Isp      string //IP提供商
	Type     string //代理类型,http/socks5
	//Anonymity  string //代理匿名度, 透明：显示真实IP, 普匿：显示假的IP, 高匿：无代理IP特征
	CreatedAt   string //获取时间
	VerifiedAt  string //校验时间
	ExpiredAt   string //过期时间(UTC),格式：2006-01-01 15:00:00
	ExpiredTime int64  //10位时间戳
	Speed       string //代理响应速度,ms
	SuccessNum  int    //验证请求成功的次数
	RequestNum  int    //验证请求的次数
	Source      string //代理源
}