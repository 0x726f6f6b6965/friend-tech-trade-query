# friend-tech-trade-query-api

## Requirement

1. As a user I need a function to query the Friend.Tech trading record from web3.

2. To reduce the latency of response, the function need to store data in the cache. 

## Request
```cmd
curl --location 'https://lnirbq5gi0.execute-api.ap-northeast-1.amazonaws.com/query' \
    --header 'Content-Type: application/json' \
    --data '{
        "tx_hash": "0xc3100e7e0bb4f89f91d2d9da6636689601f4447423562e1910e8668a5b78a987"
    }'
```

## Response 
```json
{
    "trader": "0x2B5AB508DffC087232f4df475aE29c4B7a6Aa918",
    "subject": "0x2Ff20e30D147de328a7Af12CE2c7F2520207805e",
    "share_amount": "1",
    "eth_amount": "2250000000000000",
    "protocol_eth_amount": "112500000000000",
    "subject_eth_amount": "112500000000000",
    "supply": "6"
}
```

## Sequence

```mermaid
sequenceDiagram
    actor Client
    participant API Gateway
    participant Lambda Function
    participant Dynamodb
    participant Friend.Tech Contract
    
    Client ->> API Gateway: POST a request with tx_hash
    API Gateway ->> Lambda Function: Pass the request
    Lambda Function ->> Dynamodb: Find the cache data by tx_hash
    Dynamodb -->> Lambda Function: Return the cache data
    alt cache data exist:
        Lambda Function -->> API Gateway: Return the cache as Response
        API Gateway -->> Client: Return the trading record
    else:
        Lambda Function ->> Friend.Tech Contract: Find the trading record
        Friend.Tech Contract -->> Lambda Function: Return the trading record
        Lambda Function -->> API Gateway: Return the trading record
        API Gateway -->> Client: Return the trading record
    end
```