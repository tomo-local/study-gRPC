#!/bin/bash

echo "🚀 Starting Email Service Setup..."

# Create .env file from example if it doesn't exist
if [ ! -f .env ]; then
    echo "📄 Creating .env file from .env.example..."
    cp .env.example .env
fi

echo "🐘 Starting PostgreSQL database..."
cd database
docker-compose up -d

echo "⏳ Waiting for database to be ready..."
sleep 10

echo "📧 Starting Email Service..."
cd ../app

# Initialize go modules if needed
if [ ! -f go.sum ]; then
    echo "📦 Downloading Go dependencies..."
    go mod tidy
fi

echo "🏃 Running Email Service..."
go run cmd/server/main.go
