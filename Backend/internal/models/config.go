package models

type Conf struct {
	ID              string `json:id gorm:"primaryKey"`
	ListeningPort   string `json:listening_port`
	RemoteLogServer string `json:remote_logServer`
	RateLimit       int    `json:rate_limit`
	WindowSize      int    `json:window_size`
}
