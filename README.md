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

### Get Default Categories

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

<h5>Auth User All Categories</h5>
<p>
  <em>Endpoint</em>:GET /api/v1/categories/</br>
  <em>Query Parameters</em>:None </br>
  <em>Response</em>: </br>
  <code>
    {
    "status": 200,
    "message": "Categories retrieved successfully",
    "data": [
      {
        "category_id": "35983bbe-a8f1-4b76-bf42-ec598a4791e2",
        "name": "Education",
        "description": "Tuition, books, courses, training",
        "color_code": "#4682B4",
        "is_default": true,
        "created_at": "2024-11-14T12:37:36.101612-05:00",
        "updated_at": "2024-11-14T12:37:36.101612-05:00",
        "deleted_at": null
      },
        {
          "category_id": "c48e168b-fab9-45a3-a72a-9648e4aca537",
          "name": "Utilities",
          "description": "Electricity, water, internet, phone",
          "color_code": "#A9A9A9",
          "is_default": true,
          "created_at": "2024-11-14T12:37:36.104216-05:00",
          "updated_at": "2024-11-14T12:37:36.104216-05:00",
          "deleted_at": null
        }
      ]
    }
  </br>
  {
    "status": 401,
    "message": "Token is invalid or expired"
  }
  </code>
</p>

<h5>Create Custom Category</h5>
<p>
  <em>Endpoint</em>:POST /api/v1/categories/</br>
  <em>Query Parameters</em> </br>
    <code>
    {
      "name": "Dating and Video Games!",
      "description":"Online Streaming Contents" | optional
    }
    </code> 
    </br>
  <em>Response</em>: </br>
  <code>
    {
      "status": 201,
      "message": "Category created successfully",
      "data": "05b218b7-e8c3-4b8d-9682-7d555996f2f6"
    }
    <hr>

    {
      "error": "Category name already exists"
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
  <em>Query Parameters</em>:

  </br>
  <em>Response</em>: </br>
  <code>
    {
    "status": 200,
    "message": "Category updated successfully",
    "data": {
      "category_id": "05b218b7-e8c3-4b8d-9682-7d555996f2f6",
      "user_id": "f3486758-899e-462c-98b7-ba8f691c8718",
      "name": "Adventures",
      "description": "New Grant theft auto San-Andreas",
      "color_code": "#A8C0BD",
      "is_default": false,
      "created_at": "2024-11-14T17:07:32.443566-05:00",
      "updated_at": "2024-11-14T22:06:37.9337733-05:00",
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

<h5>Delete Category</h5>
<p>
  <em>Endpoint</em>:DELETE /api/v1/categories/{categoryId} </br>
  <em>Query Parameters</em>:None </br>
  <em>Response</em>: </br>
  <code>
    {
      "status": 200,
      "message": "Category deleted successfully"
    }
  </br>
    {
      "status": 404,
      "message": "Category not found"
    }
  </code>
</p>

## Expenses

<h5>Create a Single Expense</h5>
<p>
  <em>Endpoint</em>:POST /api/v1/expenses/</br>
  <em>Query Parameters</em> </br>
    <code>
      {
        "category_id": "8c135496-ea27-446b-919e-b312394c5f36",
        "amount": 123.45,
        "date": "2024-11-14T00:00:00Z",
        "description": "Toyota Camry Car Insurance"
        "receipt_id": "your-receipt-uuid-here" // optional but can be linked to an already uploaded receipt. 
      }
    </code>
<br>
<em>Response</em>: </br>
<code>
    {
      "status": 200,
      "message": "Expense created successfully",
      "data": {
        "expense_id": "b0b87e74-b3aa-481d-a91e-d240cac56e0a",
        "user_id": "f3486758-899e-462c-98b7-ba8f691c8718",
        "category_id": "8c135496-ea27-446b-919e-b312394c5f36",
        "amount": 123.45,
        "date": "2024-11-14T00:00:00Z",
        "description": "Toyota Camry Car Insurance",
        "receipt_id": null,
        "created_at": "2024-11-14T18:00:22.735473-05:00",
        "updated_at": "2024-11-14T18:00:22.735473-05:00"
      }
    }

<hr>

    {
      "error": "Category name already exists"
    }

  </code>
</p>

## Auth User All Expenses

### Endpoint

`GET /api/v1/expenses/`

### Query Parameters

| Parameter     | Type     | Description                                  | Default | Options/Format             |
| ------------- | -------- | -------------------------------------------- | ------- | -------------------------- |
| `page`        | [int]    | The page number for pagination               | `1`     | N/A                        |
| `limit`       | [int]    | The number of items per page                 | `10`    | N/A                        |
| `start_date`  | [string] | The start date for filtering expenses        | N/A     | Format: `YYYY-MM-DD`       |
| `end_date`    | [string] | The end date for filtering expenses          | N/A     | Format: `YYYY-MM-DD`       |
| `category_id` | [string] | The ID of the category to filter expenses by | N/A     | N/A                        |
| `min_amount`  | [float]  | The minimum amount to filter the expenses by | N/A     | N/A                        |
| `max_amount`  | [float]  | The maximum amount to filter the expenses by | N/A     | N/A                        |
| `sort`        | [string] | The field by which to sort the results       | `date`  | N/A                        |
| `order`       | [string] | The order of sorting                         | `asc`   | Options: `"asc"`, `"desc"` |

### Notes:

- By default, the `page` is set to `1`, and `limit` is set to `10` if not specified.
- Date filters must follow the format `YYYY-MM-DD`.
- Sorting can be done by `date` (default), and the sorting order can be either `asc` or `desc`.

  </br>
  <em>Response</em>: </br>
  <code>
  {
    "status": 200,
    "message": "Expenses fetched successfully",
    "data": {
        "expenses": [
            {
                "expense_id": "b0b87e74-b3aa-481d-a91e-d240cac56e0a",
                "user_id": "f3486758-899e-462c-98b7-ba8f691c8718",
                "category_id": "8c135496-ea27-446b-919e-b312394c5f36",
                "amount": 123.45,
                "date": "2024-11-13T19:00:00-05:00",
                "description": "Toyota Camry Car Insurance",
                "receipt_id": null,
                "created_at": "2024-11-14T18:00:22.735473-05:00",
                "updated_at": "2024-11-14T18:00:22.735473-05:00"
            },
            {
                "expense_id": "52924f2a-f67c-4fa9-889a-60afa1518f3b",
                "user_id": "f3486758-899e-462c-98b7-ba8f691c8718",
                "category_id": "8c135496-ea27-446b-919e-b312394c5f36",
                "amount": 123.45,
                "date": "2024-11-13T19:00:00-05:00",
                "description": "Lamborghini",
                "receipt_id": null,
                "created_at": "2024-11-14T19:38:21.557713-05:00",
                "updated_at": "2024-11-14T19:38:21.557713-05:00"
            },
        ],
        "pagination": {
            "total_count": 3,
            "page": 1,
            "per_page": 10,
            "total_pages": 1
        }
    },
    "errors": null
  }

<hr>

    {
      "error": "Category name already exists"
    }

  </code>
</p>

## Listing expenses (GET /api/v1/expenses)

<h5>Get Single Expense</h5>
<p>
  <em>Endpoint</em>:GET /api/v1/expenses/{expenseId} </br>
  <em>Query Parameters</em>:None </br>
  <em>Response</em>: </br>
  <code>
  {
    "status": 200,
    "message": "Expense fetched successfully",
    "data": {
        "expense_id": "b0b87e74-b3aa-481d-a91e-d240cac56e0a",
        "user_id": "f3486758-899e-462c-98b7-ba8f691c8718",
        "category_id": "8c135496-ea27-446b-919e-b312394c5f36",
        "amount": 123.45,
        "date": "2024-11-13T19:00:00-05:00",
        "description": "Toyota Camry Car Insurance",
        "receipt_id": null,
        "created_at": "2024-11-14T18:00:22.735473-05:00",
        "updated_at": "2024-11-14T18:00:22.735473-05:00"
    }
  }
  </br>
    {
      "status": 404,
      "message": "Expense not found"
    }
  </code>
</p>

<h5>Delete Single Expense</h5>
<p>
  <em>Endpoint</em>:DELETE /api/v1/expenses/{expenseId} </br>
  <em>Query Parameters</em>:None </br>
  <em>Response</em>: </br>
  <code>
  {
    "status": 200,
    "message": "Expense deleted successfully"
  }
  </br>
    {
      "status": 404,
      "message": "Expense not found"
    }
  </code>
</p>

<h5>Update Single Expense</h5>
<p>
  <em>Endpoint</em>:PUT /api/v1/expenses/{expenseId} </br>
  <em>Query Parameters</em>:
    <code>
    {
      "amount": 90000,
      "category_id": "05b218b7-e8c3-4b8d-9682-7d555996f2f6",
      "date": "2024-11-14T00:00:00Z",
      "description":"updated expenses"
    }

    </code>

  </br>
  <em>Response</em>: </br>
  <code>
    {
      "status": 200,
      "message": "Expense deleted successfully"
    }
  </br>
   {
    "status": 200,
    "message": "Expense updated successfully",
    "data": {
        "expense_id": "f00c1593-8aca-43f1-bbd3-42d3c96c725d",
        "user_id": "f3486758-899e-462c-98b7-ba8f691c8718",
        "category_id": "05b218b7-e8c3-4b8d-9682-7d555996f2f6",
        "amount": 90000,
        "date": "2024-11-13T19:00:00-05:00",
        "description": "updated expenses",
        "receipt_id": null,
        "created_at": "2024-11-14T22:28:18.266748-05:00",
        "updated_at": "2024-11-14T22:29:49.001179-05:00"
    }
  }
  </code>
</p>

<h5>Expenses Analysis Endpoint</h5>

<p>This endpoint is more focused on aggregating and analyzing the data. It’s about providing insights based on the expenses. For example:</p></br>

<p>Total spending in a specific period (e.g., month, week).
  Expenditure by category or across time periods.
  Patterns in spending to help the user understand their spending habits.
  The Expense Analysis could use aggregate data, providing a summary or high-level view, rather than individual records, which is what CRUD operations focus on.
</p>

<em>Endpoint</em>:GET /api/v1/expenses/Analysis?start_date=2024-01-01&end_date=2024-12-31&period=month </br>
<em>Query Parameters</em>: None
</br>

<em>Response</em>: </br>
<code>
{
"status": 200,
"message": "Expense analysis fetched successfully",
"data": {
"period": "daily",
"total_spending": 91200,
"average_spending": 249.18032786885246,
"highest_expense": 90000
}
}
<br>

</code>
</br>

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
<h5>Auth User All Budgets</h5>
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

<h5>Update Singele Budget</h5>
<p> Endpoint: PUT /api/v1/budgets/{budgetId}

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
