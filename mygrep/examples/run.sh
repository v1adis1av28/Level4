#!/bin/bash

set -e

echo "Starting workers..."
go run cmd/worker/main.go --mode=worker --port=8080 &
WORKER1_PID=$!

go run cmd/worker/main.go --mode=worker --port=8081 &
WORKER2_PID=$!

go run cmd/worker/main.go --mode=worker --port=8082 &
WORKER3_PID=$!

# Ждём, пока воркеры запустятся
sleep 2

echo "Workers started. Testing with coordinator..."

# Тестируем
echo -e "line1\nerror line2\nline3\nerror line4\nline5\nerror line6" | \
go run cmd/coordinator/main.go \
  error \
  -n \
  -workers="localhost:8080,localhost:8081,localhost:8082" \
  --quorum=2

echo "Test completed."

# Убиваем воркеров
kill $WORKER1_PID $WORKER2_PID $WORKER3_PID