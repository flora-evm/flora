#!/bin/bash
set -e

echo "Setting up Flora blockchain for testing..."

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuration
CHAIN_ID="localchain_9000-1"
DENOM="flora"
HOME_DIR="$HOME/.flora"
MONIKER="test-node"

# Cleanup function
cleanup() {
    echo "Cleaning up..."
    killall florad 2>/dev/null || true
    rm -rf $HOME_DIR
}

# Check if binary exists
if [ ! -f "./build/florad" ]; then
    echo -e "${RED}Error: florad binary not found. Run 'make build' first.${NC}"
    exit 1
fi

echo -e "${GREEN}1. Cleaning previous state...${NC}"
cleanup

echo -e "${GREEN}2. Initializing chain...${NC}"
./build/florad init $MONIKER --chain-id $CHAIN_ID --home $HOME_DIR

echo -e "${GREEN}3. Creating test accounts...${NC}"
echo "test test test test test test test test test test test junk" | ./build/florad keys add validator --recover --keyring-backend test --home $HOME_DIR
./build/florad keys add alice --keyring-backend test --home $HOME_DIR
./build/florad keys add bob --keyring-backend test --home $HOME_DIR

echo -e "${GREEN}4. Adding genesis accounts...${NC}"
./build/florad genesis add-genesis-account validator 100000000000$DENOM --keyring-backend test --home $HOME_DIR
./build/florad genesis add-genesis-account alice 50000000000$DENOM --keyring-backend test --home $HOME_DIR
./build/florad genesis add-genesis-account bob 50000000000$DENOM --keyring-backend test --home $HOME_DIR

echo -e "${GREEN}5. Creating validator transaction...${NC}"
./build/florad genesis gentx validator 1000000000$DENOM \
    --chain-id $CHAIN_ID \
    --keyring-backend test \
    --home $HOME_DIR

echo -e "${GREEN}6. Collecting genesis transactions...${NC}"
./build/florad genesis collect-gentxs --home $HOME_DIR

echo -e "${GREEN}7. Validating genesis...${NC}"
./build/florad genesis validate --home $HOME_DIR

echo -e "${GREEN}8. Configuration summary:${NC}"
echo "Chain ID: $CHAIN_ID"
echo "Home Directory: $HOME_DIR"
echo "Denomination: $DENOM"
echo ""

echo -e "${YELLOW}To start the chain, run:${NC}"
echo "./build/florad start --home $HOME_DIR"
echo ""

echo -e "${YELLOW}To query the chain (after starting):${NC}"
echo "./build/florad status --home $HOME_DIR"
echo "./build/florad query bank balances \$(./build/florad keys show alice -a --keyring-backend test --home $HOME_DIR) --home $HOME_DIR"
echo ""

echo -e "${YELLOW}To test a transaction:${NC}"
echo "./build/florad tx bank send validator \$(./build/florad keys show alice -a --keyring-backend test --home $HOME_DIR) 1000$DENOM --chain-id $CHAIN_ID --keyring-backend test --home $HOME_DIR --fees 50$DENOM -y"

echo -e "\n${GREEN}Setup complete!${NC}"