package calc24

import (
	"github.com/keepzero/gocalc"
	"log"
)

func Check(str string) float64 {
	ans, err := gocalc.Calc(str)
	if err != nil {
		log.Println(err)
		return -1
	}
	log.Printf("%s=%f", str, ans)
	return ans
}
