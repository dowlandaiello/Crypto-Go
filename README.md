# Api

## Specifications

### Global

```BASH
http://108.6.212.149:8080/api/
```

#### Request Formats

In callback URLs, ${} denotes a parameter. For example, take the following request:

```BASH
http://108.6.212.149:8080/api/accounts/${username}
```

This request calls for a username, of which would be specified by replacing '${username}' with a username. All requests will return a JSON-formatted object.

### Authentication

All-Users (base callback):

```BASH
GET: http://108.6.212.149:8080/api/accounts
```

Create an Account:

```BASH
POST: http://108.6.212.149:8080/api/accounts/${username}/${email}/${password}
```

Fetch information for an Account:

```BASH
GET: http://108.6.212.149:8080/api/accounts/${username}
```

Delete an Account:

```BASH
DELETE: http://108.6.212.149:8080/api/accounts/${username}/${password}
```

### Orders

Create an Order:

```BASH
POST: http://180.6.212.149:8080/api/orders/${pair}/${type}/${amount}/${user}/${pass}
```

```JSON
{request}: returns JSON order object
```

Cancel an Order:

```BASH
DELETE: http://108.6.212.149:8080/api/accounts/${pair}/${id}/${user}/${pass}
```

Fetch information on an Order:

```BASH
GET: http://108.6.212.149:8080/api/orders/${pair}/${orderid}
```

```JSON
${OrderID}: can be retrieved via account orders or on creation
```

## Definitions

### Exchange

1. Global clocks
2. Trading pair data

### System

1. Downtime
2. System time
3. Server ping

### Individual Trading Pairs

1. Volume
2. Latest Orders (amount configurable)
3. Open
4. Close
5. Day High
6. Day Low
7. Trade Offers (amount configurable)
8. Create Order

### Individual Orders

1. Filled
2. Issuance Time
3. Amount
4. Fill Time (if filled)
5. Order Type
6. Order Fee
7. Trading Pair

### Accounts (only accessible via user token)

1. Username
2. Email
3. Hashed Password
4. Hash ID
5. Wallet Balances
6. Wallet Addresses

## Routes

### Exchange Routes

1. Get - global clock times
2. Get - Trading pair data

### System Routes

1. Get - Downtime
2. Get - System time
3. Get - Server ping

### Individual Trading Pair Routes

1. Get - Volume
2. Get - Latest Orders (amount configurable)
3. Get - Open
4. Get - Close
5. Get - Day High
6. Get - Day Low
7. Get - Trade Offers (amount configurable)
8. Post - Create Order

### Individual Order Routes

1. Get - Filled
2. Get - Issuance Time
3. Get - Amount
4. Get - Fill Time (if filled)
5. Get - Order Type
6. Get - Order Fee
7. Get - Trading Pair

### Account Routes (only accessible via user token)

1. Get - Username
2. Get - Email
3. Get - Hashed Password
4. Get - Hash ID
5. Get - Wallet Balances
6. Get - Wallet Addresses
