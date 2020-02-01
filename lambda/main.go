package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	lp "github.com/jsmorph/lambdaperf"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func HandleRequest(ctx context.Context, in lp.InMessage) (string, error) {
	var (
		then = time.Now().UTC()
	)

	out := &lp.OutMessage{
		Computes: in.Computations(),
	}

	for i := 0; i < out.Computes; i++ {
		if 0 < i {
			out.BlockTime += in.Block()
		}
		works, ms := in.Compute()
		out.Worked += works
		out.ComputeTime += ms
	}

	elapsed := time.Now().Sub(then)
	out.Elapsed = elapsed.Nanoseconds() / 1000 / 1000
	out.In = in

	js, err := json.Marshal(out)
	if err != nil {
		js = []byte(fmt.Sprintf("Marshal err %s on %#v", err, out))
	}
	return fmt.Sprintf("%s", js), nil
}

func main() {
	lambda.Start(HandleRequest)
}
