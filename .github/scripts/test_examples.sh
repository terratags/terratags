#!/bin/bash
set +e

# This script tests a specific example directory with terratags
# Usage: ./test_examples.sh <example_directory> <expected_exit_code> [config_file]

EXAMPLE_DIR=$1
EXPECTED_EXIT_CODE=$2
CONFIG_FILE=${3:-config.yaml}

echo "Testing example: $EXAMPLE_DIR (Expected exit code: $EXPECTED_EXIT_CODE, Config: $CONFIG_FILE)"

# Run terratags on the example directory
./bin/terratags -c ./examples/$CONFIG_FILE -dir ./examples/$EXAMPLE_DIR -i
ACTUAL_EXIT_CODE=$?

# Check if the exit code matches the expected exit code
if [ -n "$EXPECTED_EXIT_CODE" ] && [ "$ACTUAL_EXIT_CODE" -eq "$EXPECTED_EXIT_CODE" ]; then
  echo "✅ Test passed for $EXAMPLE_DIR"
  exit 0
else
  echo "❌ Test failed for $EXAMPLE_DIR"
  echo "Expected exit code: $EXPECTED_EXIT_CODE, Actual exit code: $ACTUAL_EXIT_CODE"
  exit 1
fi
