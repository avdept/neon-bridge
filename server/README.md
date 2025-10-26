# Neon bridge Server

A Go server for managing dashboard widgets with GORM and Gin.


## Setup

1. Initialize Go modules and install dependencies:

```bash
cd server
go mod tidy
```

2. Copy environment file:

```bash
cp .env.example .env
```

3. Run the server:

```bash
go run main.go
```

The server will start on port 8080 by default.
