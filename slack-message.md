# Slack Message for Remy

Hey @remy! 👋

Pushed liquid staking research to **`research/lst`** branch - exploring a unique approach for Flora.

## 💡 The Concept

Enable smart contracts to directly access Cosmos staking through EVM precompiles. This creates liquid staking tokens (sPETAL) that auto-compound rewards and handle slashing automatically - no oracles needed.

## 🎯 Key Benefits

- **DeFi Native**: LST tokens work seamlessly in AMMs, lending protocols
- **Capital Efficient**: Users keep earning staking rewards while using tokens in DeFi
- **Risk Managed**: Staged implementation with feature flags for safe rollout

## 📁 What's There

- [Architecture design](https://github.com/flora-evm/flora/blob/research/lst/docs/liquid-staking/architecture/01-overview.md) with precompile approach
- [Implementation strategy](https://github.com/flora-evm/flora/blob/research/lst/docs/liquid-staking/implementation/01-staged-approach.md) breaking it into testable stages
- [Working code examples](https://github.com/flora-evm/flora/tree/research/lst/docs/liquid-staking/examples)

This leverages Flora's unique EVM+Cosmos architecture to solve liquid staking in a way pure EVM or pure Cosmos chains can't.

Thoughts on the approach?