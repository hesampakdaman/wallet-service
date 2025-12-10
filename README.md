# wallet-service
In this application we provide each player to create their own wallets. Any wallet's balance can be read and updated.

## Known limitations
- Amounts and balances are integers. We don't provide functionality to work with fractions.
- No idempotency checks on endpoints. In particular, retries would be treated as a new request.

## Architecture
- **Core**: Business logic and entities.
- **Service**: Orchestration logic.
- **Adapters**: External dependencies.

## Usage

### Running locally
Ensure you have set the environment variable `DATABASE_URL`. Then run

```sh
make run
```

### Docker
Make sure you have installed Docker and Docker compose correctly. Then run

```sh
make docker-start
```

You can stop the service

```sh
make docker-stop
```

and clean up docker related resources (including volume)

```sh
make docker-clean
```

## API examples

Assuming the app is running on `localhost:8080`.

### Create wallet
```sh
curl -X POST http://localhost:8080/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "player_id": "11111111-1111-1111-1111-111111111111",
    "initial_balance": 100
  }'
# => {"wallet_id": "<wallet-id>"}
```

### Adds to the balance
```sh
curl -X PATCH http://localhost:8080/wallet/<wallet-id> \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 50,
    "type": "deposit"
  }'
```

### Subtract from the balance (withdraw)
```sh
curl -X PATCH http://localhost:8080/wallet/<wallet-id> \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 25,
    "type": "withdraw"
  }'
```

### Read balance
```sh
curl http://localhost:8080/wallet/<wallet-id>
```
