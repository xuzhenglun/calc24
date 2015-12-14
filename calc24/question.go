package calc24

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Game struct {
	Secret
	Question
}
type Question struct {
	A, B, C, D int
}

type Secret struct {
	Ans   string
	Winer string
}

func (this Game) CalcAnswer() (answer string, ifResult bool) {

	ifResult = true
	a := this.A
	b := this.B
	c := this.C
	d := this.D

	In := [...]float32{float32(a), float32(b), float32(c), float32(d)}

	I := [...]int32{
		0x00010203,
		0x00010302,
		0x00020103,
		0x00020301,
		0x00030102,
		0x00030201,
		0x01000203,
		0x01000302,
		0x01020003,
		0x01020300,
		0x01030002,
		0x01030200,
		0x02000103,
		0x02000301,
		0x02010003,
		0x02010300,
		0x02030001,
		0x02030100,
		0x03000102,
		0x03000201,
		0x03010200,
		0x03010200,
		0x03020001,
		0x03020100,
	}
	ilen := len(I)
	MARK := [...]byte{'+', '-', '*', '/'}
	for i := 0; i < ilen; i++ {
		index1 := I[i] >> 24
		index2 := I[i] >> 16 & 0x0f
		index3 := I[i] >> 8 & 0x0f
		index4 := I[i] & 0x0f

		for j := 0; j < 0x40; j++ {
			m1 := calc(j>>4, In[index1], In[index2])
			m2 := calc(j>>2&0x03, In[index2], In[index3])
			m3 := calc(j&0x03, In[index3], In[index4])

			if calc(j&0x03, calc(j>>2&0x03, m1, In[index3]), In[index4]) == 24.0 {
				return fmt.Sprintf("((%d%c%d)%c%d)%c%d=24\n", int(In[index1]), MARK[j>>4], int(In[index2]), MARK[j>>2&0x03],
					int(In[index3]), MARK[j&0x03], int(In[index4])), ifResult
			}
			if calc(j>>2&0x03, m1, m3) == 24.0 {
				return fmt.Sprintf("(%d%c%d)%c(%d%c%d)=24\n", int(In[index1]), MARK[j>>4], int(In[index2]), MARK[j>>2&0x03],
					int(In[index3]), MARK[j&0x03], int(In[index4])), ifResult
			}
			if calc(j&0x03, calc(j>>4&0x03, In[index1], m2), In[index4]) == 24.0 {
				return fmt.Sprintf("(%d%c(%d%c%d))%c%d=24\n", int(In[index1]), MARK[j>>4], int(In[index2]), MARK[j>>2&0x03],
					int(In[index3]), MARK[j&0x03], int(In[index4])), ifResult
			}
			if calc(j>>4, In[index1], calc(j&0x03, m2, In[index4])) == 24.0 {
				return fmt.Sprintf("%d%c((%d%c%d)%c%d)=24\n", int(In[index1]), MARK[j>>4], int(In[index2]), MARK[j>>2&0x03],
					int(In[index3]), MARK[j&0x03], int(In[index4])), ifResult
			}
			if calc(j>>4, In[index1], calc(j>>2&0x03, In[index2], m3)) == 24.0 {
				return fmt.Sprintf("%d%c(%d%c(%d%c%d))=24\n", int(In[index1]), MARK[j>>4], int(In[index2]), MARK[j>>2&0x03],
					int(In[index3]), MARK[j&0x03], int(In[index4])), ifResult
			}
		}
	}
	return fmt.Sprintln("No Answer!!!"), false
}

func calc(t int, a, b float32) float32 {
	switch t {
	case 0:
		return a + b
	case 1:
		return a - b
	case 2:
		return a * b
	case 3:
		return a / b
	}

	return 0.0
}

func rander() int {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := rnd.Intn(12) + 1
	return vcode

}

func New() Game {
	var game Game
	ok := false
	for {
		game.A = rander()
		game.B = rander()
		game.C = rander()
		game.D = rander()
		if game.Ans, ok = game.CalcAnswer(); ok {
			break
		}
	}
	log.Println(game)
	game.Winer = ""
	return game
}
