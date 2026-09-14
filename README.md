# AWS Employee Management System (AWSEMS)

A full-stack, cloud-native Employee Management System designed to be deployed on AWS. The application uses a modular architecture featuring a React frontend, a Go-based proxy/authentication backend, and Python AWS Lambdas for core business logic, utilizing MongoDB as the database.

## 🏗 Architecture Overview

1. **Frontend (React + Vite + TailwindCSS)**
   - Minimal and clean UI for managing employees.
   - Handles client-side validation and seamlessly communicates with the Go Backend.
2. **Backend (Go / net/http)**
   - Acts as a secure Gateway/BFF (Backend-For-Frontend).
   - Manages **AWS Cognito** OAuth2 authentication (Login, Callback, Logout) securely storing JWTs in `HttpOnly` cookies.
   - Employs **AWS Systems Manager (SSM) Parameter Store** to fetch environment variables seamlessly in the cloud.
   - Validates incoming API payloads explicitly before safely proxying valid requests down to the AWS API Gateway.
3. **Serverless Logic (Python AWS Lambdas)**
   - Python-based serverless functions for each discrete CRUD action (`create`, `delete`, `get`, `search`, `update`).
   - Built with **Pydantic** to enforce strict data modeling and business logic validation.
4. **Database**
   - **MongoDB** is used as the persistent data store for all employee records.

---

## 📂 Project Structure

```text
AWSEMS/
├── frontend/                 # React + Vite Frontend
│   ├── src/
│   │   ├── components/       # Reusable UI elements (EmployeeModal, Table, Toasts)
│   │   ├── pages/            # View components (LoginPage, EmployeesPage)
│   │   ├── services/         # Axios API communication
│   │   └── ...
│   └── package.json          
├── backend/                  # Go Backend (BFF & Auth Gateway)
│   ├── cmd/server/           # Main Go entry point (main.go)
│   ├── internal/
│   │   ├── apigateway/       # HTTP client connecting to AWS API Gateway
│   │   ├── cognito/          # AWS Cognito OIDC configuration & JWT verification
│   │   ├── config/           # Environment and SSM Parameter Store loader
│   │   ├── handler/          # Route handlers (Auth & Employee proxies)
│   │   ├── middleware/       # CORS and logging interceptors
│   │   ├── model/            # Go Structs representing internal state
│   │   └── utils/            # Validation helpers and AWS SDK integrations
│   └── go.mod
└── LambdaPython/             # Serverless Business Logic
    ├── lambdas/
    │   └── employees/        # Python handlers for API Gateway integrations
    ├── models/               # Pydantic models (employee.py, response.py)
    ├── db/                   # MongoDB connection logic
    └── requirements.txt
```

---

## 🚀 Getting Started

### 1. Prerequisites
- **Node.js** (v18+)
- **Go** (v1.22+)
- **Python** (v3.9+)
- **AWS Account** (Cognito User Pool, API Gateway, Lambda, SSM configured)
- **MongoDB** instance

### 2. Backend Setup (Go)
The Go backend requires various environment variables. Locally, create a `.env` file in the `backend/` directory:

```env
APP_ENV=local
APP_PORT=8080

# Cognito Configuration
COGNITO_REGION=us-east-1
COGNITO_USER_POOL_ID=your-pool-id
COGNITO_CLIENT_ID=your-client-id
COGNITO_CLIENT_SECRET=your-client-secret
COGNITO_DOMAIN=your-cognito-domain.auth.us-east-1.amazoncognito.com
COGNITO_REDIRECT_URI=http://localhost:8080/api/auth/callback
COGNITO_LOGOUT_URI=http://localhost:5173/login
COGNITO_OPENID_CONFIG_URL=https://cognito-idp.us-east-1.amazonaws.com/your-pool-id/.well-known/openid-configuration

# Integrations
API_GATEWAY_BASE_URL=https://your-api-gateway-id.execute-api.us-east-1.amazonaws.com/prod
FRONTEND_URL=http://localhost:5173
```
*Note: When deployed to EC2, the backend will automatically pull these variables from AWS SSM Parameter Store if the local `.env` values are absent.*

Run the Go server:
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### 3. Frontend Setup (React)
The frontend communicates directly with the Go backend at `http://localhost:8080`.
```bash
cd frontend
npm install
npm run dev
```
The app will be accessible at `http://localhost:5173`.

### 4. AWS Lambda Setup (Python)
Dependencies for the Python lambdas can be installed via:
```bash
cd LambdaPython
pip install -r requirements.txt
```
To deploy, package the scripts within `LambdaPython/lambdas/employees/` along with their dependencies into `.zip` files and upload them to AWS Lambda, placing them behind your API Gateway.

---

## 🛡️ Authentication & Authorization Flow
1. User attempts to access the frontend and clicks "Sign in".
2. The React frontend redirects to `GET /api/auth/login` on the Go Backend.
3. The Go backend generates an OAuth2 session, stores state in a temporary cookie, and redirects the user to the AWS Cognito Hosted UI.
4. User authenticates via Cognito and is redirected back to the backend callback (`/api/auth/callback`) with an authorization code.
5. The backend verifies the code, issues an `access_token` and `refresh_token` as strictly secure `HttpOnly` cookies, and redirects the user back to the React app.
6. Subsequent requests to protected Go backend routes (e.g., `POST /api/employees`) are intercepted by a JWT validation middleware before being proxied to the AWS API Gateway.

---

## ✨ Features
- **Centralized Validation**: Strict validation exists cleanly on three separate layers (React Client -> Go Gateway -> Python Lambda).
- **Error Passthrough**: Errors thrown deeply in Python lambdas are gracefully surfaced all the way to the frontend Toast notifications.
- **SSM Integrated Configs**: Go backend effortlessly transitions from local `.env` files to cloud-native AWS Systems Manager Parameter Store lookups.
- **Diff-Only Updates**: The frontend selectively patches only the properties a user has changed.

