# 🏦 Simple Bank API

A backend banking system simulation built with **Go (Golang)**, focusing on **Financial Data Integrity, Concurrency Control**, and **High Performance.**

## 🛠 Tech Stack
- **Language:** Go (Golang)
- **Framework:** Gin
- **Database:** PostgreSQL
- **Infrastructure:** Docker, Docker Compose
- **Testing:** Gomock, Testify
- **Security:** PASETO (Platform-Agnostic Security Tokens), Bcrypt
- **Communication:** gRPC & RESTful API

## ✨ Features Overview

### 🏆 Key Highlights (Technical Challenges)

- **Deadlock Prevention:** Implemented consistent lock ordering strategy to prevent database deadlocks during concurrent money transfers.

- **ACID Transactions:** Ensures atomicity for all financial operations (Transfer Record -> Entries -> Balance Update). If any step fails, the entire transaction rolls back.

- **Concurrency Control:** Utilizes **Pessimistic Locking (SELECT ... FOR UPDATE)** to handle high-concurrency balance updates without race conditions.

- **Dual Protocol:** Supports both **gRPC** (for internal high-performance communication) and **RESTful API** (for frontend clients) serving the same business logic.


### 📦 Core Modules
- **💸 Money Transfer System**
  - **Atomic Transfers:** Perform money transfers between accounts within a single database transaction.
  - **Audit Trail:** Automatically generates double-entry bookkeeping records (Entries) for every transaction.
  - **Currency Validation:** Enforces strict currency matching rules before processing transfers.

- **👤 Account Management**
  - Create and manage bank accounts with support for multiple currencies (USD, THB, etc.).
  - Secure balance inquiries with ownership validation.

- **🔐 Authentication & Security**
  - **PASETO Tokens:** Uses Platform-Agnostic Security Tokens (PASETO) for enhanced security over standard JWT.
  - Role-based access control ensuring users can only access their own accounts.

## 🚀 How to Run

```
## Clone Repository
git clone https://github.com/codepnw/simple-bank.git
cd simple-bank

## Setup Environment
cp -n .env.example .env

## Run with Docker
docker compose up --build -d
```
The API Server will start at
  - **HTTP (REST)**: `http://localhost:8080/api/v1`
  - **gRPC**: `http://localhost:9090`

**NOTE**: The database will be automatically initialized with schema and seed data located in **./scripts**

## 📨 API Documentation
**Swagger UI:** `http://localhost:8080/swagger/index.html` 

**Postman Collection:** `./docs/postman/postman_collection.json`

## 🔐 Default Test Accounts

The database comes pre-filled with the following accounts for testing concurrency and transfers:


| Email                   | Password | Initial Balance (Raw) | Display Value    | Currency | Account ID     |
| :---                    | :---     | :---                  | :---             | :---     | :---           |
| **`user1@example.com`** | `123456` | `500000`              | **5,000.00**     | THB      | **1, 2 (USD)** |
| **`user2@example.com`** | `123456` | `0`                   | **0.00**         | THB      | **3**          |
| **`rich@example.com`**  | `123456` | `100000000`           | **1,000,000.00** | THB      | **4**          |

### 💰 Currency & Amount Handling

To ensure precision and avoid floating-point errors, all monetary values in this system are stored as **integers** representing the smallest currency unit (e.g., Cents, Satang).

* **Database:** `BIGINT` (e.g., `100` = 1.00 USD)
* **API Response:** Returns raw integer values. The client is responsible for formatting.

**Example:**
* Balance: `500000` (Satang) = **5,000.00 THB**
* Amount: `100` (Cents) = **1.00 USD**
