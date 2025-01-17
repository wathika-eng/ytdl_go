#!/usr/bin/env bash

# Check if a URL was provided as an argument
if [ -z "$1" ]; then
    echo "Usage: $0 <video_url>"
    exit 1
fi

VIDEO_URL=$1
SERVICE_PID=""
SERVICE_STARTED=false

# Function to start the service
start_service() {
    echo "Starting the video download service..."
    go run main.go &
    SERVICE_PID=$!
    SERVICE_STARTED=true
    sleep 2 # Give the service some time to start
}

# Function to stop the service
stop_service() {
    if [ "$SERVICE_STARTED" = true ]; then
        echo "Stopping the video download service..."
        kill $SERVICE_PID
    fi
}

# Check if the service is already running
if ! curl -s "http://localhost:8080/api/downloads" > /dev/null; then
    start_service
fi

# Start a download and monitor progress
echo "Starting download and monitoring progress..."
RESPONSE=$(curl -s -X POST "http://localhost:8080/api/download/start" \
-H "Content-Type: application/json" \
-d "{\"url\": \"$VIDEO_URL\"}")

DOWNLOAD_ID=$(echo $RESPONSE | jq -r '.id')
if [ -z "$DOWNLOAD_ID" ]; then
    echo "Failed to retrieve download ID. Exiting."
    stop_service
    exit 1
fi
echo "Download ID: $DOWNLOAD_ID"

# Monitor progress until completion
while true; do
    STATUS=$(curl -s "http://localhost:8080/api/download/status?id=$DOWNLOAD_ID")
    PROGRESS=$(echo $STATUS | jq -r '.progress')
    CURRENT_STATUS=$(echo $STATUS | jq -r '.status')
    
    echo -ne "Progress: $PROGRESS% Status: $CURRENT_STATUS\r"
    
    if [ "$CURRENT_STATUS" = "completed" ]; then
        echo -e "\nDownload completed"
        stop_service
        exit 0
    elif [ "$CURRENT_STATUS" = "error" ]; then
        echo -e "\nDownload error"
        stop_service
        exit 1
    fi
    
    sleep 1
done