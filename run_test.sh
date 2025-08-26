#!/usr/bin/env bash
set -euo pipefail

# Ensure we clean up the container no matter what
cleanup() {
    echo "Stopping and removing Docker container..."
    docker compose down
}
trap cleanup EXIT

# Check for correct usage
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <expected_exit_code>"
    exit 1
fi

EXPECTED_EXIT_CODE="$1"

# Step 1: Build the Go project
echo "Building Go project..."
make build LINUX=1

# Step 2: Start Docker container
echo "Starting Docker container..."
docker compose up -d

# Step 3: Compile the test file inside the container
echo "Compiling test.bng inside Docker..."
docker compose exec bingus-dev bash -c "./bin/bingus test.bng"

# Step 4: Print green tick if compilation succeeded
echo -e "✅ Compilation successful"

# Step 5: Test the executable
echo "Running compiled executable..."
docker compose exec bingus-dev bash -c "./output/test"

# Get exit code of last command inside container
EXIT_CODE=$(docker compose exec bingus-dev bash -c 'echo $?')

if [ "$EXIT_CODE" -eq "$EXPECTED_EXIT_CODE" ]; then
    echo -e "✅ Test passed: exit code $EXIT_CODE matches expected $EXPECTED_EXIT_CODE"
else
    echo -e "❌ Test failed: exit code $EXIT_CODE does not match expected $EXPECTED_EXIT_CODE"
    exit 1
fi
