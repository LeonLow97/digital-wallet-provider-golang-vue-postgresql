# Digital Wallet

## Table of Contents

## Overview

Digital Wallet is a platform that enables **customer-to-customer (C2C)** transactions in multiple currencies across the globe. This application simplifies and accelerates money transfers by reducing dependence on traditional financial institutions, which often involve lengthy processing times due to regulatory checks and high transaction fees. Our solution ensures secure, fast, and cost-effective international transactions.

# Technologies

| Technology | Type                   | Version | Ports |
| ---------- | ---------------------- | :-----: | ----- |
| Golang     | Backend                |         | 8080  |
| PostgreSQL | Database               |         | 5432  |
| Redis      | Cache                  |         | 6379  |
| Vue        | Frontend               |         | 3000  |
| Docker     | Containerization       |         |       |
| Kubernetes | Orchestration          |         |       |
| REST       | Communication Protocol |         |       |

# Monolithic Architecture

Digital Wallet follows a Monolithic Architecture due to our small user base system. This architecture simplifies deployment and development but may require refactoring into a microservices-based approach as scalability needs grow.

# Assumptions

- **Regulatory Compliance & KYC**: Customers have completed Know Your Customer (KYC) and compliance checks before using the platform.
- **Stable Exchange Rates**: All exchange rates during transfers are pegged and do not fluctuate with the Forex Market, ensuring predictable transaction amounts.

# Endpoints

| Method  | Endpoint                          | Microservice           | Description                                             |
| ------- | --------------------------------- | ---------------------- | ------------------------------------------------------- |
| `GET`   | `/health`                         | API Gateway            | Health check endpoint to verify the gateway is running. |
| `POST`  | `/login`                          | Authentication Service | Endpoint to authenticate users (login).                 |
| `POST`  | `/signup`                         | Authentication Service | Endpoint for user registration (signup).                |
| `POST`  | `/logout`                         | Authentication Service | Endpoint to log out users.                              |
| `PATCH` | `/change-password`                | Authentication Service | Endpoint to change user password.                       |
| `POST`  | `/configure-mfa`                  | Authentication Service | Endpoint to set up multi-factor authentication (MFA).   |
| `POST`  | `/verify-mfa`                     | Authentication Service | Endpoint to verify MFA authentication.                  |
| `POST`  | `/password-reset/send`            | Authentication Service | Sends a password reset email.                           |
| `PATCH` | `/password-reset/reset`           | Authentication Service | Resets the user password.                               |
| `PUT`   | `/users/profile`                  | Authentication Service | Endpoint to update user profile information.            |
| `GET`   | `/users/me`                       | Authentication Service | Endpoint to retrieve current user details.              |
| `GET`   | `/balances`                       | Balance Service        | Retrieves user balances.                                |
| `GET`   | `/balances/{id}`                  | Balance Service        | Retrieves a specific balance by ID.                     |
| `GET`   | `/balances/history/{id}`          | Balance Service        | Retrieves the balance history for a given ID.           |
| `GET`   | `/balances/currencies`            | Balance Service        | Retrieves available balance currencies.                 |
| `POST`  | `/balances/deposit`               | Balance Service        | Deposits funds into a balance.                          |
| `POST`  | `/balances/withdraw`              | Balance Service        | Withdraws funds from a balance.                         |
| `PATCH` | `/balances/currency-exchange`     | Balance Service        | Performs a currency exchange.                           |
| `POST`  | `/balances/preview-exchange`      | Balance Service        | Previews a currency exchange before executing it.       |
| `POST`  | `/beneficiary`                    | Beneficiary Service    | Creates a new beneficiary.                              |
| `PUT`   | `/beneficiary`                    | Beneficiary Service    | Updates an existing beneficiary.                        |
| `GET`   | `/beneficiary/{id}`               | Beneficiary Service    | Retrieves a specific beneficiary by ID.                 |
| `GET`   | `/beneficiary`                    | Beneficiary Service    | Retrieves all beneficiaries.                            |
| `GET`   | `/wallet/{id}`                    | Wallet Service         | Retrieves a specific wallet by ID.                      |
| `GET`   | `/wallet/all`                     | Wallet Service         | Retrieves all wallets for the user.                     |
| `GET`   | `/wallet/types`                   | Wallet Service         | Retrieves available wallet types.                       |
| `POST`  | `/wallet`                         | Wallet Service         | Creates a new wallet.                                   |
| `PUT`   | `/wallet/update/{id}/{operation}` | Wallet Service         | Updates a wallet based on the operation.                |
| `POST`  | `/transaction`                    | Transaction Service    | Creates a new transaction.                              |
| `GET`   | `/transaction/all`                | Transaction Service    | Retrieves all transactions.                             |
