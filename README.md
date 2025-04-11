# 🌀 GoFiber Microservices Project

This project is a microservice-based application built with [GoFiber](https://gofiber.io/) consisting of:

- 🔐 **Auth Service** — Handles OAuth2 login & token validation  
- 📊 **Data Service** — Generates and serves time-based random data  
- 💱 **Trade Service** — Handles trade creation with price/time logic  

---

## 🔧 Tech Stack

- **Golang** (GoFiber, GORM)
- **PostgreSQL**
- **Redis (Valkey)**
- **Google Oauth**
- **JWT (`dgrijalva/jwt-go`)**
- **Air** for live reloading during development
---



---

## ⚙️ Installation & Setup

### 1. Clone the Repo

```bash
git clone https://github.com/Nicholas438/maks-go.git
```

---

### 2. Install Dependencies

Make sure Go is installed:  
👉 [Install Go](https://go.dev/doc/install)

#### 🛠️ Required Go Packages

Run this in each service directory (`auth_service`, `trade_service`, `data_service`) to install dependencies:

```bash
go get github.com/gofiber/fiber/v2
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/joho/godotenv
go get github.com/redis/go-redis/v9
go get github.com/dgrijalva/jwt-go
```

Or you can simply tidy everything:
```bash
go mod tidy
```
🔁 Install Air for Hot Reloading
```bash
go install github.com/air-verse/air@latest
```
Make sure $GOPATH/bin is in your system $PATH.

### 3. Setup Environment Variables
Create a .env file inside each service directory

### 4. Run Services with Air
Run each microservice in a separate terminal:
```bash
cd auth_service
air
```
```bash
cd data_service
air
```
```bash
cd trade_service
air
```

# 🧪 API Endpoints

🔐 Auth Service
- `POST /login` – User login and JWT generation (Receives email and password as input)
- `POST /register` – Register a new user (Receives email and password as input)
- `POST /auth` – Validate JWT and return user ID (receives token from bearer authorization header)
- `GET /google-login` – Initiate Google OAuth login
- `GET /google-callback` – Google OAuth callback handler


📊 Data Service
- `GET /bulk-trades-read` – Get all trade data
- `GET /trades-filter-coin-id/:coin_id` – Read filtered data based on coin id
- `POST /create-coin` – Creates a new coin from coin name (receives coin_name)

💱 Trade Service
- `POST /trade` – Create trade (requires Bearer token) (receives price and coin_id)
