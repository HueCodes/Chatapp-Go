# Chatapp-Go

[![Go Version](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Overview

Chatapp-Go is a lightweight, real-time chat application built with Go. It leverages WebSockets for efficient, bidirectional communication, enabling seamless messaging between users. Designed for scalability and ease of use, this project serves as a foundation for building interactive chat systems in web or mobile applications.

Key principles include simplicity in architecture, robust error handling, and modular code structure, making it ideal for developers exploring concurrent programming in Go or prototyping chat functionalities.

## Features

- **Real-Time Messaging**: Supports instant message delivery using WebSockets for low-latency communication.
- **User Authentication**: JWT-based authentication for secure user registration and login.
- **Protected WebSocket Connections**: All WebSocket connections require valid JWT tokens.
- **Message Persistence**: Optional integration with a database (e.g., SQLite) for storing chat history.
- **Concurrent Handling**: Utilizes Go's goroutines and channels for efficient multi-user support.
- **In-Memory User Store**: Simple user storage (easily replaceable with database backend).

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

4. **Set Up Environment** (Optional):
    Create a `.env` file in the root directory and add your configuration:
    ```
    JWT_SECRET=your-super-secret-key-change-in-production
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=youruser
    DB_PASSWORD=yourpassword
    DB_NAME=chatapp
    ```

## Usage

1. **Run the Server**:
   ```
   ./chatapp
   ```
   The server will start on `http://localhost:8080`.

2. **Access the Chat Interface**:
    Open your browser and navigate to `http://localhost:8080`. You'll be prompted to register or login before you can start chatting.

3. **API Endpoints**:
    - `POST /api/register` - User registration (username, email, password)
    - `POST /api/login` - User login (username, password)
    - `GET /ws` - WebSocket upgrade for real-time chat (requires JWT token)

For detailed authentication documentation, see [AUTH.md](AUTH.md).

Example WebSocket connection in JavaScript:
```javascript
// First, get a token by logging in
const loginResponse = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username: 'youruser', password: 'yourpass' })
});
const { token } = await loginResponse.json();

// Then connect to WebSocket with the token
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
ws.onopen = () => console.log('Connected');
ws.onmessage = (event) => console.log('Message:', event.data);
ws.send(JSON.stringify({type: 'text', content: 'Hello, world!'}));
```

## Project Structure

```
Chatapp-Go/
├── auth/                    # Authentication logic (JWT, password hashing)
├── handlers/                # HTTP and WebSocket handlers
├── models/                  # Data models (User, Message)
├── static/                  # Frontend files (HTML, CSS, JS)
├── store/                   # User storage (in-memory, easily replaceable)
├── go.mod                   # Go modules
├── go.sum                   # Dependency checksums
├── main.go                  # Entry point for the server
├── AUTH.md                  # Authentication documentation
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
