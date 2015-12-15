package netRelated

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/xuzhenglun/calc24-muti/conf"
	"net"
	"time"
)

var conf config.Config
var conn *net.UDPConn

func Client() {
	conf.GetConfig()
	fmt.Printf("Loaded Config: There will be %d players in a group.\nConnecing to %s:%s\n", conf.NumPreGroup, conf.ListernIp, conf.Port)
	res := make(chan *tellClient, 5)
	ans := make(chan string, 5)

	addr, err := net.ResolveUDPAddr("udp4", conf.RemoteIP+":"+conf.Port)
	if err != nil {
		panic(err)
	}
	conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}

	if conf.Name == "" || conf.Hash == "" {
		var name, passwd string
		for {
			fmt.Printf("Name yourself:")
			fmt.Scan(&name)
			if name == "" {
				fmt.Printf("Empty Name,Try again:")
			} else {
				break
			}
		}
		for {
			fmt.Printf("Enter a password:")
			fmt.Scan(&passwd)
			if passwd == "" {
				fmt.Printf("Please Enter Password:")
			} else {
				break
			}
		}
		conf.Name = name
		m := md5.New()
		m.Write([]byte(passwd))
		md5 := m.Sum(nil)
		m.Reset()
		m.Write([]byte("SALT"))
		m.Write(md5)
		conf.Hash = fmt.Sprintf("%x", m.Sum(nil))
		//conf.SaveConfig()
	}

	var input string
	done := make(chan bool, 5)
	for {
		fmt.Printf("Enter \"Ready\" to find a game.\n>>>")
		fmt.Scan(&input)
		if input == "Ready" {
			go func() {
				for {
					select {
					case <-done:
						return
					default:
						time.Sleep(2 * time.Second)
						fmt.Printf(".")
					}
				}
			}()
			sendmsg(0, "")
			response := recvmsg()

			switch response.Status {
			case 0:
				done <- true
				res <- response
				go recvall(res)
				go typing(ans)
				startGame(res, ans)
			default:
			}
		}
	}
}

func newMsg() *Information {
	var info Information
	info.ClientName = conf.Name
	info.ClientHash = conf.Hash
	return &info
}

func sendmsg(status int, str string) {
	msg := newMsg()
	msg.Status = status
	if status != 0 {
		msg.Info = str
	}
	jsonInfo, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	bash64Info := base64.StdEncoding.EncodeToString([]byte(jsonInfo))
	_, err = conn.Write([]byte(bash64Info))
	if err != nil {
		panic(err)
	}
}

func recvmsg() *tellClient {
	buf := make([]byte, 512)
	len, err := conn.Read(buf)
	if err != nil {
		panic(err)
		return nil
	}
	jsonRes, err := base64.StdEncoding.DecodeString(string(buf[0:len]))
	if err != nil {
		panic(err)
		return nil
	}
	var response tellClient
	json.Unmarshal(jsonRes, &response)
	return &response
}

func startGame(res <-chan *tellClient, input chan string) {
	for {
		select {
		case response := <-res:
			switch response.Status {
			case 0:
				if question, ok := response.Info.(map[string]interface{}); ok {
					fmt.Println("Game is ready!")
					fmt.Printf("There are 4 numbers: ")
					for _, v := range question {
						fmt.Printf("%.0f ", v.(float64))
					}
					fmt.Printf("\n")
				}
			case 1:
				if notify, ok := response.Info.(string); ok {
					fmt.Println(notify)
					go typing(input)
				}
			case 2:
				fmt.Println("▂▃▄▅▆▇█Game Over█▇▆▅▄▃▂▁")
				if notify, ok := response.Info.(map[string]interface{}); ok {
					fmt.Printf("Winder is %s and the Answer is %s.\n",
						notify["Winer"].(string), notify["Ans"].(string))
					return
				}
			default:
				continue
			}
		case str := <-input:
			sendmsg(1, str)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func typing(ans chan<- string) {
	var input string
	for {
		fmt.Printf("Answer>>>")
		fmt.Scan(&input)
		if input != "" {
			break
		}
	}
	ans <- input
}

func recvall(res chan<- *tellClient) {
	for {
		r := recvmsg()
		res <- r
		if r.Status == 2 {
			return
		}
	}
}
