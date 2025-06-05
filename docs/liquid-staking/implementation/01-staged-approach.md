# Liquid Staking Staged Implementation Plan

## Overview
This document outlines a highly sequential, testable implementation approach for liquid staking on Flora. Each stage is designed to be independently testable with clear success criteria.

## Implementation Stages

### Stage 1: Basic Infrastructure (Week 1)
**Goal**: Set up foundational types and interfaces without any actual functionality.

#### 1.1 Define Core Types
```go
// x/liquidstaking/types/types.go
type TokenizationRecord struct {
    Id               uint64
    Validator        string
    Owner            string
    SharesTokenized  sdk.Int
}

type ModuleParams struct {
    GlobalLiquidStakingCap sdk.Dec
    ValidatorLiquidCap     sdk.Dec
    Enabled                bool
}
```

#### 1.2 Create Minimal Keeper
```go
// x/liquidstaking/keeper/keeper.go
type Keeper struct {
    storeKey storetypes.StoreKey
    cdc      codec.Codec
}
```

#### 1.3 Basic Genesis
- Define genesis state structure
- Implement import/export (empty for now)

**Tests**:
- Unit tests for type validation
- Genesis import/export with empty state
- Keeper instantiation

**Success Criteria**:
- [ ] Types compile without errors
- [ ] Genesis import/export works with empty state
- [ ] Keeper can be instantiated

---

### Stage 2: State Management (Week 2)
**Goal**: Implement basic CRUD operations for tokenization records.

#### 2.1 Store Operations
```go
// Implement in keeper
func (k Keeper) SetTokenizationRecord(ctx sdk.Context, record TokenizationRecord)
func (k Keeper) GetTokenizationRecord(ctx sdk.Context, id uint64) (TokenizationRecord, bool)
func (k Keeper) GetAllTokenizationRecords(ctx sdk.Context) []TokenizationRecord
func (k Keeper) DeleteTokenizationRecord(ctx sdk.Context, id uint64)
```

#### 2.2 Record Counter
```go
func (k Keeper) GetLastTokenizationRecordID(ctx sdk.Context) uint64
func (k Keeper) SetLastTokenizationRecordID(ctx sdk.Context, id uint64)
```

#### 2.3 Params Management
```go
func (k Keeper) GetParams(ctx sdk.Context) ModuleParams
func (k Keeper) SetParams(ctx sdk.Context, params ModuleParams)
```

**Tests**:
- CRUD operations for records
- Counter increment logic
- Params get/set
- Iterator tests

**Success Criteria**:
- [ ] Can store/retrieve tokenization records
- [ ] Counter increments correctly
- [ ] Params persist across blocks

---

### Stage 3: Basic Queries (Week 3)
**Goal**: Implement gRPC/REST queries without any business logic.

#### 3.1 Query Service
```proto
service Query {
    rpc TokenizationRecord(QueryTokenizationRecordRequest) returns (QueryTokenizationRecordResponse);
    rpc TokenizationRecords(QueryTokenizationRecordsRequest) returns (QueryTokenizationRecordsResponse);
    rpc Params(QueryParamsRequest) returns (QueryParamsResponse);
}
```

#### 3.2 CLI Commands
```go
// x/liquidstaking/client/cli/query.go
func GetQueryCmd() *cobra.Command
func CmdQueryTokenizationRecord() *cobra.Command
func CmdQueryParams() *cobra.Command
```

**Tests**:
- gRPC query tests
- CLI integration tests
- REST endpoint tests

**Success Criteria**:
- [ ] Can query records via gRPC
- [ ] CLI commands work
- [ ] REST endpoints respond correctly

---

### Stage 4: Simple Tokenization Logic (Week 4)
**Goal**: Implement basic tokenization without token minting or complex validation.

#### 4.1 Basic Message Handler
```go
func (k msgServer) TokenizeShares(ctx context.Context, msg *MsgTokenizeShares) (*MsgTokenizeSharesResponse, error) {
    // 1. Basic validation
    // 2. Create record
    // 3. Store record
    // 4. Return record ID
}
```

#### 4.2 Simple Validation
- Check delegator has delegation
- Check amount is positive
- Check module is enabled

**Tests**:
- Happy path tokenization
- Validation failures
- State persistence

**Success Criteria**:
- [ ] Can create tokenization record
- [ ] Validation rejects invalid inputs
- [ ] Record ID increments correctly

---

### Stage 5: Integration with Staking Module (Week 5)
**Goal**: Connect to actual staking module for delegation queries.

#### 5.1 Staking Keeper Interface
```go
type StakingKeeper interface {
    GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, bool)
    GetValidator(ctx sdk.Context, addr sdk.ValAddress) (stakingtypes.Validator, bool)
}
```

#### 5.2 Delegation Validation
- Verify delegation exists
- Check delegation amount sufficient
- Validate validator is active

**Tests**:
- Mock staking keeper tests
- Integration tests with real staking
- Edge cases (unbonding, jailed validators)

**Success Criteria**:
- [ ] Can query real delegations
- [ ] Rejects tokenization for invalid delegations
- [ ] Handles validator state correctly

---

### Stage 6: Token Factory Integration (Week 6)
**Goal**: Create LST denominations using Token Factory.

#### 6.1 Token Factory Keeper Interface
```go
type TokenFactoryKeeper interface {
    CreateDenom(ctx sdk.Context, creator string, subdenom string) (string, error)
    Mint(ctx sdk.Context, amount sdk.Coin, to string) error
}
```

#### 6.2 LST Denom Creation
- One denom per validator
- Format: `factory/{module}/{validator}`
- Create on first tokenization

**Tests**:
- Denom creation tests
- Duplicate prevention
- Minting tests

**Success Criteria**:
- [ ] Can create LST denoms
- [ ] Denoms are unique per validator
- [ ] Can mint LST tokens

---

### Stage 7: Basic Precompile Structure (Week 7)
**Goal**: Create precompile skeleton without actual functionality.

#### 7.1 Precompile Interface
```go
type LiquidStakingPrecompile struct {
    Address() common.Address
    RequiredGas(input []byte) uint64
    Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) ([]byte, error)
}
```

#### 7.2 ABI Registration
- Define method signatures
- Implement ABI parsing
- Gas estimation placeholders

**Tests**:
- ABI parsing tests
- Gas calculation tests
- Method dispatch tests

**Success Criteria**:
- [ ] Precompile compiles
- [ ] ABI methods recognized
- [ ] Gas estimates return values

---

### Stage 8: Read-Only Precompile Methods (Week 8)
**Goal**: Implement query methods in precompile.

#### 8.1 Query Methods
```solidity
function getRecord(uint256 recordId) external view returns (address, address, uint256, address)
function getGlobalStats() external view returns (uint256, uint256, uint256)
```

#### 8.2 Address Conversion
- Implement Ethereum ↔ Cosmos address mapping
- Cache conversions for efficiency

**Tests**:
- Query accuracy tests
- Address conversion tests
- Gas consumption tests

**Success Criteria**:
- [ ] Can query records from EVM
- [ ] Address conversion works
- [ ] Gas costs are reasonable

---

### Stage 9: Simple State-Changing Precompile (Week 9)
**Goal**: Implement tokenizeShares without actual token transfers.

#### 9.1 TokenizeShares Method
- Accept validator and amount
- Create tokenization record
- Return record ID
- No token minting yet

**Tests**:
- EVM integration tests
- State persistence tests
- Revert scenarios

**Success Criteria**:
- [ ] Can call from smart contract
- [ ] State changes persist
- [ ] Reverts handled correctly

---

### Stage 10: LST Token Contract (Week 10)
**Goal**: Deploy basic ERC20 for LST tokens.

#### 10.1 Simple LST Token
```solidity
contract SimpleLSTToken is ERC20 {
    address public precompile;
    uint256 public exchangeRate = 1e18; // 1:1 initially
    
    function mint(address to, uint256 amount) external onlyPrecompile {
        _mint(to, amount);
    }
}
```

#### 10.2 Token Deployment
- Deploy one token per validator
- Store token addresses in keeper

**Tests**:
- Token deployment tests
- Minting tests
- Access control tests

**Success Criteria**:
- [ ] Can deploy LST tokens
- [ ] Minting restricted to precompile
- [ ] Basic ERC20 functions work

---

### Stage 11: Complete Tokenization Flow (Week 11)
**Goal**: Connect all pieces for basic tokenization.

#### 11.1 Full Integration
1. Call precompile from contract
2. Validate delegation
3. Create record
4. Mint LST tokens
5. Return success

**Tests**:
- End-to-end tokenization tests
- Multi-validator tests
- Error handling tests

**Success Criteria**:
- [ ] Complete tokenization works
- [ ] LST tokens received
- [ ] State consistent

---

### Stage 12: Redemption Logic (Week 12)
**Goal**: Implement basic redemption without unbonding.

#### 12.1 Redeem Method
- Burn LST tokens
- Update record
- Mark for redemption

**Tests**:
- Redemption flow tests
- Balance verification
- State consistency

**Success Criteria**:
- [ ] Can burn LST tokens
- [ ] Records updated correctly
- [ ] No token/state leaks

---

### Stage 13: Reward Distribution (Week 13)
**Goal**: Basic reward handling without auto-compounding.

#### 13.1 Reward Tracking
- Track rewards per record
- Manual claim mechanism
- Simple distribution

**Tests**:
- Reward calculation tests
- Distribution tests
- Multi-user scenarios

**Success Criteria**:
- [ ] Rewards tracked correctly
- [ ] Can claim rewards
- [ ] No double-claiming

---

### Stage 14: Exchange Rate Updates (Week 14)
**Goal**: Implement dynamic exchange rates.

#### 14.1 Rate Calculation
- Calculate based on total staked/rewards
- Update mechanism (manual first)
- Apply to minting/burning

**Tests**:
- Rate calculation tests
- Update mechanism tests
- Precision tests

**Success Criteria**:
- [ ] Rates calculate correctly
- [ ] Updates apply properly
- [ ] No precision loss

---

### Stage 15: Auto-compounding (Week 15)
**Goal**: Make rewards auto-compound.

#### 15.1 Hook Integration
- Distribution module hooks
- Automatic rate updates
- Gas optimization

**Tests**:
- Hook integration tests
- Auto-update tests
- Gas usage tests

**Success Criteria**:
- [ ] Hooks trigger updates
- [ ] Rates update automatically
- [ ] Gas costs acceptable

---

### Stage 16: Slashing Support (Week 16)
**Goal**: Handle validator slashing.

#### 16.1 Slashing Integration
- Slashing module hooks
- Rate adjustments
- User notification events

**Tests**:
- Slashing simulation tests
- Rate adjustment tests
- Event emission tests

**Success Criteria**:
- [ ] Slashing reduces LST value
- [ ] All holders affected proportionally
- [ ] Events emitted correctly

---

### Stage 17: Advanced Features (Week 17-18)
**Goal**: Add caps, limits, and governance.

#### 17.1 Risk Parameters
- Global caps
- Per-validator caps
- Rate limits

#### 17.2 Governance Integration
- Parameter updates
- Emergency pause
- Cap adjustments

**Tests**:
- Cap enforcement tests
- Governance proposal tests
- Emergency scenarios

**Success Criteria**:
- [ ] Caps enforced correctly
- [ ] Governance can update params
- [ ] Emergency pause works

---

### Stage 18: IBC Compatibility (Week 19-20)
**Goal**: Make LST tokens IBC-transferable.

#### 18.1 IBC Integration
- Register LST denoms
- Transfer enablement
- Metadata preservation

**Tests**:
- IBC transfer tests
- Metadata tests
- Multi-chain scenarios

**Success Criteria**:
- [ ] LST tokens transfer via IBC
- [ ] Metadata preserved
- [ ] Exchange rates consistent

## Testing Strategy

### Unit Test Structure
```
tests/
├── stage1/
│   ├── types_test.go
│   └── genesis_test.go
├── stage2/
│   ├── keeper_test.go
│   └── store_test.go
├── stage3/
│   ├── query_test.go
│   └── cli_test.go
└── integration/
    ├── e2e_test.go
    └── scenarios_test.go
```

### Test Execution
```bash
# Run specific stage tests
go test ./x/liquidstaking/tests/stage1/...

# Run integration tests after each stage
go test ./x/liquidstaking/tests/integration/... -tags=stage1

# Full test suite
make test-liquidstaking
```

### Success Metrics
- Each stage has isolated tests
- No stage depends on future stages
- Can deploy partial implementation
- Clear upgrade path between stages

## Rollback Strategy

Each stage can be:
1. **Deployed independently**: Via feature flags
2. **Tested in isolation**: Separate test environments
3. **Rolled back safely**: Without affecting other features
4. **Monitored separately**: Distinct metrics per stage

## Timeline Summary

- **Weeks 1-3**: Foundation (Types, Storage, Queries)
- **Weeks 4-6**: Core Logic (Tokenization, Staking, Token Factory)
- **Weeks 7-11**: EVM Integration (Precompile, LST Tokens)
- **Weeks 12-15**: Advanced Features (Redemption, Rewards)
- **Weeks 16-18**: Risk & Governance
- **Weeks 19-20**: IBC & Polish

Total: 20 weeks for full implementation with 18 independently testable stages.