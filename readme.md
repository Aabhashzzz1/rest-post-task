## Golang POST API Task

### 📝 Task Description

Create a **POST API** in Go using any Golang framework (**gin**, **mux**, or **gorilla**).  
The API should accept a JSON payload with the following fields:

- `name`: `string`
- `pan`: `string` (PAN number)
- `mobile`: `number` (mobile number)
- `email`: `string` (email ID)

### ✅ Validations

- **PAN** must follow the format: five uppercase letters, followed by four digits, and ending with a letter.  
  Example: `ABCDE1234F`
- **Mobile number** must be a **10-digit** numeric string.
- **Email** must be a valid email address.
- Use the `validator.v10` package to validate **PAN** and **mobile number**.

### 🧩 Middleware

Implement a middleware to **log the API latency** for each request.

---

### 📌 Requirements

1. Create a **router** with a POST endpoint.
2. Implement **middleware** to log the API latency.
3. Use the **validator (v10)** package for PAN and mobile validation.
4. If validation **fails**, respond with an **appropriate error message and HTTP status code**.
5. If validation **passes**, respond with a **success message and status code**.
6. Implement **dependency injection** for better code organization and testability.
7. Write **comprehensive test cases** to cover:
   - Valid inputs
   - Invalid inputs
   - Error handling
   - Middleware functionality

## 🚀 How to Run This Golang Project

### 🛠 Prerequisites

Make sure you have the following installed:

- Go 1.18 or later
- Git

--- 

### 📥 Step-by-Step Instructions

#### 1. Clone the repository

```bash
git clone https://github.com/Aabhashzzz1/rest-post-task.git`
cd rest-post-task
```

2. Initialize Go modules 
```go
go mod tidy
```

3. Run the API Server
```go
go run main.go
```

4. To check test 
```go
go test -v
```
