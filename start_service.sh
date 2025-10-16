#!/bin/bash

# 启动服务并记录PID
nohup go run main.go > auto-checkin.log 2>&1 &
echo $! > auto-checkin.pid
echo "服务已启动，PID: $(cat auto-checkin.pid)"