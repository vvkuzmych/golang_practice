#!/bin/bash

# Script to generate remaining examples for week_23

cd "$(dirname "$0")"

echo "Generating remaining examples for week_23..."
echo "This will create 80 more example files (channels, interfaces, slices, maps)"
echo ""

# Channels, Interfaces, Slices, Maps examples will be generated via Go
go run generate_examples.go

echo ""
echo "✅ All examples generated!"
echo ""
echo "📁 Structure:"
echo "  goroutines/  - 20 examples ✅ (already created)"
echo "  channels/    - 20 examples ⏳"
echo "  interfaces/  - 20 examples ⏳"
echo "  slices/      - 20 examples ⏳"
echo "  maps/        - 20 examples ⏳"
echo ""
echo "Run examples:"
echo "  cd goroutines && go run 01_basic.go"
echo "  cd channels && go run 01_basic.go"
echo "  cd interfaces && go run 01_basic.go"
echo "  cd slices && go run 01_basic.go"
echo "  cd maps && go run 01_basic.go"
