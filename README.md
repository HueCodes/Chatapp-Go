# Chatapp-Go

[![Go Version](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Overview

Chatapp-Go is a lightweight, real-time chat application built with Go. It leverages WebSockets for efficient, bidirectional communication, enabling seamless messaging between users. Designed for scalability and ease of use, this project serves as a foundation for building interactive chat systems in web or mobile applications.

Key principles include simplicity in architecture, robust error handling, and modular code structure, making it ideal for developers exploring concurrent programming in Go or prototyping chat functionalities.

## Features

- **Real-Time Messaging**: Supports instant message delivery using WebSockets for low-latency communication.
- **User Authentication**: Basic JWT-based authentication to manage user sessions securely.
- **Room Management**: Create and join chat rooms for group conversations.
- **Message Persistence**: Optional integration with a database (e.g., SQLite) for storing chat history.
- **Concurrent Handling**: Utilizes Go's goroutines and channels for efficient multi-user support.

## Prerequisites

- Go 1.21 or later
- A modern web browser for testing the client-side interface
- (Optional) PostgreSQL or SQLite for persistent storage

## Installation

1. **Clone the Repository**:
   ```
   git clone https://github.com/HueCodes/Chatapp-Go.git
   cd Chatapp-Go
   ```

2. **Install Dependencies**:
   Go modules are used for dependency management. Run:
   ```
   go mod tidy
   ```

3. **Build the Application**:
   ```
   go build -o chatapp cmd/server/main.go
   ```

4. **Set Up Environment** (Optional for Database):
   Create a `.env` file in the root directory and add your database configuration:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=youruser
   DB_PASSWORD=yourpassword
   DB_NAME=chatapp
   JWT_SECRET=your_jwt_secret_key
   ```

## Usage

1. **Run the Server**:
   ```
   ./chatapp
   ```
   The server will start on `http://localhost:8080`.

2. **Access the Chat Interface**:
   Open your browser and navigate to `http://localhost:8080`. Register or log in to start chatting.

3. **API Endpoints** (for custom clients):
   - `POST /register` - User registration
   - `POST /login` - User login
   - `GET /ws` - WebSocket upgrade for real-time chat

For detailed API documentation, refer to the [API docs](docs/api.md) or generate them using tools like Swagger.

Example WebSocket connection in JavaScript:
```javascript
const ws = new WebSocket('ws://localhost:8080/ws?token=your_jwt_token');
ws.onopen = () => console.log('Connected');
ws.onmessage = (event) => console.log('Message:', event.data);
ws.send(JSON.stringify({type: 'message', content: 'Hello, world!'}));
```

## Project Structure

```
Chatapp-Go/
├── cmd/
│   └── server/
│       └── main.go          # Entry point for the server
├── internal/
│   ├── auth/                # Authentication logic
│   ├── handlers/            # HTTP and WebSocket handlers
│   ├── models/              # Data models
│   └── services/            # Business logic services
├── pkg/
│   └── websocket/           # WebSocket utilities
├── docs/                    # Documentation
├── go.mod                   # Go modules
├── go.sum                   # Dependency checksums
└── README.md                # This file
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'Add amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a Pull Request.

Ensure code adheres to Go best practices and includes tests. Run `go test ./...` before submitting.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For questions or feedback, open an issue or reach out to the maintainer at [HueCodes](https://github.com/HueCodes). 

---

*Last updated: November 04, 2025*
