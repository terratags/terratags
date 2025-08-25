#!/bin/bash
set +e

# This script tests a specific example directory with terratags
# Usage: ./test_examples.sh <example_directory> <expected_exit_code> [config_file] [plan_file]

EXAMPLE_DIR=$1
EXPECTED_EXIT_CODE=$2
CONFIG_FILE=${3:-config.yaml}
PLAN_FILE=$4

echo "Testing example: $EXAMPLE_DIR (Expected exit code: $EXPECTED_EXIT_CODE, Config: $CONFIG_FILE)"

# Run terratags on the example directory or plan file
if [ -n "$PLAN_FILE" ]; then
  echo "Running plan validation with: ./bin/terratags -c ./examples/$CONFIG_FILE -plan ./examples/$EXAMPLE_DIR/$PLAN_FILE -i"
  ./bin/terratags -c ./examples/$CONFIG_FILE -plan ./examples/$EXAMPLE_DIR/$PLAN_FILE -i
else
  echo "Running directory validation with: ./bin/terratags -c ./examples/$CONFIG_FILE -dir ./examples/$EXAMPLE_DIR -i"
  ./bin/terratags -c ./examples/$CONFIG_FILE -dir ./examples/$EXAMPLE_DIR -i
fi

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
