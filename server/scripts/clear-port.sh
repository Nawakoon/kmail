#!/bin/bash

# check if port 8080 is already in use
PORT_PID=$(lsof -ti tcp:8080)
if [ -n "$PORT_PID" ]; then
    echo "Killing process on port 8080: $PORT_PID"
    kill $PORT_PID
fi