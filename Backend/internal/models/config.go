package models

type Conf struct {
	ID              string `json:id gorm:"primaryKey"`
	ListeningPort   string `json:"listening_port"  binding:"default=80"`
	RemoteLogServer string `json:"remote_logServer" `
}

type AppConf struct {
	ID            string `json:id gorm:"primaryKey"`
	ApplicationID string `json:application_id `
	RateLimit     int    `json:rate_limit binding:"default=50"`
	WindowSize    int    `json:window_size binding:"default=10"`
	DetectBot     bool   `json:detect_bot binding:"default=false"`
	HostName      string `json:"hostname" `
}
