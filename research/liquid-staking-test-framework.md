# Liquid Staking Test Framework

## Overview
This framework ensures each implementation stage can be tested in complete isolation with clear boundaries and mock dependencies.

## Feature Flag System

### Configuration
```go
// x/liquidstaking/types/feature_flags.go
type FeatureFlags struct {
    // Stage 1-3: Always enabled (basic infrastructure)
    
    // Stage 4: Simple tokenization
    EnableTokenization bool `json:"enable_tokenization"`
    
    // Stage 5: Staking integration
    EnableStakingIntegration bool `json:"enable_staking_integration"`
    
    // Stage 6: Token factory
    EnableTokenFactory bool `json:"enable_token_factory"`
    
    // Stage 7-9: Precompile
    EnablePrecompile bool `json:"enable_precompile"`
    
    // Stage 10-11: LST tokens
    EnableLSTTokens bool `json:"enable_lst_tokens"`
    
    // Stage 12: Redemption
    EnableRedemption bool `json:"enable_redemption"`
    
    // Stage 13-15: Rewards
    EnableRewards bool `json:"enable_rewards"`
    EnableAutoCompounding bool `json:"enable_auto_compounding"`
    
    // Stage 16: Slashing
    EnableSlashing bool `json:"enable_slashing"`
    
    // Stage 17: Advanced features
    EnableCaps bool `json:"enable_caps"`
    EnableGovernance bool `json:"enable_governance"`
    
    // Stage 18: IBC
    EnableIBC bool `json:"enable_ibc"`
}

// Default flags for each stage
func GetStageFlags(stage int) FeatureFlags {
    flags := FeatureFlags{}
    
    if stage >= 4 {
        flags.EnableTokenization = true
    }
    if stage >= 5 {
        flags.EnableStakingIntegration = true
    }
    // ... etc
    
    return flags
}
```

### Usage in Keeper
```go
func (k Keeper) TokenizeShares(ctx sdk.Context, msg *MsgTokenizeShares) error {
    flags := k.GetFeatureFlags(ctx)
    
    if !flags.EnableTokenization {
        return ErrFeatureDisabled
    }
    
    // Basic tokenization logic
    record := k.createBasicRecord(msg)
    
    if flags.EnableStakingIntegration {
        // Add staking validation
        if err := k.validateDelegation(ctx, msg); err != nil {
            return err
        }
    }
    
    if flags.EnableTokenFactory {
        // Mint tokens
        if err := k.mintLSTTokens(ctx, record); err != nil {
            return err
        }
    }
    
    return nil
}
```

## Mock Implementations

### Stage-Specific Mocks
```go
// x/liquidstaking/testutil/mocks.go

// Mock for stages before staking integration
type MockStakingKeeper struct {
    delegations map[string]stakingtypes.Delegation
    validators  map[string]stakingtypes.Validator
}

func NewMockStakingKeeper() *MockStakingKeeper {
    return &MockStakingKeeper{
        delegations: make(map[string]stakingtypes.Delegation),
        validators:  make(map[string]stakingtypes.Validator),
    }
}

func (m *MockStakingKeeper) GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, bool) {
    key := delAddr.String() + valAddr.String()
    del, found := m.delegations[key]
    return del, found
}

// Mock for stages before token factory
type MockTokenFactoryKeeper struct {
    denoms map[string]bool
    minted map[string]sdk.Coins
}

func (m *MockTokenFactoryKeeper) CreateDenom(ctx sdk.Context, creator, subdenom string) (string, error) {
    denom := fmt.Sprintf("factory/%s/%s", creator, subdenom)
    m.denoms[denom] = true
    return denom, nil
}

// Mock for stages before precompile
type MockEVMKeeper struct {
    contracts map[common.Address][]byte
}

func (m *MockEVMKeeper) DeployContract(ctx sdk.Context, code []byte) (common.Address, error) {
    addr := common.HexToAddress(fmt.Sprintf("0x%x", len(m.contracts)))
    m.contracts[addr] = code
    return addr, nil
}
```

## Test Harness

### Stage Test Suite
```go
// x/liquidstaking/tests/harness.go

type StageTestSuite struct {
    suite.Suite
    
    ctx    sdk.Context
    keeper keeper.Keeper
    stage  int
    
    // Mocks
    stakingKeeper      *MockStakingKeeper
    tokenFactoryKeeper *MockTokenFactoryKeeper
    evmKeeper          *MockEVMKeeper
}

func (s *StageTestSuite) SetupTest() {
    // Create keeper with stage-specific configuration
    s.keeper = s.createKeeperForStage(s.stage)
    s.ctx = s.createContext()
}

func (s *StageTestSuite) createKeeperForStage(stage int) keeper.Keeper {
    flags := types.GetStageFlags(stage)
    
    k := keeper.NewKeeper(
        s.storeKey,
        s.cdc,
        flags,
    )
    
    // Inject stage-appropriate dependencies
    if stage < 5 {
        k.SetStakingKeeper(s.stakingKeeper)
    } else {
        k.SetStakingKeeper(s.app.StakingKeeper) // Real keeper
    }
    
    if stage < 6 {
        k.SetTokenFactoryKeeper(s.tokenFactoryKeeper)
    } else {
        k.SetTokenFactoryKeeper(s.app.TokenFactoryKeeper)
    }
    
    return k
}
```

## Isolated Test Cases

### Stage 1-3: Foundation Tests
```go
// x/liquidstaking/tests/stage1_test.go

func TestStage1_Types(t *testing.T) {
    // Test type validation
    record := types.TokenizationRecord{
        Id:        1,
        Validator: "cosmosvaloper1...",
        Owner:     "cosmos1...",
    }
    
    require.NoError(t, record.Validate())
}

func TestStage2_StateManagement(t *testing.T) {
    suite.Run(t, &StageTestSuite{stage: 2})
}

func (s *StageTestSuite) TestSetGetRecord() {
    record := types.TokenizationRecord{Id: 1}
    s.keeper.SetTokenizationRecord(s.ctx, record)
    
    retrieved, found := s.keeper.GetTokenizationRecord(s.ctx, 1)
    s.Require().True(found)
    s.Require().Equal(record, retrieved)
}
```

### Stage 4: Tokenization Tests
```go
// x/liquidstaking/tests/stage4_test.go

func TestStage4_BasicTokenization(t *testing.T) {
    suite.Run(t, &StageTestSuite{stage: 4})
}

func (s *StageTestSuite) TestTokenizeShares_NoIntegration() {
    // Setup mock delegation
    s.stakingKeeper.SetDelegation(delegation)
    
    msg := &types.MsgTokenizeShares{
        DelegatorAddress: "cosmos1...",
        ValidatorAddress: "cosmosvaloper1...",
        Amount:           sdk.NewCoin("stake", sdk.NewInt(1000)),
    }
    
    res, err := s.keeper.TokenizeShares(s.ctx, msg)
    s.Require().NoError(err)
    s.Require().Equal(uint64(1), res.RecordId)
    
    // Verify record created
    record, found := s.keeper.GetTokenizationRecord(s.ctx, 1)
    s.Require().True(found)
    s.Require().Equal(msg.DelegatorAddress, record.Owner)
}
```

### Stage 5: Integration Tests
```go
// x/liquidstaking/tests/stage5_test.go

func TestStage5_StakingIntegration(t *testing.T) {
    suite.Run(t, &StageTestSuite{stage: 5})
}

func (s *StageTestSuite) TestTokenizeShares_WithStaking() {
    // Use real staking keeper
    // Create actual delegation
    s.app.StakingKeeper.Delegate(s.ctx, delAddr, sdk.NewInt(1000), validator)
    
    msg := &types.MsgTokenizeShares{
        DelegatorAddress: delAddr.String(),
        ValidatorAddress: valAddr.String(),
        Amount:           sdk.NewCoin("stake", sdk.NewInt(500)),
    }
    
    res, err := s.keeper.TokenizeShares(s.ctx, msg)
    s.Require().NoError(err)
    
    // Verify delegation reduced
    del, _ := s.app.StakingKeeper.GetDelegation(s.ctx, delAddr, valAddr)
    s.Require().Equal(sdk.NewInt(500), del.Shares.TruncateInt())
}
```

## Integration Test Scenarios

### Progressive Integration
```go
// x/liquidstaking/tests/integration/progressive_test.go

func TestProgressiveIntegration(t *testing.T) {
    // Test each stage builds on previous
    for stage := 1; stage <= 18; stage++ {
        t.Run(fmt.Sprintf("Stage%d", stage), func(t *testing.T) {
            suite := &StageTestSuite{stage: stage}
            
            // Run all tests up to current stage
            for i := 1; i <= stage; i++ {
                suite.RunStageTests(i)
            }
        })
    }
}
```

### Rollback Testing
```go
func TestStageRollback(t *testing.T) {
    // Start with stage 10
    app := CreateTestApp(10)
    
    // Create some state
    app.LiquidStakingKeeper.TokenizeShares(ctx, msg)
    
    // Rollback to stage 5
    app.UpdateFeatureFlags(5)
    
    // Verify basic functions still work
    record, found := app.LiquidStakingKeeper.GetTokenizationRecord(ctx, 1)
    require.True(t, found)
    
    // Verify advanced features disabled
    err := app.LiquidStakingKeeper.MintLSTTokens(ctx, record)
    require.ErrorIs(t, err, ErrFeatureDisabled)
}
```

## Benchmark Tests

### Stage Performance Benchmarks
```go
// x/liquidstaking/tests/benchmarks/stage_bench_test.go

func BenchmarkTokenization_Stage4(b *testing.B) {
    suite := &StageTestSuite{stage: 4}
    suite.SetupTest()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        suite.keeper.TokenizeShares(suite.ctx, msg)
    }
}

func BenchmarkTokenization_Stage11(b *testing.B) {
    suite := &StageTestSuite{stage: 11}
    suite.SetupTest()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        suite.keeper.TokenizeShares(suite.ctx, msg)
    }
}

// Compare performance impact of each stage
func BenchmarkStageComparison(b *testing.B) {
    stages := []int{4, 5, 6, 11, 15, 18}
    
    for _, stage := range stages {
        b.Run(fmt.Sprintf("Stage%d", stage), func(b *testing.B) {
            suite := &StageTestSuite{stage: stage}
            suite.SetupTest()
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                suite.keeper.TokenizeShares(suite.ctx, msg)
            }
        })
    }
}
```

## CI/CD Integration

### GitHub Actions Workflow
```yaml
# .github/workflows/liquidstaking-stages.yml
name: Liquid Staking Stage Tests

on:
  pull_request:
    paths:
      - 'x/liquidstaking/**'

jobs:
  stage-tests:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        stage: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18]
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Run Stage ${{ matrix.stage }} Tests
        run: |
          go test ./x/liquidstaking/tests/stage${{ matrix.stage }}_test.go -v
      
      - name: Run Integration Tests up to Stage ${{ matrix.stage }}
        run: |
          go test ./x/liquidstaking/tests/integration/... -tags=stage${{ matrix.stage }}
      
      - name: Benchmark Stage ${{ matrix.stage }}
        run: |
          go test -bench=Stage${{ matrix.stage }} ./x/liquidstaking/tests/benchmarks/...
```

### Local Testing Script
```bash
#!/bin/bash
# scripts/test-liquidstaking-stage.sh

STAGE=${1:-1}

echo "Testing Liquid Staking Stage $STAGE"

# Run stage-specific tests
echo "Running Stage $STAGE unit tests..."
go test ./x/liquidstaking/tests/stage${STAGE}_test.go -v

# Run integration tests
echo "Running integration tests up to Stage $STAGE..."
go test ./x/liquidstaking/tests/integration/... -tags=stage${STAGE}

# Run benchmarks
echo "Running Stage $STAGE benchmarks..."
go test -bench=Stage${STAGE} ./x/liquidstaking/tests/benchmarks/...

# Check coverage
echo "Checking test coverage..."
go test -cover ./x/liquidstaking/... -tags=stage${STAGE}
```

## Stage Validation Checklist

### For Each Stage
- [ ] Unit tests pass in isolation
- [ ] Integration tests pass with previous stages
- [ ] No dependencies on future stages
- [ ] Feature flags properly disable/enable
- [ ] Benchmarks show acceptable performance
- [ ] Can rollback to previous stage
- [ ] Documentation updated
- [ ] Migration path defined

### Stage Sign-off Template
```markdown
## Stage X Sign-off

**Date**: YYYY-MM-DD
**Developer**: @username
**Reviewer**: @reviewer

### Checklist
- [ ] All tests passing
- [ ] Benchmarks acceptable
- [ ] Feature flags working
- [ ] Documentation complete
- [ ] Security review done

### Test Results
- Unit Tests: X/X passing
- Integration Tests: X/X passing
- Coverage: XX%
- Benchmark: XXXns/op

### Notes
[Any additional notes or concerns]

### Approval
- [ ] Developer sign-off
- [ ] Reviewer sign-off
- [ ] Ready for next stage
```