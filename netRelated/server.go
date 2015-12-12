package netRelated

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/xuzhenglun/calc24-muti/conf"
	"log"
	"net"
)

type Server struct {
	ListernIp   string
	ListernPort string
	NumPreGroup int
}

var Conf config.Config
var groups map[string]Group

func (this *Server) Start() {
	if Conf.ListernIp != "" {
		Conf.GetConfig()
	}
	this.ListernIp = Conf.ListernIp
	this.ListernPort = Conf.Port
	this.NumPreGroup = Conf.NumPreGroup
	udpaddr, err := net.ResolveUDPAddr("udp", this.ListernIp+":"+this.ListernPort)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		panic(err)
	}

	for {
		log.Printf("Listern at %v:%v", this.ListernIp, this.ListernPort)
		buf := make([]byte, 512)
		len, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		go connectHandler(conn, len, addr, buf)
	}
}

func connectHandler(conn *net.UDPConn, len int, addr *net.UDPAddr, buf []byte) {
	buf = bytes.TrimSpace(buf)
	jsonReq, err := base64.StdEncoding.DecodeString(string(buf))
	if err != nil {
		log.Println("BAD requset")
		return
	}
	var req Information
	err = json.Unmarshal(jsonReq, &req)
	if err != nil {
		log.Panicln("Json decode fail")
	}
	log.Println(req)

	req.ClientAddr = addr

	log.Printf("%s Joined whose Hash is %s,and he/she want to do Action %d", req.ClientName, req.ClientHash, req.Status)
	log.Printf("Extre Info he/she come with is %v", req.Info)

	switch req.Status {
	case 0: //Start a new game request
		for _, group := range groups {
			if group.Now < Conf.NumPreGroup {
				group.date <- &req
				group.Now++
				return
			}
		}
		var newgroup Group
		newgroup.UUID = req.GroupUUID
		newgroup.Now = 1
		groups[req.GroupUUID] = newgroup
		newgroup.RunGroup()
		newgroup.date <- &req
	case 1:
		groups[req.GroupUUID].date <- &req
	}
}
