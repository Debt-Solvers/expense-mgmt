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

## Categories

<h5>Get Default Categories</h5>
<p>
  <em>Endpoint</em>:GET /api/v1/categories/defaults </br>
  <em>Query Parameters</em>:None </br>
  <em>Response</em>: </br>
  <code>
    {
    "status": 200,
    "message": "Fetched default categories successfully",
    "data": {
        "categories": [
            {
              "category_id": "74e21b0b-8ad8-49a1-8a01-7e74e250e713",
              "name": "Food & Dining",
              "description": "Restaurants, groceries, and food delivery",
              "color_code": "#FFD700",
              "is_default": true,
              "created_at": "2024-11-12T19:06:54.155524-05:00",
              "updated_at": "2024-11-12T19:06:54.155524-05:00",
              "deleted_at": null
            }
        ]
      }
    }
  </code>
</p>

<h5>Get Single Category Details</h5>
<p>
  <em>Endpoint</em>:GET /api/v1/categories/{categoryId} </br>
  <em>Query Parameters</em>:None </br>
  <em>Response</em>: </br>
  <code>
    {
      "status": 200,
      "message": "Category details fetched successfully",
      "data": {
          "category_id": "153789c9-3d4f-4b65-b25f-be065ecf028b",
          "name": "Others",
          "description": "Miscellaneous expenses",
          "color_code": "#A9A9A9",
          "is_default": true,
          "created_at": "2024-11-12T19:06:54.209896-05:00",
          "updated_at": "2024-11-12T19:06:54.209896-05:00",
          "deleted_at": null
      }
    }
  </br>
    {
      "status": 404,
      "message": "Category not found"
    }
  </code>
</p>

<h5>Update Category</h5>
<p>
  <em>Endpoint</em>:PUT /api/v1/categories/{categoryId} </br>
  <em>Query Parameters</em>:None </br>
  <em>Response</em>: </br>
  <code>
    {
      "status": 200,
      "message": "Category details fetched successfully",
      "data": {
          "category_id": "153789c9-3d4f-4b65-b25f-be065ecf028b",
          "name": "Others",
          "description": "Miscellaneous expenses",
          "color_code": "#A9A9A9",
          "is_default": true,
          "created_at": "2024-11-12T19:06:54.209896-05:00",
          "updated_at": "2024-11-12T19:06:54.209896-05:00",
          "deleted_at": null
      }
    }
  </br>
    {
      "status": 404,
      "message": "Category not found"
    }
  </code>
</p>

## Expenses

### Creating an expense (POST /api/v1/expenses)

### Listing expenses (GET /api/v1/expenses)

### Getting a single expense (GET /api/v1/expenses/{expenseId})

### Updating an expense (PUT /api/v1/expenses/{expenseId})

### Deleting an expense (DELETE /api/v1/expenses/{expenseId})

### Expense Analysis (GET /api/v1/expenses/analysis)

<p>This endpoint is more focused on aggregating and analyzing the data. It’s about providing insights based on the expenses. For example:</p></br>

<p>Total spending in a specific period (e.g., month, week).
  Expenditure by category or across time periods.
  Patterns in spending to help the user understand their spending habits.
  The Expense Analysis could use aggregate data, providing a summary or high-level view, rather than individual records, which is what CRUD operations focus on.
</p>

## Budgets

### POST /api/v1/budgets

<p>Description: Creates a new budget for the user.</p>

<p>
  <code>
    {
    "category_id": "d951a6bc-b346-4131-b294-fe7b33edcd59", // UUID, required, category for which the budget is set
    "amount": 500.00, // float, required, budget amount
    "start_date": "2024-12-01", // string, required, the start date in YYYY-MM-DD format
    "end_date": "2024-12-31" // string, required, the end date in YYYY-MM-DD format
    }
  </code>
</p>

<p>
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
  </code>
  <code>
    {
    "status": "error",
    "message": "Invalid input: {error_message}"
    }
  </code>
</p>

<p> 
  Endpoint: GET /api/v1/budgets

Endpoint: GET /api/v1/budgets/{budgetId}

This endpoint retrieves a specific budget by its ID.

JSON Query Parameters:
None

</p>
<p>  
  Get All Budgets (List Budgets)
  Endpoint: GET /api/v1/budgets

This endpoint will allow you to list all budgets, with optional query parameters for filtering.

JSON Query Parameters:
period: current, upcoming, past
category_id: Filter budgets by category
start_date: Filter budgets starting from this date
end_date: Filter budgets ending before this date
status: active, exceeded, upcoming

</p>
<p>
  Update Budget
  Endpoint: PUT /api/v1/budgets/{budgetId}

This endpoint will update an existing budget's amount, category, or date range.
{
"amount": 600.00,
"start_date": "2024-12-05",
"end_date": "2024-12-31",
"category_id": "a2f3b6c7-d567-492f-a8f7-b7c3b9d7e1d4"
}

  </p>

  <p> 
    Delete Budget
    Endpoint: DELETE /api/v1/budgets/{budgetId}

    This endpoint will remove a specific budget.

    JSON Body:
    None

  </p>

  <p>
    Budget Analysis (GET /api/v1/budgets/analysis)
    Request:
    Method: GET
    URL: /api/v1/budgets/analysis
    Query Parameters:
    category_id: d951a6bc-b346-4131-b294-fe7b33edcd59 (optional)
    start_date: 2024-12-01 (optional)
    end_date: 2024-12-31 (optional)

  </p>
  
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
