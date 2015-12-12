package main

import (
	"github.com/xuzhenglun/calc24-muti/netRelated"
)

func main() {
	netRelated.Conf.Port = "12345"
	netRelated.Conf.ListernIp = "0.0.0.0"
	netRelated.Conf.NumPreGroup = 4
	var server netRelated.Server
	server.Start()
}
