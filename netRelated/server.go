package netRelated

import (
	//"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/xuzhenglun/calc24-muti/conf"
	"github.com/xuzhenglun/calc24-muti/user"
	"log"
	"net"
	"sync"
)

type Server struct {
	ListernIp   string
	ListernPort string
	NumPreGroup int
}

var Conf config.Config
var groups map[string]*Group
var Clients map[string]string
var groupmux sync.RWMutex
var clientmux sync.RWMutex

func (this *Server) Start() {
	groups = make(map[string]*Group)
	Clients = make(map[string]string)

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
	log.Printf("Listern at %v:%v", this.ListernIp, this.ListernPort)

	for {
		buf := make([]byte, 2048)
		len, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		go connectHandler(conn, len, addr, buf)
	}
}

func connectHandler(conn *net.UDPConn, len int, addr *net.UDPAddr, buf []byte) {
	jsonReq, err := base64.StdEncoding.DecodeString(string(buf[0:len]))
	if err != nil {
		log.Println("BAD requset")
		return
	}
	var req Information
	err = json.Unmarshal(jsonReq, &req)
	if err != nil {
		log.Println("Json decode fail")
		return
	}

	req.ClientAddr = addr

	clientUUID := user.GenUUID(req.ClientHash, req.ClientName)
	switch req.Status {
	case 0: //Start a new game request
		groupmux.RLock()

		for _, group := range groups {
			log.Println("looking for old group")
			if group.Now < Conf.NumPreGroup {
				clientmux.Lock()
				Clients[clientUUID] = group.UUID
				clientmux.Unlock()
				group.date <- &req
				groupmux.RUnlock()
				return
			}
		}
		groupmux.RUnlock()

		log.Println("create new group")
		newgroup := NewGroup(&req, conn)
		groupmux.Lock()
		groups[newgroup.UUID] = newgroup
		groupmux.Unlock()
		newgroup.date <- &req
		clientmux.Lock()
		Clients[clientUUID] = newgroup.UUID
		clientmux.Unlock()
		newgroup.RunGroup()
	case 1: //game date
		clientmux.RLock()
		c := Clients[clientUUID]
		clientmux.RUnlock()

		groupmux.Lock()
		groups[c].date <- &req
		groupmux.Unlock()
	}
}
