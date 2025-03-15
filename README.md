# Bank Management System(BackEnd)

A simple **Go + PostgreSQL** banking application that supports **user authentication**, **accounts**, **deposits**, **withdrawals**, **transfers**, and **loan management**. All endpoints are secured with JWT-based authentication, and data is stored in a Dockerized PostgreSQL database.

## Features

1. **User Authentication**  
   - Sign up, login, and JWT-based protection for certain routes.

2. **Account Management**  
   - Create accounts for each user.  
   - Deposit and withdraw endpoints (with transaction logging).  
   - Transfer funds between accounts atomically.

3. **Transaction History**  
   - Log every deposit, withdrawal, and transfer in a `transactions` table.  
   - Retrieve transaction history for each account.

4. **Loan and Credit**  
   - Apply for loans, approve/reject them, repay partially or fully.  
   - Track outstanding balances and loan statuses.

5. **Dockerized Setup**  
   - `docker-compose.yml` to spin up both the Go application and the PostgreSQL database.

## Tech Stack

- **Golang** for the application logic  
- **Gin** (web framework) for routing  
- **GORM** (ORM) for database operations  
- **PostgreSQL** as the main database  
- **Docker + Docker Compose** for container orchestration  
- **JWT** for authentication

## Prerequisites

- **Go** (1.20+ recommended)  
- **Docker** (with Docker Compose)  
- **Git** (optional but recommended)

## Project Structure

```
bank_management/
├── handlers/
│   ├── auth.go          # Signup, Login
│   ├── account.go       # Account operations (Create, Deposit, Withdraw, Transfer)
│   ├── loan.go          # Loan operations (Apply, Approve/Reject, Repay)
├── models/
│   ├── user.go          # User model
│   ├── account.go       # Account model
│   ├── transaction.go   # Transaction model
│   └── loan.go          # Loan model
├── utils/
│   └── jwt.go           # JWT secret & token generation
├── main.go              # Entry point, routes & migrations
├── Dockerfile           # Docker instructions for Go
├── docker-compose.yml   # Docker Compose file for app + PostgreSQL
├── go.mod
├── go.sum
└── README.md            # Project documentation
```

## Setup & Installation

1. **Clone the repo** (or download):
   ```bash
   git clone https://github.com/<yourusername>/bank_management.git
   cd bank_management
   ```
2. **Build & Run with Docker Compose**:
   ```bash
   docker-compose up --build
   ```
   - This starts both the Go app and PostgreSQL database.

3. **Check the logs** to ensure the server is running on `:8080` and Postgres on `:5432`.

4. **(Optional)** If running locally without Docker:
   - Make sure PostgreSQL is installed & running.  
   - Update the DSN in `main.go` to match your local DB settings.  
   - Run:
     ```bash
     go mod tidy
     go run main.go
     ```
   - The server should run at `http://localhost:8080`.

## Usage

### Authentication

- **POST /signup**  
  Body (JSON):
  ```json
  {
    "username": "jondoe",
    "email": "jon@example.com",
    "password": "mypassword"
  }
  ```
- **POST /login**  
  Body (JSON):
  ```json
  {
    "email": "jon@example.com",
    "password": "mypassword"
  }
  ```
  Returns a **JWT token**.

### Accounts

- **POST /accounts** (Create Account)  
  ```json
  {
    "user_id": 1
  }
  ```
- **POST /accounts/:id/deposit**  
  ```json
  {
    "amount": 100.0
  }
  ```
- **POST /accounts/:id/withdraw**  
  ```json
  {
    "amount": 50.0
  }
  ```
- **POST /accounts/transfer**  
  ```json
  {
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 100.0
  }
  ```
- **GET /accounts/:id/transactions**  
  View transaction history for the given account.

### Loans

- **POST /loans/apply**  
  ```json
  {
    "user_id": 1,
    "principal": 1000,
    "interest_rate": 5.0,
    "term_months": 12
  }
  ```
- **PATCH /loans/:id/status**  
  Approve/Reject a loan. Body:
  ```json
  {
    "status": "approved"
  }
  ```
- **POST /loans/:id/repay**  
  ```json
  {
    "amount": 300
  }
  ```

(These examples assume you’re sending requests with the required `Authorization: <token>` header if endpoints are protected.)

## Testing

- **Postman / cURL**:  
  Check the endpoints with JSON bodies. 
- **Unit Tests** (planned):  
  Add `go test` coverage in `handlers/` or a dedicated `tests/` folder.  
- **Integration**:  
  Spin up a test DB and run end-to-end checks for each route.

## Environment Variables

- **JWT_SECRET**: Optionally set this to your JWT secret if you don’t want it hardcoded.  
- **DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT**: For Docker Compose or local setup.

## Roadmap / Future Features

- **Role-Based Access** (admin vs. customer)  
- **Scheduled Jobs** for interest, fees  
- **Notification System** (email/SMS alerts)  
- **KYC & Compliance** (user verification)  
- **CI/CD Pipeline** (GitHub Actions or GitLab CI)

## Contributing

1. **Fork** the project.  
2. **Create a Feature Branch**: `git checkout -b feature/your-feature`  
3. **Commit Your Changes**: `git commit -m 'Add new feature'`  
4. **Push to the Branch**: `git push origin feature/your-feature`  
5. **Create a Pull Request** in GitHub.


**Enjoy using the Bank Management System!** If you have any questions or improvements, open an issue or submit a pull request.
