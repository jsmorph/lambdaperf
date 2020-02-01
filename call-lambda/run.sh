#!/bin/bash

set -e

N=${1:-100}

echo "mb,n,work,ms" > d.csv

for I in $(seq 1 8); do
    MB=$((I*128))
    echo "MB=$MB"
    aws lambda update-function-configuration --function-name PerfTest --memory-size $MB > set-memory.log
    sleep 5
    call-lambda -spec "$(cat ../spec.json)" -n $N | tee results-$MB.json
    cat results-$MB.json | jq -r -c "[$MB,.In.N,.Worked,.Elapsed,.Computes,.ComputeTime,.BlockTime]|@csv" | tee -a d.csv
done
