package irand

import (
	"math/rand"
	"time"
)

//彩票随机算法
func GetProbability(prob []float32) int {
	v := rand.Float32()
	for i := 0; i < len(prob); i++ {
		pre := take_sum(prob, i)
		next := take_sum(prob, i+1)
		if v >= pre && v < next {
			return i
		}
	}
	return 0
}

func take_sum(prob []float32, idx int) float32 {
	var total float32 = 0
	for i := 0; i < idx; i++ {
		total += prob[i]
	}
	return total
}

func GetRand(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}
