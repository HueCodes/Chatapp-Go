# API Usage Examples

## Authentication

### Register a new user
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "token": "eyJhbGc...",
  "username": "alice",
  "user_id": "uuid-here"
}
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "password123"
  }'
```

## Room Management

### List all rooms
```bash
curl http://localhost:8080/api/rooms
```

### Create a room
```bash
curl -X POST http://localhost:8080/api/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tech Talk"
  }'
```

### Get a specific room
```bash
curl http://localhost:8080/api/rooms/1
```

## WebSocket Connection Examples

### JavaScript/Browser Example

```javascript
// After login, you'll have a token
const token = "your-jwt-token";
const roomId = 1; // General room

// Connect to WebSocket
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}&room_id=${roomId}`);

ws.onopen = () => {
  console.log('Connected to chat room', roomId);
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  if (message.type === 'typing') {
    // Handle typing indicator
    console.log(`${message.username} is ${message.is_typing ? 'typing' : 'not typing'}...`);
  } else if (message.type === 'text') {
    // Handle chat message
    console.log(`${message.username}: ${message.content}`);
  } else if (message.type === 'user_join') {
    console.log(`${message.username} joined the room`);
  } else if (message.type === 'user_left') {
    console.log(`${message.username} left the room`);
  }
};

// Send a text message
function sendMessage(content) {
  ws.send(JSON.stringify({
    type: 'text',
    content: content
  }));
}

// Send typing indicator
let typingTimeout;
function handleTyping() {
  // User started typing
  ws.send(JSON.stringify({
    type: 'typing',
    is_typing: true
  }));
  
  // Clear previous timeout
  clearTimeout(typingTimeout);
  
  // Stop typing indicator after 3 seconds
  typingTimeout = setTimeout(() => {
    ws.send(JSON.stringify({
      type: 'typing',
      is_typing: false
    }));
  }, 3000);
}

// When message is sent, stop typing indicator
function onMessageSent() {
  clearTimeout(typingTimeout);
  ws.send(JSON.stringify({
    type: 'typing',
    is_typing: false
  }));
}

// Example: Send a message
sendMessage("Hello, everyone!");
```

### Node.js Example

```javascript
const WebSocket = require('ws');

const token = 'your-jwt-token';
const roomId = 1;

const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}&room_id=${roomId}`);

ws.on('open', () => {
  console.log('Connected to room', roomId);
  
  // Send a message
  ws.send(JSON.stringify({
    type: 'text',
    content: 'Hello from Node.js!'
  }));
  
  // Simulate typing
  setTimeout(() => {
    ws.send(JSON.stringify({
      type: 'typing',
      is_typing: true
    }));
  }, 1000);
  
  setTimeout(() => {
    ws.send(JSON.stringify({
      type: 'typing',
      is_typing: false
    }));
  }, 4000);
});

ws.on('message', (data) => {
  const message = JSON.parse(data);
  console.log('Received:', message);
});

ws.on('error', (error) => {
  console.error('WebSocket error:', error);
});

ws.on('close', () => {
  console.log('Disconnected from chat');
});
```

### Python Example

```python
import asyncio
import websockets
import json

async def chat_client(token, room_id=1):
    uri = f"ws://localhost:8080/ws?token={token}&room_id={room_id}"
    
    async with websockets.connect(uri) as websocket:
        print(f"Connected to room {room_id}")
        
        # Send a message
        await websocket.send(json.dumps({
            "type": "text",
            "content": "Hello from Python!"
        }))
        
        # Start typing
        await websocket.send(json.dumps({
            "type": "typing",
            "is_typing": True
        }))
        
        # Receive messages
        async for message in websocket:
            data = json.loads(message)
            
            if data['type'] == 'typing':
                status = "typing" if data['is_typing'] else "stopped typing"
                print(f"{data['username']} is {status}")
            elif data['type'] == 'text':
                print(f"{data['username']}: {data['content']}")
            elif data['type'] in ['user_join', 'user_left']:
                print(data['content'])

# Run the client
# token = "your-jwt-token"
# asyncio.run(chat_client(token))
```

## Testing Multiple Rooms

```javascript
// Client 1 connects to room 1
const ws1 = new WebSocket('ws://localhost:8080/ws?token=token1&room_id=1');

// Client 2 connects to room 2
const ws2 = new WebSocket('ws://localhost:8080/ws?token=token2&room_id=2');

// Message sent in room 1 will NOT be received by client in room 2
ws1.onopen = () => {
  ws1.send(JSON.stringify({
    type: 'text',
    content: 'This is only in room 1'
  }));
};

ws2.onmessage = (event) => {
  // This will NOT receive the message from room 1
  console.log('Room 2 message:', JSON.parse(event.data));
};
```

## Complete Flow Example

```javascript
// 1. Register
const registerResponse = await fetch('http://localhost:8080/api/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'alice',
    email: 'alice@example.com',
    password: 'password123'
  })
});
const { token } = await registerResponse.json();

// 2. Get available rooms
const roomsResponse = await fetch('http://localhost:8080/api/rooms');
const rooms = await roomsResponse.json();
console.log('Available rooms:', rooms);

// 3. Create a new room
const newRoomResponse = await fetch('http://localhost:8080/api/rooms', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ name: 'My Private Room' })
});
const newRoom = await newRoomResponse.json();

// 4. Connect to the new room
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}&room_id=${newRoom.id}`);

ws.onopen = () => {
  console.log(`Connected to room: ${newRoom.name}`);
  
  // 5. Send a message
  ws.send(JSON.stringify({
    type: 'text',
    content: 'First message in my new room!'
  }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  // You'll receive:
  // - Historical messages (last 50)
  // - Your join message
  // - Your sent message
  // - Other users' messages
  // - Typing indicators
  
  console.log('Message:', message);
};
```

## Error Handling

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?token=invalid&room_id=999');

ws.onerror = (error) => {
  console.error('Connection error:', error);
  // Possible reasons:
  // - Invalid token (401 Unauthorized)
  // - Room doesn't exist (404 Not Found)
  // - Missing token (401 Unauthorized)
};

ws.onclose = (event) => {
  console.log('Connection closed:', event.code, event.reason);
  // Implement reconnection logic here
};
```
