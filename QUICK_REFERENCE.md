# Quick Reference Guide

## Server Start
```bash
./chatapp
```
- Creates `chatapp.db` SQLite database
- Creates default "General" room (ID: 1)
- Listens on port 8080

## API Endpoints

### Authentication
```bash
# Register
POST /api/register
Body: {"username":"alice","email":"alice@example.com","password":"pass123"}

# Login
POST /api/login
Body: {"username":"alice","password":"pass123"}

# Response: {"token":"jwt-token","username":"alice","user_id":"uuid"}
```

### Rooms
```bash
# List rooms
GET /api/rooms

# Create room
POST /api/rooms
Body: {"name":"Room Name"}

# Get room
GET /api/rooms/{id}
```

### WebSocket
```
ws://localhost:8080/ws?token={jwt-token}&room_id={room-id}
```
- `token` (required): JWT from login/register
- `room_id` (optional): Defaults to 1

## WebSocket Message Types

### Send Text Message
```json
{
  "type": "text",
  "content": "Hello, world!"
}
```

### Send Typing Indicator
```json
{
  "type": "typing",
  "is_typing": true
}
```

### Received Message Format
```json
{
  "id": 123,
  "type": "text|user_join|user_left|typing",
  "user_id": "uuid",
  "username": "alice",
  "room_id": 1,
  "content": "message content",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Received Typing Indicator
```json
{
  "type": "typing",
  "user_id": "uuid",
  "username": "bob",
  "room_id": 1,
  "is_typing": true
}
```

## JavaScript Quick Start
```javascript
// 1. Login
const res = await fetch('/api/login', {
  method: 'POST',
  headers: {'Content-Type': 'application/json'},
  body: JSON.stringify({username: 'alice', password: 'pass123'})
});
const {token} = await res.json();

// 2. Connect to room 1
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}&room_id=1`);

// 3. Handle messages
ws.onmessage = (e) => {
  const msg = JSON.parse(e.data);
  if (msg.type === 'text') console.log(`${msg.username}: ${msg.content}`);
  if (msg.type === 'typing') console.log(`${msg.username} is typing...`);
};

// 4. Send message
ws.send(JSON.stringify({type: 'text', content: 'Hello!'}));

// 5. Send typing indicator
ws.send(JSON.stringify({type: 'typing', is_typing: true}));
```

## Features

### ✅ Persistent Messages
- Messages saved to SQLite (`chatapp.db`)
- Last 50 messages loaded on connect
- Only text messages persisted

### ✅ Chat Rooms
- Multiple isolated chat rooms
- REST API for room management
- Default "General" room (ID: 1)
- Messages broadcast only within rooms

### ✅ Typing Indicators
- Real-time typing status
- Not persisted to database
- Scoped to room
- Not echoed to sender

## Database

### Location
`chatapp.db` in project root

### Tables
- `messages`: id, type, user_id, username, room_id, content, timestamp, deleted_at
- `rooms`: id, name, created_at

### Reset Database
```bash
rm chatapp.db
./chatapp  # Will recreate with default room
```

## Common Tasks

### Create a new room
```bash
curl -X POST http://localhost:8080/api/rooms \
  -H "Content-Type: application/json" \
  -d '{"name":"My Room"}'
```

### List all rooms
```bash
curl http://localhost:8080/api/rooms
```

### Connect to specific room
```javascript
const ws = new WebSocket(`ws://host/ws?token=${token}&room_id=2`);
```

### Implement typing debounce
```javascript
let timeout;
input.addEventListener('input', () => {
  ws.send(JSON.stringify({type: 'typing', is_typing: true}));
  clearTimeout(timeout);
  timeout = setTimeout(() => {
    ws.send(JSON.stringify({type: 'typing', is_typing: false}));
  }, 3000);
});
```

## Error Codes

| Code | Reason                          |
|------|---------------------------------|
| 400  | Invalid request body/parameters |
| 401  | Missing or invalid token        |
| 404  | Room not found                  |
| 409  | Room name already exists        |
| 500  | Server error                    |

## File Structure
```
chatapp/
├── main.go              # Entry point
├── auth/                # JWT & password
├── database/            # DB initialization
├── handlers/            # HTTP & WebSocket handlers
├── models/              # Data models
├── store/               # Data stores
└── static/              # Frontend files
```

## Documentation
- `FEATURES.md` - Detailed feature specs
- `API_EXAMPLES.md` - Code examples
- `IMPLEMENTATION_SUMMARY.md` - Implementation details
- `AUTH.md` - Authentication docs
- `README.md` - Project overview

## Testing
```bash
# Build
go build -o chatapp

# Run
./chatapp

# Test in browser
# 1. Open http://localhost:8080
# 2. Register/login
# 3. Start chatting!
```
