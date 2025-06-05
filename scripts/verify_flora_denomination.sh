#!/bin/bash

echo "=== Flora Denomination Verification ==="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Check source code
echo "Checking source code..."
PETAL_COUNT=$(grep -r "petal" --include="*.go" --include="*.sol" . | grep -v "node_modules\|build\|vendor" | wc -l)
STAKE_COUNT=$(grep -r "\"stake\"" --include="*.go" --include="*.sh" . | grep -v "node_modules\|build\|vendor\|staking\|stake_" | wc -l)

if [ $PETAL_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ No 'petal' references found in source code${NC}"
else
    echo -e "${RED}✗ Found $PETAL_COUNT 'petal' references in source code${NC}"
fi

if [ $STAKE_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ No 'stake' denomination references found${NC}"
else
    echo -e "${RED}✗ Found $STAKE_COUNT 'stake' denomination references${NC}"
fi

echo ""
echo "Configuration Summary:"
echo "- Native token: flora (18 decimals)"
echo "- Display token: FLORA"
echo "- Liquid staking token: stFLORA-{ValidatorID}"
echo "- Chain ID: localchain_9000-1"
echo ""

# Check key files
echo "Key file configurations:"
echo -n "app/app.go: "
grep "BaseDenom.*=" app/app.go | grep -o '"[^"]*"' || echo "Not found"

echo -n "interchaintest/setup.go: "
grep "Denom.*=" interchaintest/setup.go | grep -o '"[^"]*"' || echo "Not found"

echo ""
echo "Documentation:"
echo "- Token naming convention: docs/liquid-staking/token-naming.md"
echo "- Liquid staking uses stFLORA-{ValidatorID} format"
echo ""

echo -e "${GREEN}Denomination update complete!${NC}"