package calc24

import (
	"github.com/keepzero/gocalc"
	"log"
)

func Check(str string) bool {
	ans, err := gocalc.Calc(str)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Printf("%s=%f", str, ans)
	if ans == 24.0 {
		return true
	} else {
		return false
	}
}
