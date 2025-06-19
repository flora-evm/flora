# Makefile for liquid staking precompile

.PHONY: generate-abi
generate-abi:
	@echo "Generating ABI bindings..."
	@go run generate_abi.go
	@echo "ABI bindings generated successfully"

.PHONY: generate-solidity
generate-solidity:
	@echo "Solidity interface is manually maintained in ILiquidStaking.sol"

.PHONY: test
test:
	@echo "Running precompile tests..."
	@go test -v ./...

.PHONY: all
all: generate-abi test

.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	@rm -f abi_generated.go

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  generate-abi     Generate Go bindings from ABI JSON"
	@echo "  test            Run precompile tests"
	@echo "  all             Generate ABI and run tests"
	@echo "  clean           Remove generated files"
	@echo "  help            Show this help message"