package main

import (
	"github.com/xuzhenglun/calc24-muti/netRelated"
)

func main() {
	var server netRelated.Server
	server.Port = "12345"
	server.ListernIp = "0.0.0.0"
	server.Listen()
}
