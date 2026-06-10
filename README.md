<div align="center">

# Task Manager REST API

**Lightweight REST API for task management built with Go**

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-1.10-00ADD8?style=flat&logo=go&logoColor=white)](https://github.com/gin-gonic/gin)
[![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=flat&logo=mysql&logoColor=white)](https://www.mysql.com/)

[Report a Bug](https://github.com/Cesarks81/Task-manager-REST-API/issues) · [Request a Feature](https://github.com/Cesarks81/Task-manager-REST-API/issues)

</div>

---

## Table of Contents

- [Description](#description)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Environment Variables](#environment-variables)
- [Usage](#usage)
- [API Reference](#api-reference)
- [Project Structure](#project-structure)

---

## Description

Task Manager REST API is a backend service built with Go and Gin that provides a clean HTTP interface for managing tasks. It supports creating, reading, updating, and deleting tasks, as well as searching by title or status. All data is persisted in a MySQL database using the standard `database/sql` package with full connection pool configuration.

---

## Features

### Task Management
- Create tasks with title, description, and status
- Retrieve all tasks or a single task by ID
- Update task title, description, and status
- Delete tasks by ID
- Search tasks by `title` or `status` via query parameters

### Status Lifecycle
- Three accepted statuses: `new`, `ongoing`, `completed`
- `completedAt` timestamp is automatically set when a task reaches `completed` status
- Input validation rejects unknown or missing statuses

### Data Layer
- Raw SQL queries — no ORM overhead
- Connection pool tuning (max open/idle connections, max lifetime)
- Parameterized queries to prevent SQL injection

---

## Tech Stack

| Layer | Technology | Version |
|---|---|---|
| Language | Go | 1.24.x |
| HTTP framework | Gin | 1.10.x |
| Database | MySQL | 8.0+ |
| MySQL driver | go-sql-driver/mysql | 1.9.x |

---

## Architecture

```
┌─────────────────────────────────────────────────┐
│               HTTP Client / curl                │
└───────────────────────┬─────────────────────────┘
                        │ HTTP / REST
┌───────────────────────▼─────────────────────────┐
│              Gin Router (main.go)               │
│         handlers/task.go — route logic          │
└───────────────────────┬─────────────────────────┘
                        │ database/sql
┌───────────────────────▼─────────────────────────┐
│         db/task_queries.go — SQL layer          │
└───────────────────────┬─────────────────────────┘
                        │ TCP
┌───────────────────────▼─────────────────────────┐
│               MySQL — table: task               │
└─────────────────────────────────────────────────┘
```

The project follows a simple **layered pattern**: the handler registers routes and handles HTTP concerns, the `db` package owns all SQL queries, and `models` defines the shared data structures.

---

## Prerequisites

- **Go** ≥ 1.21
- **MySQL** running locally or remotely
- A database with the `task` table created:

```sql
CREATE TABLE task (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      VARCHAR(50) NOT NULL,
    createdat   DATETIME NOT NULL,
    completedat DATETIME DEFAULT NULL
);
```

---

## Installation & Setup

### 1. Clone the repository

```bash
git clone https://github.com/Cesarks81/Task-manager-REST-API.git
cd Task-manager-REST-API/app
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure environment variables

```bash
# Copy the example file and fill in your credentials
cp .env.example .env
```

### 4. Start the server

```bash
# Export the DSN and run
export DB_DSN="user:password@tcp(127.0.0.1:3306)/datago?parseTime=true&loc=Europe%2fMadrid"
go run .
```

The server will be available at `http://localhost:8080`.

---

## Environment Variables

Create the file `app/.env` with the following variable:

```env
# MySQL connection string (DSN)
DB_DSN=user:password@tcp(127.0.0.1:3306)/datago?parseTime=true&loc=Europe%2fMadrid
```

| Variable | Description |
|---|---|
| `DB_DSN` | Full MySQL Data Source Name including user, password, host, port, and database name |

> **Security note:** Never commit the `.env` file to the repository. It is already included in `.gitignore`.

---

## Usage

### Create a task

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Study Go", "description": "Review concurrency patterns", "status": "new"}'
```

### Get all tasks

```bash
curl http://localhost:8080/api/tasks
```

### Get a task by ID

```bash
curl http://localhost:8080/api/tasks/1
```

### Search by status

```bash
curl "http://localhost:8080/api/tasks/search?status=ongoing"
```

### Search by title

```bash
curl "http://localhost:8080/api/tasks/search?title=study%20go"
```

### Update a task

```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Study Go", "description": "Review concurrency patterns", "status": "completed"}'
```

### Delete a task

```bash
curl -X DELETE http://localhost:8080/api/tasks/1
```

---

## API Reference

### Tasks

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/tasks` | List all tasks |
| `GET` | `/api/tasks/:id` | Get a task by ID |
| `GET` | `/api/tasks/search` | Search tasks by `title` or `status` (query param) |
| `POST` | `/api/tasks` | Create a new task |
| `PUT` | `/api/tasks/:id` | Update an existing task |
| `DELETE` | `/api/tasks/:id` | Delete a task |

### Request body — POST / PUT

```json
{
  "title": "Task title",
  "description": "Task description",
  "status": "new"
}
```

### Valid statuses

| Status | Description |
|---|---|
| `new` | Task has been created but not started |
| `ongoing` | Task is currently in progress |
| `completed` | Task is done — `completedAt` is set automatically |

### Response — GET /api/tasks

```json
[
  {
    "id": 1,
    "title": "Study Go",
    "description": "Review concurrency patterns",
    "status": "completed",
    "createdAt": "2025-06-01T10:00:00Z",
    "completedAt": "2025-06-02T15:30:00Z"
  }
]
```

---

## Project Structure

```
Task-manager-REST-API/
└── app/
    ├── main.go               # Entry point — DB setup, connection pool, Gin server
    ├── handlers/
    │   └── task.go           # Route registration and HTTP handler logic
    ├── db/
    │   └── task_queries.go   # SQL queries: get, update, delete
    ├── models/
    │   └── tasks.go          # Task struct and accepted status map
    ├── .env.example          # Environment variable template
    ├── go.mod                # Module definition and dependencies
    └── go.sum                # Dependency checksums
```

---

<div align="center">

Developed by [César Ramos Morón](https://github.com/Cesarks81)

</div>
