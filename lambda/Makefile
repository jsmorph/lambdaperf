all:	main.zip

main.zip: main
	build-lambda-zip --output main.zip main

main: main.go
	go build -o main main.go

.PHONY: create-lambda
create-lambda: main.zip
	aws lambda create-function --function-name PerfTest --runtime go1.x --zip-file fileb://main.zip --handler main --role arn:aws:iam::918528020511:role/LambdaPerf

.PHONY: update-lambda
update-lambda: main.zip
	aws lambda update-function-code --function-name PerfTest --zip-file fileb://main.zip

.PHONY: create-role
create-role:
	aws iam create-role --role-name LambdaPerf --assume-role-policy-document file://trust-policy.json

.PHONY: attach-policy
attach-policy:
	aws iam attach-role-policy --role-name LambdaPerf --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole

.PHONY: test-lambda
test-lambda:
	aws lambda invoke --function-name PerfTest output.json --payload "$(cat ../spec.json)" --log-type Tail | tee result.json
	cat result.json | jq -c -r .LogResult | base64 --decode
	cat output.json | jq -r .

.PHONY: set-memory
set-memory:
	aws lambda update-function-configuration --function-name PerfTest --memory-size 128

