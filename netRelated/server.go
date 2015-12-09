package netRelated

import (
	"crypto/md5"
	"fmt"
	"github.com/xuzhenglun/calc24-muti/calc24"
	"io"
	"log"
	"math/rand"
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
var done chan bool

func (this Server) Listen() {
	done = make(chan bool, 4)
	this.ch = make(chan string, 4)
	addr, err := net.ResolveUDPAddr("udp4", ":"+this.Port)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("Listen at " + this.ListernIp + ":" + this.Port)

	Question.Winer = "NULL"
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
			log.Printf("%s", string(buf[:]))
			log.Println("BAD Request")
			return
		}
		_, err = conn.WriteToUDP([]byte("200"), addr)

		if len(data) >= 15 && data[10:15] == "ready" || this.PeopleNum < 4 {
			log.Println("one People is ready!")
			log.Println(addr.Port)
			this.gamer[this.PeopleNum] = string(buf[32 : 10+32]) //dirty hack
			this.PeopleNum++
			go this.obverser(conn, addr)
			return
		}

		//log.Println("HERE IM HERE")
		if this.ready && Question.Winer == "NULL" {
			//if len(data) < 10 {
			this.ch <- data
			//}
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
			if Question.Winer != "NULL" {
				break
			}

		}
	}()
}

func (this *Server) playgame() {
	for {
		if this.PeopleNum < 4 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		Question.A = rander()
		Question.B = rander()
		Question.C = rander()
		Question.D = rander()

		var ifResultable bool
		Question.Ans, ifResultable = Question.CalcAnswer()
		if ifResultable == false {
			log.Println("unable to resolve, aborded")
			continue
		} else {
			log.Println(Question.Ans)
		}
		log.Println("GAME IS ON")
		Question.Winer = "NULL"
		this.ready = true

		for {
			ans := <-this.ch
			log.Println("get ans: " + ans[10:])
			if calc24.Check(ans[10:]) == 24.0 {
				fmt.Printf("%s\n", ans[:10])
				fmt.Printf("%x\n", ans[:10])

				Question.Winer = fmt.Sprintf("%s", ans[:10])
				this.ready = false
				this.PeopleNum = 0
				log.Println("GAME OVER")
				for i := 0; i < 4; i++ {
					done <- true
				}
				break
			} else {
				//conn.WriteToUDP("Wrong Ans\n Take your time,No pressure\n", addr)
			}
		}
	}
}

func rander() int {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := rnd.Intn(12) + 1
	return vcode
}

func (this Server) obverser(conn *net.UDPConn, addr *net.UDPAddr) {
	log.Println("obverse ready")
	<-done

	winmesg := Question.Winer + " Winded,Typing \"ready\"to have another try!\n"
	log.Println(winmesg)
	_, err := conn.WriteToUDP([]byte(winmesg), addr)
	if err != nil {
		log.Println(err)
	}

	log.Println("obverse exiting")
}
