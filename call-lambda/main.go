package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	lp "github.com/jsmorph/lambdaperf"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func main() {

	var (
		spec = flag.String("spec", "{}", "Lambda spec")
		n    = flag.Int("n", 3, "Number of calls")
	)

	flag.Parse()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{})
	// Region: aws.String(is.Getenv("AWS_DEFAULT_REGION"))})

	var in lp.InMessage
	if err := json.Unmarshal([]byte(*spec), &in); err != nil {
		panic(err)
	}

	for i := 0; i < *n; i++ {
		in.N = i
		in.T = time.Now().UTC()
		js, err := json.Marshal(&in)
		if err != nil {
			panic(err)
		}

		result, err := client.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("PerfTest"),
			Payload:      js,
		})

		if err != nil {
			panic(err)
		}
		var x string
		if err = json.Unmarshal(result.Payload, &x); err != nil {
			panic(err)
		}

		fmt.Println(x)
	}
}
