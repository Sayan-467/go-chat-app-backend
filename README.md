# Go Chat App Backend 💬

![Go](https://img.shields.io/badge/Go-1.21-blue)
![Gin](https://img.shields.io/badge/Framework-Gin-green)
![WebSocket](https://img.shields.io/badge/RealTime-WebSockets-orange)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue)
![License](https://img.shields.io/badge/License-MIT-lightgrey)

A **real-time chat backend** built using **Go (Gin Framework)**, **WebSockets**, and **PostgreSQL**.  
Includes authentication, message persistence, and scalable architecture for real-time applications.

## 🚀 Features
- JWT-based **Signup & Login**
- Real-time messaging using WebSockets
- Chat history stored and retrieved from PostgreSQL
- Support for both **public chatrooms** and **1-on-1 DMs**
- Modular project structure for scalability

## 🛠️ Tech Stack
- **Backend:** Go (Gin, GORM)
- **Database:** PostgreSQL
- **WebSockets:** Gorilla WebSocket
- **Auth:** JWT
- **ORM:** GORM

## 🧩 Project Structure
## 🚀 Features
- JWT-based **Signup & Login**
- Real-time messaging using WebSockets
- Chat history stored and retrieved from PostgreSQL
- Support for both **public chatrooms** and **1-on-1 DMs**
- Modular project structure for scalability

## 🛠️ Tech Stack
- **Backend:** Go (Gin, GORM)
- **Database:** PostgreSQL
- **WebSockets:** Gorilla WebSocket
- **Auth:** JWT
- **ORM:** GORM

---

## ⚙️ Setup Instructions

### 1️⃣ Clone the repository
```bash
git clone https://github.com/<your-username>/chat-app-backend.git
cd chat-app-backend
```

### 2️⃣ Setup .env file
```bash
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=chatappDB
```

### 3️⃣ Run the app
```bash
go mod tidy
go run cmd/main.go
```

### Expected output:
🚀 Server running on port: 8080
Database connected and migrated successfully

--- 
## 📜 License
MIT License © 2025 [Syed Sayan]

### Built with ❤️ in Go and PostgreSQL