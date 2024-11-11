# Debt Solver - Expense Management Microservice

This repository contains the **Expense Management** microservice for the **Debt Solver** project, a mobile application that enables users to track and manage their finances. This service handles key functionalities such as budgeting, categorizing expenses, managing receipts, and generating expense reports.

## Key Features

- **Expense Tracking**: Record and categorize user expenses for easier tracking.
- **Budget Allocation**: Allow users to set and monitor spending limits across various categories.
- **Receipts Management**: Upload and process receipts using OCR for automatic expense entry.
- **Expense Reports**: Generate insights and summaries of spending habits.

## Technologies Used

- **Golang & Gin**: For building the service.
- **PostgreSQL**: For data storage (expenses, budgets, categories, and receipts).
- **GORM**: For ORM database interactions.
- **JWT**: For user authorization on protected endpoints.
- **Viper**: For configuration management.

## Directory Structure

```plaintext
expense-service/
│
├── cmd/
│   └── expense-service/
│       └── main.go                  # Entry point for the application
│
├── configs/
│   └── config.yaml                  # Configuration file for the service
│
├── db/
│   └── migrate.go                   # Database migrations for creating tables
│
├── internal/
│   ├── common/
│   │   └── common.go                # Common utility functions
│   ├── controller/
│   │   ├── budget_controller.go     # Controller for budget management
│   │   ├── category_controller.go   # Controller for expense categories
│   │   ├── expense_controller.go    # Controller for expense entries
│   │   └── receipt_controller.go    # Controller for receipt handling and OCR
│   ├── middleware/
│   │   └── auth_middleware.go       # Middleware for JWT-based route protection
│   ├── model/
│   │   ├── budget.go                # Budget model and database interactions
│   │   ├── category.go              # Category model and database interactions
│   │   ├── expense.go               # Expense model and database interactions
│   │   └── receipt.go               # Receipt model and database interactions
│   └── routes/
│       └── routes.go                # Define routes for all expense-related endpoints
│
├── utils/
│   ├── response.go                  # Utility functions for handling responses
│   └── ocr_utils.go                 # Utility functions for OCR processing
│
├── Dockerfile                       # Dockerfile for building the container
├── go.mod                           # Go module file
└── README.md                        # Project documentation
```

## Setup and Installation

git clone https://github.com/debt-solver/expense-service.git
cd expense-service

## Setup and PostgreSQL

docker run --name debt-solver-expense-db -e POSTGRES_PASSWORD=yourpassword -d -p 5432:5432 postgres

## Install Dependencies

go mod tidy

## Run Database Migrateions

go run db/migrate.go

## Run the Application

go run cmd/expense-service/main.go

## Build and Run with Docker

docker build -t expense-service .
docker run -p 8081:8081 expense-service

## API Endpoints

### Categories

### Expenses

### Budgets

POST /api/v1/budgets
Description: Creates a new budget for the user.
<code>{
"category_id": "d951a6bc-b346-4131-b294-fe7b33edcd59", // UUID, required, category for which the budget is set
"amount": 500.00, // float, required, budget amount
"start_date": "2024-12-01", // string, required, the start date in YYYY-MM-DD format
"end_date": "2024-12-31" // string, required, the end date in YYYY-MM-DD format
}
</code>
<code>
{
"status": "success",
"message": "Budget created successfully",
"data": {
"budget_id": "uuid", // The generated budget ID
"user_id": "uuid", // User's ID
"category_id": "uuid", // Associated category ID
"amount": 500.00, // Budget amount
"start_date": "2024-12-01", // Budget start date
"end_date": "2024-12-31" // Budget end date
}
}

{
"status": "error",
"message": "Invalid input: {error_message}"
}

</code>

### Receipts

## Environment Varibles

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=expense_service_db

JWT_SECRET=DebtSolverSecret
JWT_EXPIRATION_HOURS=24

## License

&copy This project is open-source and licensed under the MIT License.

## Contributions

Contributions are welcome! Feel free to open an issue or submit a pull request.
