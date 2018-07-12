# Api

[![Build Status](https://travis-ci.com/MitsukoMegumi/Crypto-Go.svg?branch=master)](https://travis-ci.com/MitsukoMegumi/Crypto-Go) [![Go Report Card](https://goreportcard.com/badge/github.com/mitsukomegumi/crypto-go)](https://goreportcard.com/report/github.com/mitsukomegumi/crypto-go)

## Specifications

### Global

```BASH
http://108.6.212.149:8080/api/
```

#### Request Formats

In callback URLs, ?{}= denotes a parameter. For example, take the following request:

```BASH
http://108.6.212.149:8080/api/accounts/user?username=
```

This request calls for a username, of which would be specified by adding a username after '?username='. All requests will return a JSON-formatted object.

### Authentication

All-Users (base callback):

```BASH
GET: http://108.6.212.149:8080/api/accounts
```

Create an Account:

```BASH
POST: http://108.6.212.149:8080/api/accounts/create?username=?email=?password=
```

Fetch Account Information:

```BASH
GET: http://108.6.212.149:8080/api/accounts/user?username=
```

Update Account Information:

```BASH
POST: http://108.6.212.149:8080/api/accounts/update?username=?email=?newpassword=?oldpassword=
```

Fetch Account Wallet Private Keys:

```BASH
GET: http://108.6.212.149:8080/api/accounts/keys?username=?password=
```

Delete an Account:

```BASH
DELETE: http://108.6.212.149:8080/api/accounts/remove?username=?password=
```

Fetch an Account's Balance:

```BASH
POST: http://108.6.212.149:8080/api/deposit?username=?symbol=
```

#### General Account Request Parameters

```JSON
1. username: username for account
```

```JSON
2. email: specified email for account
```

```JSON
3. password: account password
```

```JSON
4. symbol: trading symbol ("BTC", "LTC", "ETH")
```

```JSON
5. newpassword: desired newpassword (used only in update request)
```

```JSON
5. oldpassword: raw account old password (used only in update request)
```

### Orders

### Route Specifications

Create an Order:

```BASH
POST: http://180.6.212.149:8080/api/orders?pair=?type=?amount=?fillprice=?username=?password=

Note: before creating an order, make sure to fetch the balance of the issuing account
```

Cancel an Order:

```BASH
DELETE: http://108.6.212.149:8080/api/orders?pair=?OrderID=?username=?password=
```

Fetch Order Information:

```BASH
GET: http://108.6.212.149:8080/api/orders/order?pair=?OrderID=
```

Update Order Information:

```BASH
POST: http://108.6.212.149:8080/api/orders/update?pair=?OrderID=?username=?password=?updatedfill=?updatedamount=

Note: before updating an order, make sure to fetch the balance of the issuing account
```

Fetch All Orders for Trading Pair:

```BASH
GET: http://108.6.212.149:8080/api/orders?pair=
```

#### General Order Parameters

```JSON
1. pair: specific trading pair (e.g. "BTC-ETH")
```

```JSON
2. type: string specifying buy or sell (e.g. "BUY")
```

```JSON
3. amount: specifies amount to trade, can have decimal
```

```JSON
4. username: username of user to issue order (e.g. "satoshi")
```

```JSON
5. fillprice: price at which order is to fill
```

```JSON
6. password: password of user to issue order
```

```JSON
7. OrderID: id of order, found under JSON tag "orderid"
```

```JSON
8. updatefill: desired new fill price (used only in update request)
```

```JSON
9. updateamount: desired new amount of coins to trade (used only in update request)
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
