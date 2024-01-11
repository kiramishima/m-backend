# BondApp Backend

-----

# **Stack**

- Router: [Chi 🚀](https://github.com/go-chi/chi)
- Logger: [Zap ⚡](https://github.com/uber-go/zap)
- Mocks: [gomock 💀](https://github.com/uber-go/mock)
- Asserts: [testify 💀](https://github.com/stretchr/testify)
- DI: [fx 🤖](https://github.com/uber-go/fx)
- Deploy: [Docker 🐳](https://www.docker.com)
- Database: [MariaDB](https://mariadb.org/download/?t=mariadb)

## **Deploy with Docker**
- The `scripts/build-container.sh` file is used to build the docker container.
- The `scripts/run-container.sh` file is used to run/execute the docker container generate with `scripts/build-container.sh`.
- Service run in port `8080`.
- For change the default port (8080), provide the environment variable `PORT`. If you change that in `scripts/build-container.sh`, you'll need change also in `scripts/run-container.sh`.

## **Deploy with Compose**
To deploy in composer, create a basic docker-compose.yaml file like:

```yaml
version: "3"
services:
   database:
      image: mariadb:latest
      hostname: database
      ports:
        - 3306:3306
      env_file:
        - ./.env
   nats_server:
      image: nats:2.10.7-alpine3.18
      ports:
         - 4222:4222
      env_file:
        - ./config
      volumes:
        - ./test.sh:/opt/test.sh
   redis_server:
      image: redis:7.2.4-alpine
      ports:
        - 6379:6379
   bondsapp:
      image: bondsapp-backend
      ports:
        - 8080:8080
      env_file:
        - ./.env
```

### Run the docker containers

Run the docker containers using docker-compose

```
docker compose up -d
```

# Deploy in local
- Install [golang](https://golang.org/dl)
- Install [Task CLI](https://taskfile.dev/) for executing the task of the taskfile.
    - The command `task run` launch the service. Default port is 8080

---
## Summary of API Specification

### Endpoint: Sign-In

* Path: `/v1/auth/sign-in`
* Method: `POST`
* Payload: {email: string, password: string}
* Response: JSON Response.

Description:

Takes in a JSON data for authenticate an user. It returns a token authentication or error message.

Example of Responses:
```json
{ "token": "7fb1377bb22349d9a31a-5a02701dd310" }
```

```json
{ "error": "Wrong password" }
```

### Endpoint: Sign-Up

* Path: `/v1/auth/sign-up`
* Method: `POST`
* Payload: {email: string, password: string, name: string}
* Payload Rules:
  * Name: gte=6, alphanum, required
  * Email: gte=6, email, required
  * Password: gte=6, alphanum, required 
* Response: JSON Response.

Description:

Takes in a JSON data for register a new user. It returns a success message or error message.

Example of Responses:
```json
{ "message": "Something" }
```

```json
{ "error": "Wrong password" }
```

### Endpoint: ListBonds

* Path: `/v1/bonds`
* Method: `GET`
* Auth: Bearer Token
* Response: JSON Response.

Description:

Return the list of user bonds. Required a authentication token

Example of Responses:
```json
{ 
  "data": [
    {
      "id": 1,
      "bond_id": "35as43a-23as4d32a-2s22a-1s22a",
      "name": "AX23",
      "price": 1500.0000,
      "number": 200,
      "currency": 1,
      "created_by": "solid_snake",
      "created_by_id": 1,
      "on_sale": false,
      "is_owner": true,
      "status": "on_hold",
      "created_at": "10/01/2024 13:26:25",
      "updated_at": ""
    },
    {
      "id": 2,
      "bond_id": "35as43a-23as4d32a-2s22a-1s22a",
      "name": "AX24",
      "price": 500.0000,
      "number": 400,
      "currency": 1,
      "created_by": "solid_snake",
      "created_by_id": 1,
      "on_sale": true,
      "is_owner": true,
      "status": "on_sale",
      "created_at": "10/01/2024 13:26:25",
      "updated_at": ""
    }
  ]
}
```

```json
{ "data": [] }
```

### Endpoint: CreateBond

* Path: `/v1/bonds`
* Method: `POST`
* Payload: {name: string, number: int, price: float, currency_id: int}
* Payload Rules:
  * name: Length >= 4
  * number: Min: 1, Max: 10000
  * price: Min: 0, Max: 1000000000.0000
  * currency_id: Default: 1
* Response: JSON Response.

Description:

Takes in a JSON data for create a new bond. Default status is `on_hold`.
Required a authentication token.

Example of Responses:
```json
{ "message": "Success bond created" }
```

```json
{ "error": "fail creating a bond" }
```

### Endpoint: UpdateBond

* Path: `/v1/bonds/{id}`
* Method: `PATCH`
* Payload Rules:
  * name: Length >= 4
  * number: Min: 1, Max: 10000
  * price: Min: 0, Max: 1000000000.0000
  * currency_id: Default: 1
  * status: on_hold | on_sale
* Response: JSON Response.

Description:

To update a user bond.
Takes a JSON data for update the bond.
Required a authentication token.

Example of Responses:
```json
{ "message": "Success bond created" }
```

```json
{ "error": "fail creating a bond" }
```

### Endpoint: DeleteBond

* Path: `/v1/bonds/{id}`
* Method: `DELETE`
* Response: JSON Response.

Description:

To Delete a user bond.
Required a authentication token.

Example of Responses:
```json
{ "message": "The bond was deleted" }
```

```json
{ "error": "fail deleting a bond" }
```

### Endpoint: MarketBondList

* Path: `/v1/market`
* Method: `GET`
* Auth: Bearer Token
* Response: JSON Response.

Description:

Return the list of bonds on sale. Required a authentication token

Example of Responses:
```json
{ 
  "data": [
    {
      "id": 1,
      "bond_uuid": "35as43a-23as4d32a-2s22a-1s22a",
      "name": "AX23",
      "price": 1500.0000,
      "available": 200,
      "currency": 1,
      "created_by": "seller_1",
      "created_by_id": 1,
      "is_owner": false,
      "status": "available",
      "created_at": "10/01/2024 13:26:25",
      "updated_at": ""
    },
    {
      "id": 2,
      "bond_id": "35as43a-23as4d32a-2s22a-1s22a",
      "name": "AX24",
      "price": 500.0000,
      "available": 0,
      "currency": 1,
      "created_by": "seller_2",
      "created_by_id": 1,
      "is_owner": false,
      "status": "bought",
      "created_at": "10/01/2024 13:26:25",
      "updated_at": ""
    }
  ]
}
```

```json
{ "data": [] }
```

### Endpoint: GetMarketBondByID

* Path: `/v1/market/{id`
* Method: `GET`
* Auth: Bearer Token
* Response: JSON Response.

Description:

Return a specify bond by provided id. Required a authentication token

Example of Responses:
```json
{ 
  "data": {
      "id": 1,
      "bond_uuid": "35as43a-23as4d32a-2s22a-1s22a",
      "name": "AX23",
      "price": 1500.0000,
      "available": 200,
      "currency": 1,
      "created_by": "seller_1",
      "created_by_id": 1,
      "is_owner": false,
      "status": "available",
      "created_at": "10/01/2024 13:26:25",
      "updated_at": ""
    }
}
```

```json
{ "error": "error message" }
```

### Endpoint: MarketBuyBond

* Path: `/v1/market/{id}/buy`
* Method: `POST`
* Payload: {seller_id: int, order: int}
* Payload Rules:
  * order: Min: 1, Max: 10000
* Response: JSON Response.

Description:

To buy a bond available in the market
Required a authentication token.

Example of Responses:
```json
{ "message": "Success. Your bought is in process." }
```

```json
{ "error": "fail buying a bond" }
```

### Endpoint: MarketSellBond

* Path: `/v1/market/sell`
* Method: `POST`
* Payload Rules:
  * bond_id: Required
  * num_sell: Min: 1, Max: 10000
* Response: JSON Response.

Description:

To sell a user bond.
Takes a JSON data for update the bond.
Required a authentication token.

Example of Responses:
```json
{ "message": "Success. Your bond is available to sell in the market" }
```

```json
{ "error": "failed putting the bond on sale" }
```

---

