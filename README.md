# Go Chat App

A real-time chat application built with Go, WebSockets, and vanilla JavaScript.

## Features

- **Real-time messaging** - Messages are delivered instantly using WebSockets
- **User management** - Users can set custom usernames and see join/leave notifications
- **Message history** - New users see the last 50 messages when they join
- **Responsive design** - Works on desktop and mobile devices
- **Connection handling** - Automatic reconnection on connection loss
- **Clean UI** - Modern, intuitive interface with message timestamps

## Technology Stack

- **Backend**: Go with Gorilla WebSocket and Mux router
- **Frontend**: Vanilla JavaScript, HTML5, CSS3
- **Real-time Communication**: WebSockets
- **Architecture**: Hub-based message broadcasting

## Project Structure

```
chatapp/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── handlers/
│   ├── hub.go             # WebSocket hub for managing connections
│   ├── client.go          # WebSocket client connection handling
│   └── handlers.go        # HTTP route handlers
├── models/
│   └── message.go         # Message and user data structures
└── static/
    ├── style.css          # Application styles
    └── app.js             # Frontend JavaScript logic
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Git (for cloning dependencies)

### Installation & Running

1. **Clone or navigate to the project directory**
   ```bash
   cd /Users/hugh/Documents/Projects/Chatapp-Go
   ```

2. **Initialize Go modules and install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

4. **Open your browser**
   Navigate to `http://localhost:8080`

The server will start on port 8080. You can open multiple browser tabs to test the chat functionality.

## Usage

1. **Join the chat**: Enter a username when prompted
2. **Send messages**: Type in the message box and press Enter or click Send
3. **Change username**: Click the "Change Username" button in the header
4. **Multiple users**: Open multiple browser tabs/windows to simulate multiple users

## Configuration

### Default Settings

- **Server Port**: 8080
- **Message History**: Last 100 messages stored in memory
- **Recent Messages**: Last 50 messages sent to new users
- **Max Message Length**: 500 characters
- **Max Username Length**: 20 characters
- **WebSocket Settings**:
  - Read timeout: 60 seconds
  - Write timeout: 10 seconds
  - Ping interval: 54 seconds
  - Max message size: 512 bytes

### Customization

You can modify these settings in the respective files:

- **Server port**: Change in `main.go`
- **Message limits**: Modify constants in `handlers/client.go`
- **UI styling**: Edit `static/style.css`
- **Frontend behavior**: Update `static/app.js`

## API Endpoints

### HTTP Routes

- `GET /` - Serves the main chat interface
- `GET /ws?username=<name>` - WebSocket endpoint for chat connections
- `GET /static/*` - Serves static files (CSS, JS)

### WebSocket Messages

#### Client to Server
```json
{
  "type": "text",
  "content": "Hello, world!",
  "username": "john_doe"
}
```

#### Server to Client
```json
{
  "id": "uuid-here",
  "type": "text|user_join|user_left|system",
  "username": "john_doe",
  "content": "Hello, world!",
  "timestamp": "2025-11-03T10:30:00Z"
}
```

## Development

### Adding New Features

1. **New message types**: Add to `models/message.go` and handle in `handlers/hub.go`
2. **Additional routes**: Add to `main.go` and implement in `handlers/handlers.go`
3. **Frontend features**: Extend `static/app.js` and style in `static/style.css`

### Security Considerations

For production deployment, consider adding:

- User authentication and authorization
- Rate limiting for messages
- Input validation and sanitization
- HTTPS/WSS encryption
- CORS configuration
- Database persistence
- Message moderation

## Troubleshooting

### Common Issues

1. **"Connection failed"**: Ensure the server is running on port 8080
2. **"Module not found"**: Run `go mod tidy` to install dependencies
3. **WebSocket errors**: Check browser console for detailed error messages
4. **Port already in use**: Kill existing processes or change port in `main.go`

### Debug Mode

To enable debug logging, add this to `main.go`:
```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

## License

This project is open source and available under the MIT License.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

---

**Happy chatting! 🚀**