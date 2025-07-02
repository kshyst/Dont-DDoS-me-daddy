package models

type ReqData struct {
	UserIp         string `json:"user_ip"`
	RequestAddress string `json:"request_address"`
}

type RedisSaveData struct {
	ipAddr     string
	url        string
	timeStamp  int64
	expiration int
}
