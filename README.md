## Problem Statement:  
A notification service that facilitates the sending of notifications thorugh email, slack, in-app notification   

## Requirements:  
1. Send Notifications through channels like Email, Slack, In-app. 
2. Producers tells notification service. 
    a. to send the message instantly. 
    b. to send the message in future. 
3. Producers can design their own template. 

---

## Project Structure

```
.
├── datastore/        # In-memory datastore
├── handler/          # HTTP handlers for templates and notifications
├── model/            # Data models and enums
├── service/          # Business logic, scheduling, heap operations
├── main.go           # Application entrypoint
├── go.mod
├── go.sum
└── README.md
```

---

## Getting Started

### Prerequisites
- Go 1.18+
- Git

### Clone and Build
```bash
git clone https://github.com/<your-username>/notification-service.git
cd notification-service
go mod tidy
go build
```

### Run the Service
```bash
go run main.go
```

The server will start at:
```
http://localhost:8080
```

---

## API Endpoints

### 1. Create a Template
**POST** `/producer/template`

**Request Body:**
```json
{
  "name": "welcome",
  "subject": "Hi {0}",
  "message": "Welcome {0} to {1}"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Successfully inserted for welcome"
}
```

---

### 2. Send a Notification
**POST** `/producer/notify`

**Request Body (Immediate Send):**
```json
{
  "to": "user@example.com",
  "from": "noreply@example.com",
  "message": {
    "subject": "Hello",
    "body": "This is an immediate message"
  },
  "channel": "email"
}
```

**Request Body (With Template and Delay):**
```json
{
  "to": "user@example.com",
  "from": "noreply@example.com",
  "template": "welcome",
  "time": 25,
  "message": {
    "subplaceholder": ["Alice"],
    "bodyplaceholder": ["Alice", "GoLang Community"]
  },
  "channel": "email"
}
```
