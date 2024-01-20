# Golang Clean Architecture

- [Example 1](https://github.com/Ikhlashmulya/golang-clean-architecture/)
- [amitshekhariitbhu](https://github.com/amitshekhariitbhu/go-backend-clean-architecture)
- [bxcodec](https://github.com/bxcodec/go-clean-arch)

## Clean Architecture Folder Structure

```
server
├── Dockerfile
├── README.md
├── cmd
│   └── api
│       └── main.go
├── config
│   └── development.yaml
├── delivery
│   └── http
│       └── handler
│           └── auth_handler.go
├── domain
│   └── user.go
├── dto
│   └── auth.go
├── exception
│   ├── error_user.go
│   ├── errors.go
│   └── response
│       └── messages.go
├── go.mod
├── go.sum
├── infrastructure
│   ├── config.go
│   └── database.go
├── init-db.sql
├── repository
│   └── user_repository.go
├── tests
├── usecase
│   └── auth_usecase.go
└── utils
    └── json_response.go
```

## Folder Structure Explanation

- **cmd**: Contains the main application entry point (`main.go`) and any additional entry points.
- **config**: Holds configuration files for different environments (e.g., `development.yaml`).
- **delivery/http/handler**: Manages HTTP request handling, handling incoming requests for authentication in this case.
- **domain**: Defines the core domain entities and business logic, such as the `User` entity.
- **dto**: Contains Data Transfer Objects (DTOs) that define the data structures exchanged between layers.
- **exception**: Handles error management, including specific errors related to users and general error messages.
- **infrastructure**: Deals with infrastructure concerns, such as database and configuration setup.
- **repository**: Implements the repository pattern, providing a way to interact with data storage (e.g., database).
- **usecase**: Implements use cases or application-specific business logic, such as authentication use cases.
- **utils**: Houses utility functions or helper methods, like `json_response.go` for managing JSON responses.
- **init-db.sql**: Initialization script for setting up the database schema and initial data.
- **tests**: Reserved for test files, including unit tests and integration tests.
- **README.md**: Documentation file providing an overview of the project structure and usage.

## Endpoints

| Method | Endpoint                        | Description                                                                             | Completed |
| :----: | ------------------------------- | --------------------------------------------------------------------------------------- | :-------: |
|  POST  | `/api/v1/users/register`        | Register a new user with required information like email, password and personal details |    [X]    |
|  POST  | `/api/v1/users/login`           | Authenticate a user and generate an access token for authentication                     |    [X]    |
|  POST  | `/api/v1/users/logout`          | Invalidate the user's access token, logging them out                                    |    [X]    |
|  POST  | `/api/v1/wallet`                | Create a new wallet for a user                                                          |    [X]    |
|  GET   | `/api/v1/wallet/all`            | Retrieve all wallets for a user                                                         |    [X]    |
|  GET   | `/api/v1/wallet`                | Retrieve one wallet for a user                                                          |    [X]    |
|  POST  | `/api/v1/balance/deposit`       | Deposit money into the user's wallet                                                    |    [X]    |
|  POST  | `/api/v1/balance/withdraw`      | Withdraw money from the user's wallet                                                   |    [X]    |
|  POST  | `/api/v1/transactions/transfer` | Transfer money between two user wallets                                                 |    [ ]    |
|  GET   | `/api/v1/transactions/history`  | Retrieve the transaction history for a user                                             |    [ ]    |
|  GET   | `/api/v1/users/profile`         | Retrieve the user's profile information                                                 |    [ ]    |
|  PUT   | `/api/v1/users/profile`         | Update the user's profile information                                                   |    [ ]    |
|  PUT   | `/api/v1/users/password`        | Change the user's password                                                              |    [ ]    |
|  POST  | `/api/v1/token/refresh`         | Request a new access token using a refresh token                                        |    [ ]    |
|  POST  | `/api/v1/exchange`              | Allow users to exchange one currency for another                                        |    [ ]    |
|  POST  | `/api/v1/kyc`                   | Implement KYC processes for user verification                                           |    [ ]    |
| DELETE | `/api/v1/users/close`           | Close a user's account and deactivate the associated wallet                             |    [ ]    |
