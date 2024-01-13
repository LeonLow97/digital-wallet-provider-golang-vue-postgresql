# REST API Documentation

## Endpoints Summary

| Method | Description                  | Endpoint                    |
| :----: | ---------------------------- | --------------------------- |
|  POST  | User registration            | `/api/v1/register`          |
|  POST  | User login                   | `/api/v1/login`             |
|  GET   | User profile details         | `/api/v1/user/profile`      |
|  PUT   | Update user profile          | `/api/v1/user/profile`      |
|  POST  | Add beneficiary              | `/api/v1/beneficiary`       |
|  PUT   | Edit beneficiary information | `/api/v1/recipient/{id}`    |
| DELETE | Delete beneficiary           | `/api/v1/recipient/{id}`    |
|  POST  | Initiate money transfer      | `/api/v1/transfer`          |
|  GET   | Transaction history          | `/api/v1/transactions`      |
|  GET   | Transaction details          | `/api/v1/transactions/{id}` |
|  POST  | Deposit funds to account     | `/api/v1/deposit`           |
|  POST  | Withdraw funds from account  | `/api/v1/withdraw`          |
|  GET   | Real-time exchange rates     | `/api/v1/exchange-rates`    |
|  POST  | Subscribe to notifications   | `/api/v1/notifications`     |
