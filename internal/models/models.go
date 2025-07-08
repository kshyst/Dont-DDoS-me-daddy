package models

type ReqData struct {
	UserIp         string `json:"user_ip"`
	RequestAddress string `json:"request_address"`
}

type RedisSaveData struct {
	IPAddr     string `json:"ipAddr"`
	URL        string `json:"url"`
	TimeStamp  int64  `json:"timeStamp"`
	Expiration int    `json:"expiration"` // in seconds
}
