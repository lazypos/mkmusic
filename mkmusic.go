package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	max_level = 4 //最大音差
)

//0,3,4,5,6,7,11,12,13,14,15,16,17,21,22,23,24,25 0
//5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23
func InitYinMap(m map[int]int) {
	m[0] = 5
	m[1] = 6
	m[2] = 7
	m[3] = 7
	m[4] = 8
	m[5] = 8
	m[6] = 8
	m[7] = 8
	m[8] = 9
	m[9] = 9
	m[10] = 9
	m[11] = 9
	m[12] = 10
	m[13] = 10
	m[14] = 10
	m[15] = 10
	m[16] = 10
	for j := 11; j < 18; j++ {
		for i := 17; i < 24; i++ {
			m[i+(j-10)*7] = j
		}
	}
	m[66] = 18
	m[67] = 18
	m[68] = 18
	m[69] = 18
	m[70] = 18
	m[71] = 19
	m[72] = 19
	m[73] = 19
	m[74] = 19
	m[75] = 20
	m[76] = 20
	m[77] = 20
	m[78] = 20
	m[79] = 21
	m[80] = 21
	m[81] = 22
	m[82] = 5
	m[83] = 5
	m[84] = 5
	m[85] = 23
	m[86] = 23
	m[87] = 23
	m[88] = 23
}

const (
	pai = 60 // 60拍
)

func InitPaiMap(m map[int]int) {
	m[0] = 1
	for i := 1; i < 4; i++ {
		m[i] = 2
	}
	for i := 4; i < 12; i++ {
		m[i] = 3
	}
	for i := 12; i < 20; i++ {
		m[i] = 4
	}
	for i := 20; i < 24; i++ {
		m[i] = 5
	}
	m[24] = 6
}

func GetNextYin(m map[int]int, last int) (int, int) {
	if last == 0 {
		start := rand.Intn(6) + 11
		return start, start
	}
	for {
		s := rand.Intn(88)
		yin := m[s]
		if int(math.Abs(float64(yin-last))) <= max_level {
			if yin == 5 || yin == 23 {
				return yin, 0
			}
			if yin < 11 {
				return yin, yin - 3
			}
			if yin > 17 {
				return yin, yin + 3
			}
			return yin, yin
		}
	}
}

func main() {
	mYin := make(map[int]int)
	mPai := make(map[int]int)
	InitYinMap(mYin)
	InitPaiMap(mPai)

	rand.Seed(time.Now().UnixNano())

	last := 0
	yin := 0
	buf := bytes.NewBufferString("")
	for i := 0; i < 60; i++ {
		n := rand.Intn(24)
		count := mPai[n]
		for i := 0; i < count; i++ {
			last, yin = GetNextYin(mYin, last)
			buf.WriteString(fmt.Sprintf(`%d,`, yin))
		}
		buf.WriteString("|")
	}
	log.Println(buf.String())
}
