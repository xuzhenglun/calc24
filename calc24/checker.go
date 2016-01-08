package calc24

import (
	"github.com/keepzero/gocalc"
	"log"
	"regexp"
	"sort"
	"strconv"
)

func (this *Question) Check(str string) bool {

	r, err := regexp.Compile(`\d+`)
	if err != nil {
		log.Println(err)
		return false
	}
	ansstr := r.FindAllString(str, -1)

	if len(ansstr) > 4 {
		return false
	}

	ansnum := make([]int, 4)
	for i, v := range ansstr {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Println(err)
			return false
		}
		ansnum[i] = n
	}
	log.Println("asd", ansnum)
	sort.Ints(ansnum)

	if this.A != ansnum[0] {
		return false
	}
	if this.B != ansnum[1] {
		return false
	}
	if this.C != ansnum[2] {
		return false
	}
	if this.D != ansnum[3] {
		return false
	}

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
} //ToDO:  auth 4 number weather is belog to the game
