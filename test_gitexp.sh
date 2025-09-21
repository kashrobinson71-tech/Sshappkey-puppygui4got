#!/bin/bash

# Simple test script for gitexp tool
echo "Testing gitexp tool..."

# Build the tool
echo "Building gitexp..."
go build -o gitexp main.go

# Test help
echo -e "\n=== Testing help ==="
./gitexp --help

# Test config list
echo -e "\n=== Testing config list ==="
./gitexp config list

# Test environment configuration
echo -e "\n=== Testing environment configuration ==="
GIT_USER_NAME="Test User" GIT_USER_EMAIL="test@example.com" ./gitexp env

# Test SSH help
echo -e "\n=== Testing SSH help ==="
./gitexp ssh --help

echo -e "\n=== All tests completed ==="