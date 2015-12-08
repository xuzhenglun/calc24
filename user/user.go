package user

import (
	"crypto/md5"
	"encoding/hex"
)

type User struct {
	Name          string
	Passwd        string
	WantedPlayers int
}

func (this User) GetUserPassedHash() string {
	hash := md5.New()
	hash.Write([]byte(this.Passwd))
	return hex.EncodeToString(hash.Sum(nil))
}
