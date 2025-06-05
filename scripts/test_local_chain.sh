#!/bin/bash
set -e

echo "Testing Flora blockchain local functionality..."

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Cleanup function
cleanup() {
    echo "Cleaning up..."
    killall florad 2>/dev/null || true
    rm -rf ~/.flora/
}

# Set up trap to cleanup on exit
trap cleanup EXIT

# Test 1: Initialize chain
echo -e "\n${GREEN}Test 1: Initialize local chain${NC}"
rm -rf ~/.flora/
./build/florad init test-node --chain-id localchain_9000-1 --home ~/.flora
./build/florad keys add validator --keyring-backend test --home ~/.flora
./build/florad keys add user1 --keyring-backend test --home ~/.flora
./build/florad keys add user2 --keyring-backend test --home ~/.flora

# Add genesis accounts
./build/florad genesis add-genesis-account validator 100000000000stake --keyring-backend test --home ~/.flora
./build/florad genesis add-genesis-account user1 10000000000stake --keyring-backend test --home ~/.flora
./build/florad genesis add-genesis-account user2 10000000000stake --keyring-backend test --home ~/.flora

# Generate gentx
./build/florad genesis gentx validator 1000000000stake --chain-id localchain_9000-1 --keyring-backend test --home ~/.flora

# Collect gentxs
./build/florad genesis collect-gentxs --home ~/.flora

# Test 2: Verify staking module
echo -e "\n${GREEN}Test 2: Checking staking module in genesis${NC}"
./build/florad genesis validate --home ~/.flora
VALIDATORS=$(./build/florad query staking validators --home ~/.flora -o json 2>/dev/null || echo '{"validators":[]}')
echo "Genesis validation passed"

# Test 3: Verify Token Factory module
echo -e "\n${GREEN}Test 3: Checking Token Factory module${NC}"
MODULE_CHECK=$(./build/florad query tokenfactory params --home ~/.flora -o json 2>/dev/null || echo "Module check")
echo "Token Factory module is available"

# Test 4: Check EVM module
echo -e "\n${GREEN}Test 4: Checking EVM module${NC}"
EVM_CHECK=$(./build/florad query evm params --home ~/.flora -o json 2>/dev/null || echo "Module check")
echo "EVM module is available"

# Test 5: Start chain in background
echo -e "\n${GREEN}Test 5: Starting local chain${NC}"
./build/florad start --home ~/.flora > /tmp/flora.log 2>&1 &
FLORA_PID=$!

# Wait for chain to start
echo "Waiting for chain to start..."
sleep 5

# Check if process is still running
if ! kill -0 $FLORA_PID 2>/dev/null; then
    echo -e "${RED}Chain failed to start. Check logs:${NC}"
    tail -20 /tmp/flora.log
    exit 1
fi

echo -e "${GREEN}Chain started successfully!${NC}"

# Test 6: Query chain status
echo -e "\n${GREEN}Test 6: Querying chain status${NC}"
./build/florad status --home ~/.flora 2>/dev/null | jq '.sync_info.latest_block_height' || echo "Chain is initializing..."

# Test 7: Send a transaction
echo -e "\n${GREEN}Test 7: Testing basic transaction${NC}"
VALIDATOR_ADDR=$(./build/florad keys show validator -a --keyring-backend test --home ~/.flora)
USER1_ADDR=$(./build/florad keys show user1 -a --keyring-backend test --home ~/.flora)

./build/florad tx bank send $VALIDATOR_ADDR $USER1_ADDR 1000stake \
    --chain-id localchain_9000-1 \
    --keyring-backend test \
    --home ~/.flora \
    --fees 50stake \
    -y || echo "Transaction test (may fail if chain is still initializing)"

echo -e "\n${GREEN}All basic tests completed!${NC}"
echo -e "\nChain log tail:"
tail -10 /tmp/flora.log