#!/bin/zsh

# This script tests rate limiting by making multiple requests to a specified URL
HOST="http://localhost:8080"
ENDPOINT="/aaaCeCSK"
FULL_URL="$HOST$ENDPOINT"
TOTAL_REQUESTS=100
DELAY_BETWEEN_REQUESTS=0.1  # in seconds
RATE_LIMIT_EXCEEDED_COUNT=0
RATE_LIMIT_STATUS_CODE=429
SUCCESS_STATUS_CODE=200

echo "Starting rate limiting test..."
for i in {1..$TOTAL_REQUESTS}; do
    RESPONSE_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$FULL_URL")
    if [ "$RESPONSE_CODE" -eq "$RATE_LIMIT_STATUS_CODE" ]; then
        RATE_LIMIT_EXCEEDED_COUNT=$((RATE_LIMIT_EXCEEDED_COUNT + 1))
        echo "Request $i: Rate limit exceeded (HTTP $RESPONSE_CODE)"
    elif [ "$RESPONSE_CODE" -eq "$SUCCESS_STATUS_CODE" ]; then
        echo "Request $i: Success (HTTP $RESPONSE_CODE)"
    else
        echo "Request $i: Unexpected response (HTTP $RESPONSE_CODE)"
    fi
    sleep $DELAY_BETWEEN_REQUESTS
done

