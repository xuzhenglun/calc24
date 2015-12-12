package config

import (
	//"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	NumPreGroup int
	RemoteIP    string
	Port        string
	ListernIp   string
}

func (this *Config) GetConfig() bool {
	file, err := os.Open("config.json")
	if err != nil {
		log.Panicln("Config file is not exist")
	}
	defer file.Close()

	buf := make([]byte, 1024)
	conf := make([]byte, 0)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		conf = append(conf, buf[:n]...)
	}

	err = json.Unmarshal(conf, &this)
	if err != nil {
		return false
	}
	log.Println(this)
	return true
}
