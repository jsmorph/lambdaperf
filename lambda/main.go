package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	lp "github.com/jsmorph/lambdaperf"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func HandleRequest(ctx context.Context, in lp.InMessage) (string, error) {
	log.Printf("N=%d %#v", in.N, in)

	out := in.Run()

	js, err := json.Marshal(out)
	if err != nil {
		js = []byte(fmt.Sprintf("Marshal err %s on %#v", err, out))
	}
	return fmt.Sprintf("%s", js), nil
}

func main() {
	lambda.Start(HandleRequest)
}
