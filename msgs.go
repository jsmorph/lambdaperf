package lambdaperf

import (
	"encoding/json"
	"strconv"
	"time"

	"gonum.org/v1/gonum/stat/distuv"
)

type InMessage struct {
	T           time.Time
	ComputeTime distuv.Gamma
	BlockTime   distuv.Gamma
	Computes    distuv.Poisson
	N           int
}

type OutMessage struct {
	In          InMessage
	ComputeTime int64
	Worked      int
	BlockTime   int64
	Computes    int
	Elapsed     int64
}

func (in *InMessage) work(keys, rounds int) {
	for i := 0; i < rounds; i++ {
		m := make(map[string]interface{})
		for j := 0; j < keys; j++ {
			var (
				k = strconv.Itoa(j)
				v = strconv.Itoa(-j)
			)
			m[k] = v
		}
		js, _ := json.Marshal(&m)
		m = make(map[string]interface{})
		json.Unmarshal(js, &m)
	}
}

func (in *InMessage) Compute() (int, int64) {
	if in.ComputeTime.Alpha == 0 {
		in.ComputeTime.Alpha = 7
		in.ComputeTime.Beta = 1
	}
	var (
		ms      = int64(in.ComputeTime.Rand())
		then    = time.Now().UTC()
		elapsed time.Duration
		count   int
	)
	for {
		in.work(10, 100)
		count++
		if elapsed = time.Now().Sub(then); elapsed > time.Millisecond*time.Duration(ms) {
			break
		}
	}
	return count, elapsed.Nanoseconds() / 1000 / 1000
}

func (in *InMessage) Block() int64 {
	if in.BlockTime.Alpha == 0 {
		in.BlockTime.Alpha = 7
		in.BlockTime.Beta = 1
	}
	ms := int64(in.BlockTime.Rand())
	time.Sleep(time.Millisecond * time.Duration(ms))
	return ms
}

func (in *InMessage) Computations() int {
	if in.Computes.Lambda == 0 {
		in.Computes.Lambda = 4
	}
	return int(in.Computes.Rand())
}
