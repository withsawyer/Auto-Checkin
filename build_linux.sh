#!/bin/bash

# 定义项目名称和输出目录
PROJECT_NAME="Auto-Checkin"
OUTPUT_DIR="dist"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 编译为 Linux 可执行文件
GOOS=linux GOARCH=amd64 go build -o "$OUTPUT_DIR/$PROJECT_NAME-linux-amd64"

# 检查编译是否成功
if [ $? -eq 0 ]; then
    echo "Build successful! Executable saved to $OUTPUT_DIR/$PROJECT_NAME-linux-amd64"
else
    echo "Build failed."
    exit 1
fi