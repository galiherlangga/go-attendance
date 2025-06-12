# Go Payroll Processing System 🧾

A payroll management system built with **Go (Gin Gonic)**, **PostgreSQL**, and **Redis**.  
It supports employee attendance tracking, overtime, reimbursements, payslip generation, and admin reporting.

---

## 🌟 Features

- JWT-based authentication (Admin & Employee roles)
- Monthly payroll periods with prorated salary
- Overtime and reimbursement handling
- Payslip generation with detailed breakdowns
- Admin summary reporting
- Redis caching for performance
- REST API using Gin Gonic
- Interactive Swagger documentation

---

## 📚 Prerequisites

Ensure you have the following installed:

- [Go 1.24+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## 🛠 Tech Stack

- **Backend:** Go + Gin
- **Database:** PostgreSQL
- **Cache:** Redis
- **Docs:** Swagger via swaggo
- **Containerized with:** Docker Compose

---

## 📁 Project Structure

```
.
├── app/
├── config/
├── docs/
├── pkg/
├── routes/
├── tests/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── README.md
├── .env.example
└── .env
```

---

## ⚙️ Environment Variables

Create a `.env` file in the root directory:

```env
# PostgreSQL
DB_HOST=db
DB_PORT=5432
DB_USER=dev
DB_PASSWORD=yourpassword
DB_NAME=go-attendance

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASS=password
REDIS_DB=0

# Server
PORT=8010
JWT_SECRET=your_jwt_secret
```

---

## 🐳 Run with Docker Compose

Build and run the application:

```bash
docker-compose up --build
```

- App: `http://localhost:8010`
- Swagger UI: `http://localhost:8010/swagger/index.html`
- PostgreSQL: Port `5432`
- Redis: Port `6379`

---

## 🧪 API Documentation (Swagger UI)

Interactive API docs available at:

```
http://localhost:8010/swagger/index.html
```

If you change handler annotations or models, regenerate docs using:

```bash
swag init
```

---

## 🔐 Authentication

Uses JWT stored in HTTP-only cookies.

### 🔑 Login

```http
POST /auth/login
```

Example Response:

```json
{
  "access_token": "....",
  "refresh_token": "...."
}
```

The `access_token` will be stored in cookies and used for protected endpoints.

If testing in Swagger UI, manually add a cookie:

```
Name: access_token
Value: <your JWT token>
```

---

## 🧰 Development Notes

If you want to run locally without Docker:

```bash
go mod tidy
go run main.go
```

Ensure PostgreSQL and Redis are running with correct `.env` config.

---

## ✅ Health Check

Verify if the API server is running:

```bash
http://localhost:8010/health
```

---

## 🧪 Common Issues

### ❌ Server returned nothing

Ensure Gin binds to `0.0.0.0` in `main.go`:

```go
r.Run("0.0.0.0:" + port)
```

### ❌ Can't connect to DB or Redis?

Ensure your `.env` uses service names in Docker:

```env
DB_HOST=db
REDIS_HOST=redis
```

---

## 🧹 Stopping & Cleaning Up

To stop the Docker containers:

```bash
docker-compose down
```

To remove volumes (Postgres data):

```bash
docker-compose down -v
```

---

## 🙌 Contributing

Feel free to fork, open issues or PRs. This is a work in progress and your feedback is welcome!

---

## 📄 License

MIT License
