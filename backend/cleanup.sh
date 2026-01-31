#!/bin/bash

# Backend Code Cleanup Script
# This script performs various code cleanup tasks

echo "======================================"
echo "🧹 Backend Code Cleanup"
echo "======================================"
echo ""

# Run go vet
echo "1️⃣ Running go vet..."
go vet ./response ./constants ./validators ./errors ./repositories ./services ./utils ./config ./models ./middleware
if [ $? -eq 0 ]; then
    echo "✅ go vet passed"
else
    echo "❌ go vet found issues"
fi
echo ""

# Run go fmt
echo "2️⃣ Running go fmt..."
go fmt ./response ./constants ./validators ./errors ./repositories ./services ./utils ./config ./models ./middleware
echo "✅ go fmt completed"
echo ""

# Check for unused dependencies
echo "3️⃣ Checking for unused dependencies..."
go mod tidy
echo "✅ go mod tidy completed"
echo ""

# Build to ensure everything compiles
echo "4️⃣ Testing compilation..."
go build -o /tmp/lottery-test-build .
if [ $? -eq 0 ]; then
    echo "✅ Build successful"
    rm -f /tmp/lottery-test-build
else
    echo "❌ Build failed"
fi
echo ""

echo "======================================"
echo "✅ Code cleanup completed!"
echo "======================================"
