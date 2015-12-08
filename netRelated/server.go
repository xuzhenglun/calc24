package netRelated

import (
	"crypto/md5"
	"fmt"
	"github.com/xuzhenglun/calc24-muti/calc24"
	"io"
	"log"
	"net"
	"time"
)

type Server struct {
	ListernIp string
	Port      string
	PeopleNum int
	ready     bool
	ch        chan string
	gamer     [4]string
	net.UDPAddr
}

var Question calc24.Game

func (this Server) Listen() {
	this.ch = make(chan string, 4)
	Question.Winer = "NULL"

	addr, err := net.ResolveUDPAddr("udp4", ":"+this.Port)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("Listen at " + this.ListernIp + ":" + this.Port)

	go this.playgame()

	for {
		this.udpHandler(conn)
	}
}

func (this *Server) udpHandler(conn *net.UDPConn) {
	var buf [512]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		log.Println(err)
		return
	}

	for id, item := range this.gamer {
		//log.Printf("item:%x\nbuf:%x\n", item, string(buf[32:42]))
		if item == string(buf[10:10+32]) || this.PeopleNum < 4 {
			break
		} else if id < 4 {
			continue
		}
		log.Println("Sorry,Server is full")
		return
	}

	go func() {
		var i int
		for i = 42; i < 500; i++ {
			if buf[i] == 0 {
				break
			}
		}
		data := string(buf[32:i])
		m := md5.New()
		io.WriteString(m, data)
		hash := fmt.Sprintf("%x", m.Sum(nil))
		if string(buf[:32]) != hash {
			log.Printf("%x", string(buf[:32]))
			log.Println("BAD Request")
			return
		}
		_, err = conn.WriteToUDP([]byte("200"), addr)

		if len(data) >= 15 && data[10:15] == "ready" && this.PeopleNum < 4 {
			log.Println("one People is ready!")
			log.Println(addr.Port)
			this.gamer[this.PeopleNum] = string(buf[32 : 10+32]) //dirty hack
			this.PeopleNum++
			return
		}

		//log.Println("HERE IM HERE")
		if this.ready && Question.Winer == "NULL" {
			this.ch <- data
		}
	}()

	go func() {
		for {
			if Question.Winer != "NULL" {
				winmesg := Question.Winer + " Winded"
				log.Println(winmesg)
				_, err := conn.WriteToUDP([]byte(winmesg), addr)
				if err != nil {
					log.Println(err)
				}
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	go func() {
		for {
			if this.ready == true {
				ques := fmt.Sprintf("%d, %d, %d, %d \n", Question.A, Question.B, Question.C, Question.D)
				_, err := conn.WriteToUDP([]byte(ques), addr)
				if err != nil {
					log.Println(err)
				}
				//log.Println("WTF")
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (this *Server) playgame() {
	for {
		if this.PeopleNum < 3 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		Question.A = 2
		Question.B = 4
		Question.C = 9 //Mock data
		Question.D = 9 //ToDo: Generate Question

		var ifResultable bool
		Question.Ans, ifResultable = Question.CalcAnswer()
		if ifResultable == false {
			log.Println("unable to resolve, aborded")
			continue
		}
		log.Println("GAME IS ON")
		this.ready = true

		for {
			ans := <-this.ch
			log.Println("get ans: " + ans[10:])
			if calc24.Check(ans[10:]) == "24" {
				Question.Winer = ans[:10]
			}
		}
	}
}
