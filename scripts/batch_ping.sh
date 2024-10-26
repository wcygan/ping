#!/bin/bash

# Number of times to send the request
N=${REQUEST_COUNT:-10}

# URL and headers for the request
URL=${PING_URL:-"http://localhost:8080/ping.v1.PingService/Ping"}
CONTENT_TYPE="Content-Type: application/json"

# Function to get the current timestamp in milliseconds
get_current_timestamp_ms() {
    echo $(( $(date +%s) * 1000 ))
}

# Loop to send the request N times
for ((i=1; i<=N; i++))
do
    TIMESTAMP_MS=$(get_current_timestamp_ms)
    DATA="{\"timestamp_ms\": $TIMESTAMP_MS}"
    
    response=$(curl -s -o /dev/null -w "%{http_code}" -X POST $URL -H "$CONTENT_TYPE" -d "$DATA")
    
    if [ "$response" -ne 200 ]; then
        echo "Request $i failed with status code $response" >> error.log
    else
        echo "Request $i succeeded"
    fi
done