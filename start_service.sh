#!/bin/bash

# 定义可执行文件路径和PID文件路径
EXECUTABLE="./dist/auto-checkin-linux-amd64"
PID_FILE="./auto-checkin.pid"

# 检查可执行文件是否存在
if [ ! -f "$EXECUTABLE" ]; then
    echo "Error: Executable file not found at $EXECUTABLE"
    exit 1
fi

# 后台运行可执行文件，并记录PID
nohup "$EXECUTABLE" > /dev/null 2>&1 &
PID=$!

# 将PID写入文件
echo "$PID" > "$PID_FILE"
echo "Started auto-checkin with PID: $PID"

# 检查进程是否在运行
if ps -p "$PID" > /dev/null; then
    echo "Process is running with PID: $PID"
else
    echo "Failed to start the process."
    rm -f "$PID_FILE"
    exit 1
fi