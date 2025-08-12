# Service-Module

A RESTful Go micro-service for managing services and their versions. Supports listing services with filtering, sorting, and pagination, as well as fetching individual services with their versions.

## Table of Contents

- Tech Stack

- Setup & Run

- Features

- API Endpoints

- Design & Architecture

- Error Handling

- Extensibility & Future Scope

- Testing


##  Tech Stack

- Language: Go (1.20+ recommended)

- Web Framework: Gin

- ORM: GORM

- Database: Any supported by GORM (e.g., MySQL, PostgreSQL, SQLite)

- Configuration: Environment variables managed via .env file and godotenv


## Setup & Run

### Prerequisites

- Go installed (>= 1.20 recommended)

- Database server running (e.g., MySQL)

### Steps

1. Clone the repo:
    - git clone https://github.com/Shivaygupta/Go-Service-Module

2. Set up db with the details provided in .env file or create custom db details and update db details in the .env file.

3. Build and run the application:
    - go run main.go

4. The Service will be accessible at http://localhost:8080


## Features

- List services with:

    - Filtering by name/description (case-insensitive)

    - Sorting by allowed fields (id, name, created_at, updated_at)

    - Pagination with page & limit parameters

- Retrieve a service by ID including all its versions

- Create service with versions and description(Optional)

- Delete service by id including all its versions(Optional)

- Clean separation of concerns: handlers, services, repositories

- Structured error handling with consistent API error responses

- Middleware for logging and recovery

- Configuration via environment variables


## API Endpoints

#### 1. Get List of Services

- curl 
    - Request Type: GET
    - http://localhost:8080/api/v1/services?page=1&limit=5&name=payment&sortBy=name&order=asc

- Response

     ``` {
        "data": [
            {
                "id": ,
                "name": "",
                "description": "",
                "created_at": "",
                "updated_at": "",
                "versions": [
                    {
                        "id": ,
                        "service_id": ,
                        "version": "",
                        "created_at": ""
                    }
                ]
            }
         ],
        "total_count": ,
        "page": ,
        "limit": ,
        "total_pages": 
    }     

-  Error Response

```
    {
        "error": "",
        "message": "",
        "status_code": 
    }
```

#### 2. Get Service by ID

- curl
    - Request Type: GET
    - http://localhost:8080/api/v1/services

- Response

``` {
        "data": [
            {
                "id": ,
                "name": "",
                "description": "",
                "created_at": "",
                "updated_at": "",
                "versions": [
                    {
                        "id": ,
                        "service_id": ,
                        "version": "",
                        "created_at": ""
                    }
                ]
            }
         ],
        "total_count": ,
        "page": ,
        "limit": ,
        "total_pages": 
    }
```

Error Response:
```
    {
        "error": "",
        "message": "",
        "status_code": 
    }
```
#### 3. Create Service

- curl
    - Request Type: POST
    - Endpoint: /api/v1/services
    - Request Body:
    ```
    {
        "name": "Try Service",
        "description": "Handles all payment transactions",
        "versions": [
            {
                "version": "v1.0"
            },
            {
                "version": "v1.1"
            }
        ]
    }
    ```
    - Response:
    ```
    {
        "data": {
            "id": 3,
            "name": "Try Service",
            "description": "Handles all payment transactions",
            "created_at": "2025-08-11T21:35:21.54+05:30",
            "updated_at": "2025-08-11T21:35:21.54+05:30",
            "versions": [
                {
                    "id": 5,
                    "service_id": 3,
                    "version": "v1.0",
                    "created_at": "2025-08-11T21:35:21.541+05:30"
                },
                {
                    "id": 6,
                    "service_id": 3,
                    "version": "v1.1",
                    "created_at": "2025-08-11T21:35:21.541+05:30"
                }
            ]
        }
    }
    ```

    - Error Response
        ```
        {
            "error": "",
            "message": "",
            "status_code": 
        }
        ```

    #### 4. Delete Service By Id

    - curl
        - Request Type: POST
        - Endpoint: /api/v1/services
        - Path Param: id
    
    - Response:
        - "message": "Service deleted successfully"

    - Error Response
        ```
        {
            "error": "service not found",
            "message": "Service with give id doesnt exists",
            "status_code": 404
        }
        ```

## Design & Architecture

### Why Gin (web framework) and GORM (ORM)

#### Gin

- Lightweight, high-performance HTTP framework built for Go.

- Minimal overhead and fast request routing — good for production APIs.

- Middleware-first design makes adding logging, auth, recovery easy.

- Good defaults for JSON binding & validation via binding tags.

#### GORM

- Mature ORM for Go that simplifies DB interactions while still allowing raw SQL.

- Handles common concerns: migrations (AutoMigrate), associations (Preload), transactions, constraints.

- Speeds development (less boilerplate SQL) while preserving performance tuning options.

- Works with multiple DB backends — useful for CI (sqlite) vs production (postgres/mysql).

#### Tradeoff summary

- Faster development & clarity with ORM vs hand-written SQL performance edge. We keep the option to drop into raw SQL for hot paths.

- Gin gives speed + ergonomics; for extremely high-throughput needs you could swap to a lower-level handler later.

### High-Level Design (HLD)

#### Goals

- Provide REST API for Services and ServiceVersions.

- Support list (filter/sort/paginate), fetch-by-id (with versions), and create.

- Clean separation of concerns.

#### Components

- HTTP/API Layer (Gin) — routes, request validation, response formatting

- Middleware — logging, recovery, request ID, global error handler

- Service Layer — business logic, input validation beyond request-level, orchestration

- Repository Layer (GORM) — DB queries, mapping to models, pagination logic

- Database — Relational DB (Postgres/MySQL/SQLite for tests)

- Config — environment-driven configuration loader

- Observability — structured logs, metrics (future), traces (future)


#### Assumptions & Extensibility Considerations

- Framework Selection (Gin)
    - Chosen for its simplicity, performance, and active community support.
    - The routing structure is modular to allow easy integration of new API versions or additional middleware without disrupting existing endpoints.

- Error Handling Strategy
    - Implemented a centralized error-handling middleware to ensure consistent API responses.
    - Future scope includes supporting multiple error formats (custom structs) to accommodate different client needs, such as mobile vs. web consumers.

- Routing Structure
    - Grouped routes under /api/v1 to enable versioned API releases without breaking existing integrations.
    - Prepared for plug-and-play route modules to add new service domains quickly.

- Middleware Usage
    - Standard middlewares (Logger, Recovery, custom ErrorHandler) applied globally to enforce uniform behavior.
    - Designed so domain-specific middleware can be introduced per route group in the future.

- Service Layer Design
    - Business logic is separated from the controller layer for testability and reusability.
    - Structure allows easy addition of new service methods without modifying existing handlers.

- Persistence Layer
    - Repository pattern used to abstract DB operations, enabling future migration to a different database with minimal changes to upper layers.
    - Current implementation optimized for basic CRUD, but extensible to support caching layers (e.g., Redis) or search engines (e.g., Elasticsearch).

- Search Functionality Scalability
    - Using SQL LIKE queries for search, suitable for small to moderate datasets and low search complexity.
    -  Search logic is abstracted in a service layer, enabling migration to Elasticsearch (or similar) for advanced capabilities such as full-text search, relevance ranking, autocomplete, and near real-time indexing with minimal code changes.

- Configuration Management
    - Environment-based configuration to allow seamless deployment in different environments (dev, staging, production).
    - Can be extended to support dynamic runtime configuration updates.

### Low-Level Design (LLD)

#### Package responsibilities

- models — domain structs (Service, ServiceVersion). GORM tags + JSON tags.

- repositories — ServiceRepository interface + serviceRepo impl (GORM). Query building, pagination, allowed sort fields.

- services — ServiceService interface + impl. Business rules & validation; call repo; wrap errors in AppError.

- handlers — HTTP handlers; parse query/path/body params; call service; return APIResponse or call c.Error() for middleware.

- middleware — Logger, Recovery, ErrorHandler (translates AppError → HTTP response).

- errors — AppError type (StatusCode, Message, Debug) + pre-defined errors.

- config — load environment, return Config struct.

- main — InitializeApp() does DI wiring: db → repo → service → handler → router.

#### Database schema (relational)

- services table: id (PK), name (unique,index), description, created_at, updated_at

- service_versions table: id (PK), service_id (FK → services.id), version (unique per service), created_at

- Add FK constraints: ON DELETE CASCADE

#### Error flow
- Repo returns *AppError or wraps underlying error.

- Service may convert repo errors to AppError (e.g., NotFound, InvalidInput).

- Handler when encountering an error does c.Error(err) and returns.

- ErrorHandler middleware inspects c.Errors, if first error is *errors.AppError uses StatusCode and Message to respond; logs Debug.

- Unknown errors → middleware logs stack trace and returns generic 500 message.

### Extensibility & Future Scope
- Authentication & Authorization: Add secure login and role-based access so only authorized users can access or modify data.

- Full CRUD Support: Allow creating, updating, and deleting services and versions in addition to reading them.

- Bulk Operations: Support adding or updating multiple services at once to handle large datasets efficiently.

- Advanced Search: Enable full-text search to let users find services quickly by name, description, or keywords.

- Error Handling Extensibility: Currently, errors use a single standard format, but future requirements may need multiple error structures for different contexts. This will ensure more precise error representation, better debugging, and improved client-side handling.

- Caching: Introduce a caching layer (e.g., Redis) to speed up responses and reduce database load.

- Observability: Add metrics, logging, and tracing to monitor performance and troubleshoot issues faster.

- Testing: Implement unit and integration tests to ensure reliability and prevent future bugs.

- Soft Deletes & Version History: Maintain historical versions even after deletion for rollback and auditing purposes.

- API Documentation: Add OpenAPI/Swagger-based interactive documentation for easier onboarding of new developers and consumers.

### Error Handling
- Uses a structured error model (AppError) that separates:

    - HTTP status code

    - User-friendly error message

    - Internal debug info (not exposed to clients)

- Centralized error handling middleware ensures consistent API responses

- Errors are wrapped with context to help with debugging

### Testing Plan

#### Testing Objective

- Validate correctness of business logic.

- Ensure database interactions work as expected.

- Verify HTTP handlers respond properly.

- Confirm error handling and edge cases.

- Enable safe refactoring by providing regression tests.

#### Test Cases 

##### Fetch Service by ID
- Fetch Service by Valid ID

    - Input: Existing service ID
    - Expect: Return service data with all associated versions, HTTP 200 OK

- Fetch Service by Non-Existent ID

    - Input: ID that does not exist in DB
    - Expect: Return HTTP 404 Not Found with appropriate error message

- Fetch Service by Invalid ID (e.g., negative or zero)

    - Input: Invalid ID format or values
    - Expect: Return HTTP 400 Bad Request with validation error

- Fetch Service by ID with DB Error (e.g., connection lost)
    
    - Simulate DB failure during fetch
    - Expect: Return HTTP 500 Internal Server Error with generic error message


##### List Services with Filters and Pagination
- List Services Without Filters

    - Input: Empty filter parameters
    - Expect: Return first page with default limit, sorted by default field and order

- List Services with Valid Filters (name, sort, order, page, limit)

    - Input: Valid filter values
    - Expect: Return filtered, sorted, and paginated list of services

- List Services with Invalid Filter Values

    - Input: Invalid sort field, order, page, or limit (e.g., negative page)
    - Expect: Defaults applied; no errors, valid response

-   List Services When No Services Exist

    - Input: Empty DB
    - Expect: Return empty data array, HTTP 200 OK

-   List Services with Partial DB Failure (e.g., during count or fetch)
    - Simulate DB query failure
    - Expect: Return HTTP 500 Internal Server Error

