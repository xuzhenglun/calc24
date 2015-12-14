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

func (this *User) GenUUID(passwdhash string) string {
	m := md5.New()
	m.Write([]byte(passwdhash))
	m.Write([]byte(this.Name))
	this.UUID = fmt.Sprintf("%x", m.Sum(nil))
	log.Println("Gen UUID :" + this.UUID + " from User :" + this.Name)
	return this.UUID
}

func GenUUID(hash, name string) string {
	m := md5.New()
	m.Write([]byte(hash))
	m.Write([]byte(name))
	UUID := fmt.Sprintf("%x", m.Sum(nil))
	return UUID
}
