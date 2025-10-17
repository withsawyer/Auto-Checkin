#!/bin/bash

# 定义PID文件路径
PID_FILE="./auto-checkin.pid"

# 检查PID文件是否存在
if [ ! -f "$PID_FILE" ]; then
    echo "Error: PID file not found at $PID_FILE"
    exit 1
fi

# 读取PID
PID=$(cat "$PID_FILE")

# 检查进程是否在运行
if ps -p "$PID" > /dev/null; then
    echo "Stopping process with PID: $PID"
    kill "$PID"
    sleep 2
    if ps -p "$PID" > /dev/null; then
        echo "Failed to stop the process. Trying force kill..."
        kill -9 "$PID"
    fi
else
    echo "Process with PID: $PID is not running."
fi

# 清理PID文件
rm -f "$PID_FILE"
echo "PID file removed."