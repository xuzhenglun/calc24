package user

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
)

type User struct {
	UUID string
	Addr *net.UDPAddr
	Name string
}

func (this User) GenUUID(passwdhash string) string {
	m := md5.New()
	m.Write([]byte(passwdhash))
	m.Write([]byte(this.Name))
	this.UUID = fmt.Sprintf("%s", m.Sum(nil))
	log.Println("Gen UUID :" + this.UUID + "from User :" + this.Name)
	return this.UUID
}
