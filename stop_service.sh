#!/bin/bash

# 读取PID文件并结束进程
if [ -f "auto-checkin.pid" ]; then
    PID=$(cat auto-checkin.pid)
    kill $PID
    rm -f auto-checkin.pid
    echo "服务已停止，PID: $PID"
else
    echo "未找到PID文件，服务可能未运行"
fi