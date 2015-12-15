package netRelated

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/xuzhenglun/calc24-muti/calc24"
	"github.com/xuzhenglun/calc24-muti/user"
	"log"
	"math/rand"
	"net"
	"time"
)

type Group struct {
	Clients  map[string]user.User
	Question string
	conn     *net.UDPConn
	date     chan *Information
	UUID     string
	Now      int
}

type Information struct {
	Status     int
	ClientHash string
	ClientName string
	Info       string
	ClientAddr *net.UDPAddr
}

type tellClient struct {
	Status int
	Info   interface{}
}

func (this Group) NotifiyAll(date []byte) {
	for _, client := range this.Clients {
		go func(date []byte, client *net.UDPAddr) {
			log.Printf("Notifiy %v with date :%v", client, date)
			this.conn.WriteToUDP(
				[]byte(base64.StdEncoding.EncodeToString(date)), client)
		}(date, client.Addr)
	}
}

func (this Group) NotifiyOne(date []byte, addr *net.UDPAddr) {
	log.Printf("Notifiy %v with date :%v", addr, date)
	this.conn.WriteToUDP(
		[]byte(base64.StdEncoding.EncodeToString(date)), addr)
}

func (this Group) RunGroup() {
	log.Println("Grouper running")
	for {
		if this.Now < Conf.NumPreGroup { //ToDo: Just in converion. Change to len(map) later
			req := <-this.date
			var newClient user.User
			newClient.Addr = req.ClientAddr
			newClient.Name = req.ClientName
			clientUUID := newClient.GenUUID(req.ClientHash)
			this.Clients[clientUUID] = newClient
			this.Now++
			groups[this.UUID] = &this
			log.Println("group Add a new client,Now have " + string(this.Now+'0') + " and we need " + string(Conf.NumPreGroup+'0'))
		} else {
			break
		}
	}

	log.Println("players is ready!")
	game := calc24.New()
	var info tellClient
	info.Info = game.Question
	info.Status = 0
	jsonQuset, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
		return
	}
	this.NotifiyAll(jsonQuset)
	//game is ready

	for { //check answers
		ans := <-this.date
		if ans.Status != 1 {
			return
		}

		if ans.Info != "" && game.Check(ans.Info) {
			var info tellClient
			info.Status = 2
			var secret calc24.Secret
			secret.Winer = ans.ClientName
			secret.Ans = ans.Info
			info.Info = secret
			jsonInfo, err := json.Marshal(info)
			if err != nil {
				log.Println(err)
				return
			}
			this.NotifiyAll(jsonInfo)
			break
		} else {
			var info tellClient
			info.Status = 1
			info.Info = "Wrong Answer"
			jsonInfo, err := json.Marshal(info)
			if err != nil {
				log.Println(err)
				return
			}
			this.NotifiyOne([]byte(jsonInfo), ans.ClientAddr)
		}
	}
}

func NewGroup(req *Information, conn *net.UDPConn) *Group {
	var newgroup Group
	newgroup.Now = 0
	newgroup.Clients = make(map[string]user.User)
	newgroup.date = make(chan *Information, Conf.NumPreGroup)
	newgroup.conn = conn
	newgroup.GenUUID()
	return &newgroup
}

func (this *Group) GenUUID() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vocode := rnd.Intn(1<<17 - 1)
	m := md5.New()
	m.Write([]byte(string(vocode)))
	this.UUID = fmt.Sprintf("%x", m.Sum(nil))
}
