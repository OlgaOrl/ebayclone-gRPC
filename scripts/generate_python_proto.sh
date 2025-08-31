#!/bin/bash

set -e

echo "=== Generating Python protobuf files ==="

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo "Error: Python 3 is not installed."
    exit 1
fi

# Install required Python packages
echo "Installing required Python packages..."
pip3 install grpcio grpcio-tools

# Generate Python protobuf files
echo "Generating Python protobuf files..."
python3 -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. proto/ebayclone.proto

echo "✓ Python protobuf files generated successfully!"
echo "You can now run the Python client example:"
echo "python3 client/example.py"
