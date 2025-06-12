# Go Payroll Processing System 🧾

A payroll management system built with **Go (Gin Gonic)**, **PostgreSQL**, and **Redis**.  
It supports employee attendance tracking, overtime, reimbursements, payslip generation, and admin reporting.

## 🌟 Features

- JWT-based authentication (Admin & Employee roles)
- Monthly payroll periods with prorated salary
- Overtime and reimbursement handling
- Payslip generation with detailed breakdowns
- Admin summary reporting
- Redis caching for performance
- REST API using Gin Gonic

---

## 🛠 Tech Stack

- **Backend:** Go + Gin
- **Database:** PostgreSQL
- **Cache:** Redis
- **Containerized with:** Docker Compose

---

## 📁 Project Structure

```
.
├── app/
├── config/
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

### Build and run with Docker Compose

```bash
docker-compose up --build
```

- App will be available at: `http://localhost:8010`
- PostgreSQL will run on port `5432`
- Redis will run on port `6379`

---

## ✅ Health Check

Once the services are up:

- Visit `http://localhost:8010/health` to test the API
- Use tools like Postman or cURL to interact with other endpoints

---

## 🧪 Common Issues

### ❌ Server returned nothing

Make sure Gin binds to `0.0.0.0` in `main.go`:

```go
r.Run("0.0.0.0:" + port)
```

### ❌ Can't connect to DB or Redis?

Ensure `.env` variables use container service names, like:

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
