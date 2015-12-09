package netRelated

import (
	"crypto/md5"
	"fmt"
	"io"
	//"log"
	"net"
	"time"
)

const (
	IP   = "127.0.0.1"
	PORT = ":12345"
)

func Client() {
	addr, err := net.ResolveUDPAddr("udp4", IP+PORT)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}

	var name string
	var buf [512]byte

	go func() {
		for {
			data := make([]byte, 100)
			conn.Read(data)
			if fmt.Sprintf("%s", data[:3]) == "200" {
				fmt.Println("Connect is Good")
			} else {
				fmt.Printf("%s", data)
			}
		}
	}()

	for {
		fmt.Printf("Name yourself:")
		fmt.Scanf("%s", &name)
		if name == "" {
			fmt.Println("Empty Name,Try again")
		} else {
			fmt.Println("Wait for other players")
			break
		}
	}

	for i := 0; i < len(name) && i < 10; i++ {
		buf[32+i] = name[i]
	}

	hashlize(&buf, "ready")
	conn.Write([]byte(buf[:]))

	time.Sleep(time.Second)

	for {
		var answer string
		fmt.Printf(">>>")
		fmt.Scanf("%s", &answer)
		hashlize(&buf, answer)
		conn.Write([]byte(buf[:]))
	}
}

func hashlize(buf *[512]byte, str string) {
	if str == "" {
		return
	}
	m := md5.New()
	for i := 0; i < len(str); i++ {
		buf[10+32+i] = str[i]
	}
	io.WriteString(m, string(buf[32:len(str)+32+10]))
	hash := fmt.Sprintf("%x", m.Sum(nil))
	for i := 0; i < len(hash); i++ {
		buf[i] = hash[i]
	}
	for i := 32 + 10 + len(str); i < len(buf); i++ {
		buf[i] = 0
	}
}
