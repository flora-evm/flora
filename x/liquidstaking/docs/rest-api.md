# Liquid Staking REST API Documentation

## Overview

The liquid staking module exposes REST endpoints through the gRPC gateway. All endpoints support JSON responses and follow standard Cosmos SDK REST conventions.

## Base URL

```
http://localhost:1317
```

Replace with your node's API endpoint address.

## Query Endpoints

### Get Module Parameters

Retrieve the current liquid staking module parameters.

```http
GET /cosmos/liquidstaking/v1/params
```

#### Response

```json
{
  "params": {
    "enabled": true,
    "min_liquid_stake_amount": "1000000",
    "global_liquid_staking_cap": "0.250000000000000000",
    "validator_liquid_cap": "0.100000000000000000"
  }
}
```

### Get Tokenization Record

Retrieve a specific tokenization record by ID.

```http
GET /cosmos/liquidstaking/v1/records/{record_id}
```

#### Parameters

- `record_id` (path, required): The tokenization record ID

#### Response

```json
{
  "record": {
    "id": "1",
    "validator_address": "floravaloper1abc...",
    "owner": "flora1xyz...",
    "shares_denomination": "shares/floravaloper1abc...",
    "liquid_staking_token_denom": "liquidstake/floravaloper1abc.../1",
    "shares_amount": "1000000.000000000000000000",
    "status": "TOKENIZATION_RECORD_STATUS_ACTIVE",
    "created_at": "2024-01-15T10:30:00Z",
    "redeemed_at": null
  }
}
```

#### Error Response

```json
{
  "code": 5,
  "message": "tokenization record not found: 999",
  "details": []
}
```

### List All Tokenization Records

Retrieve all tokenization records with optional pagination.

```http
GET /cosmos/liquidstaking/v1/records
```

#### Query Parameters

- `pagination.key` (optional): The key for the next page
- `pagination.offset` (optional): Number of items to skip
- `pagination.limit` (optional): Maximum number of items to return
- `pagination.count_total` (optional): Set to true to return total count
- `pagination.reverse` (optional): Reverse the result order

#### Response

```json
{
  "records": [
    {
      "id": "1",
      "validator_address": "floravaloper1abc...",
      "owner": "flora1xyz...",
      "shares_denomination": "shares/floravaloper1abc...",
      "liquid_staking_token_denom": "liquidstake/floravaloper1abc.../1",
      "shares_amount": "1000000.000000000000000000",
      "status": "TOKENIZATION_RECORD_STATUS_ACTIVE",
      "created_at": "2024-01-15T10:30:00Z",
      "redeemed_at": null
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

### Get Records by Validator

Retrieve all tokenization records for a specific validator.

```http
GET /cosmos/liquidstaking/v1/validators/{validator_address}/records
```

#### Parameters

- `validator_address` (path, required): The validator's address

#### Response

```json
{
  "records": [
    {
      "id": "1",
      "validator_address": "floravaloper1abc...",
      "owner": "flora1xyz...",
      "shares_denomination": "shares/floravaloper1abc...",
      "liquid_staking_token_denom": "liquidstake/floravaloper1abc.../1",
      "shares_amount": "1000000.000000000000000000",
      "status": "TOKENIZATION_RECORD_STATUS_ACTIVE",
      "created_at": "2024-01-15T10:30:00Z",
      "redeemed_at": null
    }
  ]
}
```

### Get Records by Owner

Retrieve all tokenization records owned by a specific address.

```http
GET /cosmos/liquidstaking/v1/owners/{owner_address}/records
```

#### Parameters

- `owner_address` (path, required): The owner's address

#### Response

```json
{
  "records": [
    {
      "id": "1",
      "validator_address": "floravaloper1abc...",
      "owner": "flora1xyz...",
      "shares_denomination": "shares/floravaloper1abc...",
      "liquid_staking_token_denom": "liquidstake/floravaloper1abc.../1",
      "shares_amount": "1000000.000000000000000000",
      "status": "TOKENIZATION_RECORD_STATUS_ACTIVE",
      "created_at": "2024-01-15T10:30:00Z",
      "redeemed_at": null
    }
  ]
}
```

### Get Total Liquid Staked

Retrieve the total amount of tokens liquid staked across all validators.

```http
GET /cosmos/liquidstaking/v1/total_liquid_staked
```

#### Response

```json
{
  "amount": "5000000"
}
```

### Get Validator Liquid Staked

Retrieve the amount of tokens liquid staked for a specific validator.

```http
GET /cosmos/liquidstaking/v1/validators/{validator_address}/liquid_staked
```

#### Parameters

- `validator_address` (path, required): The validator's address

#### Response

```json
{
  "amount": "1000000"
}
```

## Transaction Endpoints

Transactions are submitted through the standard Cosmos SDK transaction endpoints. The liquid staking module messages can be included in transactions sent to:

```http
POST /cosmos/tx/v1beta1/txs
```

### Message Types

#### TokenizeShares

```json
{
  "@type": "/cosmos.liquidstaking.v1.MsgTokenizeShares",
  "delegator_address": "flora1xyz...",
  "validator_address": "floravaloper1abc...",
  "shares": {
    "denom": "stake",
    "amount": "1000000.000000000000000000"
  },
  "owner_address": "flora1xyz..."
}
```

#### RedeemTokens

```json
{
  "@type": "/cosmos.liquidstaking.v1.MsgRedeemTokens",
  "owner_address": "flora1xyz...",
  "amount": {
    "denom": "liquidstake/floravaloper1abc.../1",
    "amount": "1000000"
  }
}
```

### Complete Transaction Example

```json
{
  "body": {
    "messages": [
      {
        "@type": "/cosmos.liquidstaking.v1.MsgTokenizeShares",
        "delegator_address": "flora1xyz...",
        "validator_address": "floravaloper1abc...",
        "shares": {
          "denom": "stake",
          "amount": "1000000.000000000000000000"
        },
        "owner_address": "flora1xyz..."
      }
    ],
    "memo": "Tokenizing delegation shares",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [
        {
          "denom": "flora",
          "amount": "1000"
        }
      ],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

## Status Codes

- `200 OK`: Successful query
- `400 Bad Request`: Invalid parameters
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Common Response Fields

### Tokenization Record Status

- `TOKENIZATION_RECORD_STATUS_UNSPECIFIED`: Unspecified status (should not occur)
- `TOKENIZATION_RECORD_STATUS_ACTIVE`: Record is active, tokens are in circulation
- `TOKENIZATION_RECORD_STATUS_REDEEMED`: Tokens have been redeemed, record is archived

### Timestamps

All timestamps are in RFC3339 format (e.g., `2024-01-15T10:30:00Z`).

### Amount Fields

- Numeric amounts (e.g., `amount`) are returned as strings to preserve precision
- Share amounts include 18 decimal places of precision
- Token amounts are integers without decimals

## Pagination

Endpoints that return lists support pagination through query parameters:

```http
GET /cosmos/liquidstaking/v1/records?pagination.limit=10&pagination.offset=20
```

The response includes a `pagination` object with:
- `next_key`: Base64 encoded key for the next page
- `total`: Total count of items (if `count_total` was set to true)

## Error Handling

Errors follow the gRPC error model with additional details:

```json
{
  "code": 3,
  "message": "invalid argument: validator address cannot be empty",
  "details": [
    {
      "@type": "type.googleapis.com/google.rpc.BadRequest",
      "field_violations": [
        {
          "field": "validator_address",
          "description": "cannot be empty"
        }
      ]
    }
  ]
}
```

Common error codes:
- `3`: Invalid argument
- `5`: Not found
- `13`: Internal error

## Rate Limiting

REST endpoints may be subject to rate limiting. Check response headers:
- `X-RateLimit-Limit`: Maximum requests per window
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Unix timestamp when the window resets

## Authentication

Public query endpoints do not require authentication. Transaction submission requires a properly signed transaction with valid account credentials.

## CORS

The REST API supports CORS for browser-based applications. Ensure your node is configured to allow CORS from your application's domain.

## Examples

### cURL Examples

Get module parameters:
```bash
curl -X GET "http://localhost:1317/cosmos/liquidstaking/v1/params"
```

Get a specific record:
```bash
curl -X GET "http://localhost:1317/cosmos/liquidstaking/v1/records/1"
```

Get records by owner with pagination:
```bash
curl -X GET "http://localhost:1317/cosmos/liquidstaking/v1/owners/flora1xyz.../records?pagination.limit=10"
```

### JavaScript Example

```javascript
// Fetch module parameters
async function getParams() {
  const response = await fetch('http://localhost:1317/cosmos/liquidstaking/v1/params');
  const data = await response.json();
  return data.params;
}

// Get tokenization record
async function getRecord(recordId) {
  const response = await fetch(`http://localhost:1317/cosmos/liquidstaking/v1/records/${recordId}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch record: ${response.statusText}`);
  }
  const data = await response.json();
  return data.record;
}

// Get total liquid staked with error handling
async function getTotalLiquidStaked() {
  try {
    const response = await fetch('http://localhost:1317/cosmos/liquidstaking/v1/total_liquid_staked');
    const data = await response.json();
    return BigInt(data.amount);
  } catch (error) {
    console.error('Failed to fetch total liquid staked:', error);
    throw error;
  }
}
```

## WebSocket Support

For real-time updates, use the Tendermint WebSocket endpoint to subscribe to events:

```javascript
const ws = new WebSocket('ws://localhost:26657/websocket');

ws.onopen = () => {
  // Subscribe to liquid staking events
  ws.send(JSON.stringify({
    jsonrpc: '2.0',
    method: 'subscribe',
    id: '1',
    params: {
      query: "tm.event='Tx' AND liquidstaking.action='tokenize_shares'"
    }
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('New tokenization event:', data);
};
```

## Further Resources

- [Cosmos SDK REST Documentation](https://docs.cosmos.network/main/core/grpc_rest)
- [gRPC Gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Liquid Staking Module CLI Documentation](./cli-usage.md)