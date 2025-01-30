package models

type Conf struct {
	ID              string `json:id`
	ListeningPort   string `json:listening_port`
	RemoteLogServer string `json:remote_logServer`
}
