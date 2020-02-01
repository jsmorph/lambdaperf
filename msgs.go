package lambdaperf

import (
	"encoding/json"
	"strconv"
	"time"

	"gonum.org/v1/gonum/stat/distuv"
)

type InMessage struct {
	T          time.Time
	StepKeys   int
	StepRounds int
	WorkSteps  distuv.Poisson
	Works      distuv.Poisson
	BlockTime  distuv.Gamma
	N          int
}

type OutMessage struct {
	In        InMessage
	WorkTime  int64
	Worked    int
	Steps     int
	BlockTime int64
	Elapsed   int64
	Version   string
}

func (in *InMessage) WorkStep() {
	if in.StepKeys <= 0 {
		in.StepKeys = 10
	}
	if in.StepRounds <= 0 {
		in.StepRounds = 10
	}
	for i := 0; i < in.StepRounds; i++ {
		m := make(map[string]interface{})
		for j := 0; j < in.StepKeys; j++ {
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

func (in *InMessage) Work() (int, time.Duration) {
	var (
		then  = time.Now()
		steps = int(in.WorkSteps.Rand())
	)
	for i := 0; i < steps; i++ {
		in.WorkStep()
	}
	return steps, time.Now().Sub(then)
}

func (in *InMessage) Block() time.Duration {
	then := time.Now()
	if in.BlockTime.Alpha == 0 {
		in.BlockTime.Alpha = 7
		in.BlockTime.Beta = 1
	}
	var (
		ms = int64(in.BlockTime.Rand())
		d  = time.Millisecond * time.Duration(ms)
	)
	time.Sleep(d)
	return time.Now().Sub(then)
}

func (in *InMessage) Run() *OutMessage {
	if in.T.IsZero() {
		in.T = time.Now().UTC()
	}

	var (
		then = time.Now().UTC()
		out  = &OutMessage{
			Worked:  1 + int(in.Works.Rand()),
			Version: "1",
		}
	)

	for i := 0; i < out.Worked; i++ {
		if 0 < i {
			out.BlockTime += in.Block().Nanoseconds() / 1000 / 1000
		}
		steps, d := in.Work()
		out.Steps += steps
		out.WorkTime += d.Nanoseconds() / 1000 / 1000
	}

	elapsed := time.Now().Sub(then)
	out.Elapsed = elapsed.Nanoseconds() / 1000 / 1000
	out.In = *in

	return out
}
