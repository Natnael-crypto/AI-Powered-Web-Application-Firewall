package models

type Conf struct {
	ID              string `json:"id" gorm:"primaryKey"`
	ListeningPort   string `json:"listening_port"  binding:"default=80"`
	RemoteLogServer string `json:"remote_logServer" `
}

type SystemEmail struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Email  string `json:"email" binding:"required,email"`
	Active bool   `json:"active" binding:"default=true"`
}

type AppConf struct {
	ID              string  `json:"id" gorm:"primaryKey"`
	ApplicationID   string  `json:"application_id" gorm:"unique;not null"`
	RateLimit       int     `json:"rate_limit" binding:"default=50"`
	WindowSize      int     `json:"window_size" binding:"default=10"`
	BlockTime       int     `json:"block_time" binding:"default=10"`
	DetectBot       bool    `json:"detect_bot" binding:"default=false"`
	HostName        string  `json:"hostname"`
	MaxPostDataSize float64 `json:"max_post_data_size" binding:"default=5"`
	Tls             bool    `json:"tls"`
}
