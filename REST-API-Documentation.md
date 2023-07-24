# REST API Documentation

### [GET] `http://localhost:4000/user`

- Request Header
  - Key: 'Authorization'
  - Value: 'Bearer <JWT_Token>'

```js
// Response
{
    "username": "Alice",
    "mobile_number": "+65 90399012",
    "currency": "SGD",
    "balance": 60000.00
}
```

| Error Code | Status                | Description                                       |
| ---------- | --------------------- | ------------------------------------------------- |
| 200        | OK                    | Successfully retrieve user details.               |
| 500        | Internal Server Error | Unexpected exception occurred on the server-side. |

---

### [GET] `http://localhost:4000/beneficiaries`

- Request Header
  - Key: 'Authorization'
  - Value: 'Bearer <JWT_Token>'

```js
// Response:
{
    "beneficiaries": [
        {
            "beneficiary_id": 2,
            "beneficiary_name": "Bob",
            "mobile_number": "+65 89230122",
            "currency": "SGD"
        },
        {
            "beneficiary_id": 3,
            "beneficiary_name": "Charlie",
            "mobile_number": "+1 555-123-4567",
            "currency": "USD"
        },
        {
            "beneficiary_id": 4,
            "beneficiary_name": "David",
            "mobile_number": "+49 1234567890",
            "currency": "EUR"
        }
    ]
}
```

| Error Code | Status                | Description                                        |
| ---------- | --------------------- | -------------------------------------------------- |
| 200        | OK                    | Successfully retrieve user details.                |
| 400        | Bad Request           | Service layer error with customised error message. |
| 500        | Internal Server Error | Unexpected exception occurred on the server-side.  |

---

### [GET] `http://localhost:4000/transactions`

- Request Header
  - Key: 'Authorization'
  - Value: 'Bearer <JWT_Token>'

```js
// Response:
{
    "transactions": [
        {
            "sender_name": "Alice",
            "beneficiary_name": "Charlie",
            "amount_transferred": 10,
            "amount_transferred_currency": "SGD",
            "amount_received": 10,
            "amount_received_currency": "USD",
            "status": "CONFIRMED",
            "date_transferred": "2023-07-02T11:00:23.512197Z",
            "date_received": "2023-07-02T11:00:23.512197Z"
        },
        {
            "sender_name": "Alice",
            "beneficiary_name": "Charlie",
            "amount_transferred": 500,
            "amount_transferred_currency": "SGD",
            "amount_received": 370,
            "amount_received_currency": "USD",
            "status": "CONFIRMED",
            "date_transferred": "2023-07-02T10:59:36.418704Z",
            "date_received": "2023-07-02T10:59:36.418704Z"
        }
    ]
}
```

| Error Code | Status                | Description                                        |
| ---------- | --------------------- | -------------------------------------------------- |
| 200        | OK                    | Successfully retrieve user details.                |
| 400        | Bad Request           | Service layer error with customised error message. |
| 500        | Internal Server Error | Unexpected exception occurred on the server-side.  |

---

### [POST] `http://localhost:4000/transaction`

- Request Header
  - Key: 'Authorization'
  - Value: 'Bearer <JWT_Token>'

```js
// Request Body:
{
    "beneficiary_name": "Bob",
    "mobile_number": "+65 89230122",
    "amount_transferred": "1000",
    "amount_transferred_currency": "SGD",
}
```

| Error Code | Status                | Description                                        |
| ---------- | --------------------- | -------------------------------------------------- |
| 201        | Created               | Successfully created a transaction.                |
| 400        | Bad Request           | Service layer error with customised error message. |
| 500        | Internal Server Error | Unexpected exception occurred on the server-side.  |

---
