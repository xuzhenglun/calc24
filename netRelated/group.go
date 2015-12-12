package netRelated

import (
	"encoding/base64"
	"encoding/json"
	"github.com/xuzhenglun/calc24-muti/calc24"
	"github.com/xuzhenglun/calc24-muti/user"
	"log"
	"net"
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
	GroupUUID  string
	ClientHash string
	ClientName string
	ClientAddr *net.UDPAddr
	Info       string
}

type tellClient struct {
	Status int
	Info   interface{}
}

func (this Group) NotifiyAll(date []byte) {
	for _, client := range this.Clients {
		go func(date []byte) {
			this.conn.WriteToUDP(
				[]byte(base64.StdEncoding.EncodeToString(date)), client.Addr)
		}(date)
	}
}

func (this Group) RunGroup() {
	for {
		if this.Now < 4 {
			req := <-this.date
			var newClient user.User
			newClient.Addr = req.ClientAddr
			newClient.Name = req.ClientName
			this.Clients[newClient.GenUUID(req.ClientHash)] = newClient
		} else {
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
			ans := <-this.date
			if calc24.Check(ans.Info) {
				var info tellClient
				info.Status = 1
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
			}
		}
	}
}
