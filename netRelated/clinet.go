package netRelated

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/xuzhenglun/calc24-muti/conf"
	"net"
	"os"
	"strings"
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
			name = scan()
			if name == "" {
				fmt.Printf("Empty Name,Try again:")
			} else {
				break
			}
		}
		for {
			fmt.Printf("Enter a password:")
			passwd = scan()
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

	done := make(chan bool, 5)
	for {
		fmt.Printf("Enter \"Ready\" to find a new game.\n>>>")
		input := scan()
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

				//fmt.Println("▂▃▄▅▆▇█Game Over█▇▆▅▄▃▂▁")
				if notify, ok := response.Info.(map[string]interface{}); ok {
					winer := notify["Winer"].(string)
					ans := notify["Ans"].(string)
					if winer != conf.Name {
						fmt.Println("\n▂▃▄▅▆▇█ Defeated █▇▆▅▄▃▂▁\n")
						fmt.Println("▂▃▄▅▆▇█You Losted█▇▆▅▄▃▂▁\n")
						fmt.Printf(
							"Winder is %s and the Answer is %s.\nPress \"Enter\" to continue\n",
							winer, ans)
						<-input
						return
					} else {
						fmt.Printf("Congrations, You Wined.\n")
						return
					}
				}
			default:
				continue
			}
		case str := <-input:
			if str == "" {
				return
			}
			sendmsg(1, str)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func typing(ans chan<- string) {
	var input string
	fmt.Printf("Answer>>>")
	input = scan()
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

func scan() string {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	return strings.TrimSpace(s)
}
