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

## Environment Varibles

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=expense_service_db

JWT_SECRET=DebtSolverSecret
JWT_EXPIRATION_HOURS=24

# API Endpoints

## Categories

### Get Default Categories

- **Endpoint**: `/api/v1/categories/defaults`
- **Method**: `GET`
- **Description**: Returns default categories.
- **Query Parameters**: `None`
- **Response**:
  #### Success
  ```json
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
  ```

### Auth User All Categories

- **Endpoint**: `GET /api/v1/categories/`
- **Query Parameters**: `None`
- **Response**:
  #### Success Response:
  ```json
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
  ```

#### Error Response

Occurs when the user provides an invalid or expired token.

```json
{
	"status": 401,
	"message": "Token is invalid or expired"
}
```

### Create Custom Category

- **Endpoint**: `POST /api/v1/categories/`
- **Description**: This endpoint allows users to create a custom category. Each category must have a unique name.

- **Query Parameters**:
  ```json
  {
  	"name": "Dating and Video Games!", // Required. The unique name of the category.
  	"description": "Online Streaming Contents" // Optional. A brief description of the category.
  }
  ```
- **Response**:
  #### Success
  ```json
  {
  	"status": 201,
  	"message": "Category created successfully",
  	"data": "05b218b7-e8c3-4b8d-9682-7d555996f2f6" // The unique ID of the created category.
  }
  ```
  #### Error
  ```json
  {
  	"error": "Category name already exists"
  }
  ```

### Recommendations:

- **Unique Category Name**: Ensure the category name is unique. If the name already exists, return an error with the message `"Category name already exists"`.
- **Description**: Validate the description field (if provided), ensuring it’s not excessively long and is meaningful.
- **Color Code**: If the `color_code` is provided, ensure it is in the correct format (e.g., `#RRGGBB`).

### Get Single Category Details

- **Endpoint**: `GET /api/v1/categories/{categoryId}`
- **Description**: This endpoint retrieves the details of a specific category by its unique ID. It provides information such as the category name, description, color code, default status, and timestamps for creation and updates.

- **Path Parameters**:

  - `categoryId` (required): The unique identifier of the category you want to retrieve.

- **Query Parameters**: None

- **Response**:

  #### Success

  ```json
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
  ```

  #### Error

  ```json
  {
  	"status": 404,
  	"message": "Category not found"
  }
  ```

### Update Category

- **Endpoint**: `PUT /api/v1/categories/{categoryId}`
- **Description**: This endpoint allows users to update the details of an existing category using its unique ID. You can modify attributes such as the category name, description, color code, and default status.

- **Path Parameters**:

  - `categoryId` (required): The unique identifier of the category to be updated.

- **Request Body**:

  - `name` (required): The updated name of the category.
  - `description` (optional): A brief description of the category.
  - `color_code` (optional): A hex color code to represent the category.
  - `is_default` (optional): A boolean value indicating if the category is the default one.

  **Example Request Body**:

  ```json
  {
  	"name": "Adventures",
  	"description": "New Grant theft auto San-Andreas",
  	"color_code": "#A8C0BD",
  	"is_default": false
  }
  ```

- **Response**:

  #### Success

  ```json
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
  ```

  #### Error

  ```json
  {
  	"status": 404,
  	"message": "Category not found"
  }
  ```

### Recommendations:

- **Valid categoryId**: Ensure the `categoryId` is a valid UUID and exists in the system. Return a `404` error if the category is not found.
- **Data Validation**: When updating the category name, ensure that the new name is unique. If the name already exists, return an appropriate error (`400` or `409`).
- **Description**: Validate that the description, if provided, is a valid string and within acceptable length limits (e.g., 255 characters).
- **Color Code**: If a `color_code` is provided, ensure it follows a valid format (e.g., `#RRGGBB`).

### Delete Category

- **Endpoint**: `DELETE /api/v1/categories/{categoryId}`
- **Description**: This endpoint allows users to delete a category by its unique ID. Deleting a category will remove it from the system. It is recommended to implement soft deletion (e.g., setting a `deleted_at` timestamp) rather than a permanent removal.

- **Path Parameters**:

  - `categoryId` (required): The unique identifier of the category to be deleted.

- **Query Parameters**: None

- **Response**:

  #### Success

  ```json
  {
  	"status": 200,
  	"message": "Category deleted successfully"
  }
  ```

  #### Error

  ```json
  {
  	"status": 404,
  	"message": "Category not found"
  }
  ```

### Recommendations:

- **Valid categoryId**: Ensure the `categoryId` is a valid UUID and exists in the system. Return a `404` error if the category is not found.
- **Deletion Logic**: Consider handling any dependencies that might exist with other resources (e.g., if the category is associated with expenses).
- **Access Control**: Verify that the user has permission to delete the specified category.

## Expenses

### Create a Single Expense

- **Endpoint**: `POST /api/v1/expenses/`
- **Description**: This endpoint allows users to create a new expense entry. The expense is associated with a category, includes a specific amount, date, and description, and can optionally be linked to a receipt.

- **Request Body**:

  - `category_id` (required): The unique identifier of the category the expense belongs to.
  - `amount` (required): The amount of the expense.
  - `date` (required): The date when the expense occurred in ISO 8601 format (e.g., `2024-11-14T00:00:00Z`).
  - `description` (required): A brief description of the expense.
  - `receipt_id` (optional): The unique identifier for a receipt, if the expense is associated with an uploaded receipt.

  **Example Request Body**:

  ```json
  {
  	"category_id": "8c135496-ea27-446b-919e-b312394c5f36",
  	"amount": 123.45,
  	"date": "2024-11-14T00:00:00Z",
  	"description": "Toyota Camry Car Insurance",
  	"receipt_id": "your-receipt-uuid-here" // Optional
  }
  ```

- **Response**:

  #### Success

  ```json
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
  ```

  #### Error

  ```json
  {
  	"error": "Category not found" // Example error message if the category_id doesn't exist
  }
  ```

### Recommendations:

- **Category Validation**: Ensure that the `category_id` exists in the system. If it does not, return a `400` error indicating that the category is invalid.
- **Amount Validation**: Ensure that the `amount` is a positive number. If it is zero or negative, return a `400` error.
- **Date Validation**: Ensure that the `date` is in the correct format (`YYYY-MM-DD`), and it’s not a future date unless necessary.
- **Description**: Make sure that the `description` is meaningful and within acceptable length limits (e.g., no more than 255 characters).
- **Optional Fields**: If the `receipt_id` is provided, validate that it exists and is a valid reference to a receipt in the system.

### List All Expenses

- **Endpoint**: `GET /api/v1/expenses/`
- **Description**: This endpoint allows users to retrieve a list of all expenses associated with their account. It supports pagination and filtering by date range, category, amount, and sorting by date.

- **Query Parameters**:

| Parameter     | Type     | Description                                    | Default | Options/Format             |
| ------------- | -------- | ---------------------------------------------- | ------- | -------------------------- |
| `page`        | [int]    | The page number for pagination                 | `1`     | N/A                        |
| `limit`       | [int]    | The number of items per page                   | `10`    | N/A                        |
| `start_date`  | [string] | The start date for filtering expenses          | N/A     | Format: `YYYY-MM-DD`       |
| `end_date`    | [string] | The end date for filtering expenses            | N/A     | Format: `YYYY-MM-DD`       |
| `category_id` | [string] | The ID of the category to filter expenses by   | N/A     | N/A                        |
| `min_amount`  | [float]  | The minimum amount to filter the expenses by   | N/A     | N/A                        |
| `max_amount`  | [float]  | The maximum amount to filter the expenses by   | N/A     | N/A                        |
| `sort`        | [string] | The field by which to sort the results         | `date`  | `date`                     |
| `order`       | [string] | The order of sorting (ascending or descending) | `asc`   | Options: `"asc"`, `"desc"` |

### Notes:

- **Pagination**: If no `page` or `limit` is specified, defaults are set to `page=1` and `limit=10`.
- **Date Filters**: The `start_date` and `end_date` parameters must follow the format `YYYY-MM-DD`.
- **Sorting**: The `sort` field defaults to `date`. Sorting order (`asc` or `desc`) can be specified using the `order` parameter.

### Example Request:

```http
GET /api/v1/expenses/?page=1&limit=10&start_date=2024-11-01&end_date=2024-11-30&category_id=8c135496-ea27-446b-919e-b312394c5f36&min_amount=50&max_amount=500&sort=date&order=desc
```

- **Response**:

  #### Success

  ```json
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
  			}
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
  ```

  #### Error

  ```json
  {
  	"error": "Category not found" // Example error message if the category_id does not exist
  }
  ```

### Recommendations:

- **Pagination**: Always use the `page` and `limit` parameters to prevent retrieving a large set of data. Ensure that pagination is handled correctly in the response.
- **Date Range**: Ensure the `start_date` and `end_date` parameters follow the correct format (`YYYY-MM-DD`).
- **Sorting**: Validate the `sort` and `order` parameters to ensure they are set correctly (e.g., `sort` should default to `date`, and `order` should be either `asc` or `desc`).
- **Category Validation**: If the `category_id` is provided, ensure that it exists in the system. If not, return a `404` error indicating that the category was not found.
- **Amount Filters**: Validate the `min_amount` and `max_amount` filters to ensure they are numeric values.

### Get Single Expense

- **Endpoint**: `GET /api/v1/expenses/{expenseId}`
- **Description**: This endpoint allows users to retrieve the details of a single expense by its unique `expenseId`.

### Query Parameters

None.

### Path Parameters

| Parameter   | Type   | Description                              | Options/Format |
| ----------- | ------ | ---------------------------------------- | -------------- |
| `expenseId` | string | The unique ID of the expense to retrieve | UUID format    |

- **Response**:

  #### Success

  ```json
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
  ```

  #### Error

  ```json
  {
  	"status": 404,
  	"message": "Expense not found"
  }
  ```

### Recommendations:

- **Valid expenseId**: Ensure that the `expenseId` passed in the URL is a valid UUID. If an invalid or non-existing `expenseId` is provided, return a `404` error with the message `"Expense not found"`.
- **Error Handling**: The `404` error should be returned if the `expenseId` does not exist in the system. You may also consider handling other types of errors (e.g., unauthorized access).
- **Access Control**: Make sure that the user can only access their own expenses. You can validate the `user_id` and ensure it matches the authenticated user's ID.

### Delete Single Expense

- **Endpoint**:

`DELETE /api/v1/expenses/{expenseId}`

- **Query Parameters**: None

- **Response**:

  ##### Success

  ```json
  {
  	"status": 200,
  	"message": "Expense deleted successfully"
  }
  ```

  #### Error

  ```json
  {
  	"status": 404,
  	"message": "Expense not found"
  }
  ```

### Update Single Expense

- **Endpoint**:
  `PUT /api/v1/expenses/{expenseId}`

- **Query Parameters**:
  ```json
  {
  	"amount": 90000,
  	"category_id": "05b218b7-e8c3-4b8d-9682-7d555996f2f6",
  	"date": "2024-11-14T00:00:00Z",
  	"description": "updated expenses"
  }
  ```
- **Response**:

  #### Success

  ```json
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
  ```

  #### Error

  ```json
  {
  	"status": 404,
  	"message": "Expense not found"
  }
  ```

### Expenses Analysis Endpoint

This endpoint is designed to aggregate and analyze expense data. It provides insights and summaries of the user's spending habits over specific time periods or by categories. The analysis may include:

- Total spending within a specific period (e.g., month, week).
- Expenditure by category or over time.
- Spending patterns to help the user understand their financial habits.

This endpoint aggregates data and offers a high-level overview rather than providing individual expense records, which is the focus of CRUD operations.

- **Endpoint**:
  `GET /api/v1/expenses/Analysis`
- **Example Query**:

```http
GET /api/v1/expenses/analysis?start_date=2024-01-01&end_date=2024-12-31&period=month&page=1&per_page=1&category_id=46bbdd03-d7d0-4a29-8f1e-31f9cfcc666e
```

- **Query Parameters**:

| Parameter     | Type   | Description                                                       | Default | Options/Format                          |
| ------------- | ------ | ----------------------------------------------------------------- | ------- | --------------------------------------- |
| `start_date`  | string | The start date for the analysis period. Format: `YYYY-MM-DD`      | N/A     | Format: `YYYY-MM-DD`                    |
| `end_date`    | string | The end date for the analysis period. Format: `YYYY-MM-DD`        | N/A     | Format: `YYYY-MM-DD`                    |
| `period`      | string | The period for which to analyze expenses (e.g., `month`, `week`). | `month` | Options: `"month"`, `"week"`, `"daily"` |
| `category_id` | uuid   | The category id to be analyze (e.g., `uuid`).                     |

- **Response**:

  #### Success

  ```json
  {
  	"status": 200,
  	"message": "Expense analysis fetched successfully",
  	"data": {
  		"period": "month",
  		"total_spending": 9150.45,
  		"average_spending": 1525.075,
  		"highest_expense": 4300,
  		"category_breakdown": [
  			{
  				"category_id": "0ec4e2ba-4623-4380-b1d0-eb3d0b0c3e6f",
  				"percentage": 8.201236004786649,
  				"total": 750.45
  			},
  			{
  				"category_id": "46bbdd03-d7d0-4a29-8f1e-31f9cfcc666e",
  				"percentage": 30.59958799840445,
  				"total": 2800
  			}
  		],
  		"most_frequent_category": {
  			"category_id": "46bbdd03-d7d0-4a29-8f1e-31f9cfcc666e",
  			"count": 2
  		},
  		"daily_average": 9150.45,
  		"pagination": {
  			"total_count": 2,
  			"page": 1,
  			"per_page": 10,
  			"total_pages": 1
  		}
  	}
  }
  ```

## Budgets

### Create Budget

This endpoint allows users to create a new budget. A budget is defined for a specific category and time period, with a set amount to track expenses against.

- **Endpoint**:
  `POST /api/v1/budgets`

- **Request Body**:

| Parameter     | Type   | Description                                           | Format                           | Required |
| ------------- | ------ | ----------------------------------------------------- | -------------------------------- | -------- |
| `category_id` | string | The UUID of the category for which the budget is set. | Format: `UUID`                   | Yes      |
| `amount`      | float  | The amount for the budget.                            | Format: Decimal (e.g., `500.00`) | Yes      |
| `start_date`  | string | The start date for the budget period.                 | Format: `YYYY-MM-DD`             | Yes      |
| `end_date`    | string | The end date for the budget period.                   | Format: `YYYY-MM-DD`             | Yes      |

#### Example Request Body

```json
{
	"category_id": "d951a6bc-b346-4131-b294-fe7b33edcd59",
	"amount": 500.0,
	"start_date": "2024-12-01",
	"end_date": "2024-12-31"
}
```

- **Response**:

  #### Success

  ```json
  {
  	"status": "success",
  	"message": "Budget created successfully",
  	"data": {
  		"budget_id": "uuid",
  		"user_id": "uuid",
  		"category_id": "uuid",
  		"amount": 500.0,
  		"start_date": "2024-12-01",
  		"end_date": "2024-12-31"
  	}
  }
  ```

  #### Error

  ```json
  {
  	"message": "Invalid input: Key: 'BudgetInput.Amount' Error:Field validation for 'Amount' failed on the 'gt' tag",
  	"data": null,
  	"errors": null
  }
  {
    "message": "Budget period overlaps with an existing budget for the same category",
    "data": null,
    "errors": null
  }

  ```

### Get Single Budget

This endpoint retrieves detailed information about a specific budget by its ID.

- **Endpoint**:
  `GET /api/v1/budgets/{budgetId}`

- **Query Parameters**:
- **None** (The budget ID is passed as part of the URL.)

- **Response**:

  #### Success

  ```json
  {
  	"status": "success",
  	"message": "Budget fetched successfully",
  	"data": {
  		"budget_id": "uuid", // The ID of the budget
  		"user_id": "uuid", // User's ID
  		"category_id": "uuid", // Category associated with the budget
  		"amount": 500.0, // The budgeted amount
  		"start_date": "2024-12-01", // Budget start date
  		"end_date": "2024-12-31" // Budget end date
  	}
  }
  ```

  #### Error

  ```json
  {
  	"status": "error",
  	"message": "Budget not found"
  }
  ```

### List All Budgets

This endpoint allows users to list all budgets with optional query parameters for filtering based on various criteria.

- **Endpoint**: `GET /api/v1/budgets`

- **Query Parameters**:

| Parameter     | Type   | Description                                                | Default | Options/Format                   |
| ------------- | ------ | ---------------------------------------------------------- | ------- | -------------------------------- |
| `category_id` | string | Filter budgets by category ID                              | None    | UUID                             |
| `start_date`  | string | Filter budgets starting from this date                     | None    | Format: `YYYY-MM-DD`             |
| `end_date`    | string | Filter budgets ending before this date                     | None    | Format: `YYYY-MM-DD`             |
| `period`      | string | Filter budgets by period: `current`, `upcoming`, `past`    | None    | `current`, `upcoming`, `past`    |
| `status`      | string | Filter budgets by status: `active`, `exceeded`, `upcoming` | None    | `active`, `exceeded`, `upcoming` |

- **Response**:

  #### Success

  ```json
  {
  	"status": "success",
  	"message": "Budgets fetched successfully",
  	"data": [
  		{
  			"budget_id": "uuid", // The ID of the budget
  			"user_id": "uuid", // User's ID
  			"category_id": "uuid", // Category associated with the budget
  			"amount": 500.0, // The budgeted amount
  			"start_date": "2024-12-01", // Budget start date
  			"end_date": "2024-12-31" // Budget end date
  		},
  		{
  			"budget_id": "uuid", // The ID of the budget
  			"user_id": "uuid", // User's ID
  			"category_id": "uuid", // Category associated with the budget
  			"amount": 300.0, // The budgeted amount
  			"start_date": "2024-01-01", // Budget start date
  			"end_date": "2024-01-31" // Budget end date
  		}
  	]
  }
  ```

  #### Error

  ```json
  {
  	"status": "error",
  	"message": "Failed to fetch budgets"
  }
  ```

### Update Single Budget

This endpoint allows users to update an existing budget. The user can modify the amount, category, or date range of an existing budget.

- **Endpoint**: `PUT /api/v1/budgets/{budgetId}`

- **Request Body**:

The body of the request should contain the fields that need to be updated. The following parameters are required:

| Parameter     | Type   | Description                                   | Format/Options                                      |
| ------------- | ------ | --------------------------------------------- | --------------------------------------------------- |
| `amount`      | float  | The new amount for the budget.                | Positive float (e.g., 600.00)                       |
| `category_id` | string | The category ID to associate with the budget. | UUID (e.g., `d951a6bc-b346-4131-b294-fe7b33edcd59`) |
| `start_date`  | string | The start date for the updated budget.        | Date format `YYYY-MM-DD`                            |
| `end_date`    | string | The end date for the updated budget.          | Date format `YYYY-MM-DD`                            |

- **Example Request Body**:

  ```json
  {
  	"amount": 600.0,
  	"category_id": "a2f3b6c7-d567-492f-a8f7-b7c3b9d7e1d4",
  	"start_date": "2024-12-05",
  	"end_date": "2024-12-31"
  }
  ```

- **Response**:

  #### Success

  ```json
  {
  	"status": "success",
  	"message": "Budget updated successfully",
  	"data": {
  		"budget_id": "uuid",
  		"user_id": "uuid",
  		"category_id": "a2f3b6c7-d567-492f-a8f7-b7c3b9d7e1d4",
  		"amount": 600.0,
  		"start_date": "2024-12-05",
  		"end_date": "2024-12-31"
  	}
  }
  ```

  #### Error

  ```json
  {
  	"status": "error",
  	"message": "Invalid input: {error_message}"
  }

  {
    "status": "error",
    "message": "Budget not found"
  }

  ```

### Delete Budget

This endpoint allows users to delete a specific budget.

- **Endpoint**: `DELETE /api/v1/budgets/{budgetId}`

- **Request Body**: `None`

- **Example Request**:

- **No JSON body is needed** for this endpoint.

- **Response**:

  #### Success

  ```json
  {
  	"status": "success",
  	"message": "Budget deleted successfully"
  }
  ```

  #### Error

  ```json
  {
  	"status": "error",
  	"message": "Budget not found"
  }
  ```

### Budget Analysis

This endpoint allows users to fetch an analysis of budgets, including details on spending and budget status for different categories.

- **Endpoint**: `GET /api/v1/budgets/analysis`

- **Query Parameters** (Optional):

  - **`category_id`**: (string, optional) The ID of the category to filter the analysis. If not provided, all categories are included.
  - **`start_date`**: (string, optional) The start date for the analysis period in the format `YYYY-MM-DD`.
  - **`end_date`**: (string, optional) The end date for the analysis period in the format `YYYY-MM-DD`.

- **Example Request**:

  ```http
  GET /api/v1/budgets/analysis?category_id=09880493-bf02-4d5a-87df-e515d0c39dc1&start_date=2024-11-01&end_date=2024-11-30
  ```

- **Response**

  #### Success

  ```json
  {
  	"status": 200,
  	"message": "Budget analysis fetched successfully",
  	"data": [
  		{
  			"category_id": "09880493-bf02-4d5a-87df-e515d0c39dc1",
  			"category": "Groceries",
  			"budgeted_amount": 500.0,
  			"total_spent": 650.25,
  			"remaining_budget": -150.25,
  			"percentage_spent": 130.05,
  			"exceeds_budget": true
  		},
  		{
  			"category_id": "e3d5f0ba-4623-11ec-81d3-0242ac130003",
  			"category": "Transport",
  			"budgeted_amount": 200.0,
  			"total_spent": 180.0,
  			"remaining_budget": 20.0,
  			"percentage_spent": 90.0,
  			"exceeds_budget": false
  		}
  	],
  	"errors": null
  }
  ```

  #### Error

  ```json
  {
  	"status": 400,
  	"message": "Invalid date range provided",
  	"errors": {
  		"start_date": "Must be a valid date in YYYY-MM-DD format",
  		"end_date": "Must be a valid date in YYYY-MM-DD format"
  	}
  }
  ```

### Receipts

## License

&copy This project is open-source and licensed under the MIT License.

## Contributions

Contributions are welcome! Feel free to open an issue or submit a pull request.
