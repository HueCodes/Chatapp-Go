# Chat Rooms, Persistent Messages, and Typing Indicators

This document describes the implementation of three new features in Chatapp-Go:

## 1. Persistent Message Storage (SQLite)

### Database Setup
- **Location**: `chatapp.db` (SQLite file in project root)
- **Schema**: Auto-migrated using GORM
- **Tables**:
  - `messages`: Stores all chat messages with user info and room association
  - `rooms`: Stores chat room information

### Message Model
```go
type Message struct {
    ID        uint            // Primary key
    Type      MessageType     // text, user_join, user_left, system, typing
    UserID    string          // User who sent the message
    Username  string          // Username of sender
    RoomID    uint            // Room where message was sent
    Content   string          // Message content
    Timestamp time.Time       // When message was sent
    DeletedAt gorm.DeletedAt  // Soft delete support
}
```

### Features
- **Automatic persistence**: Text messages are saved to database asynchronously
- **Message history**: New users connecting to a room receive the last 50 messages
- **Only text messages are persisted**: System messages (join/leave) and typing indicators are not saved

### API
No direct API - persistence happens automatically through WebSocket messages.

---

## 2. Chat Rooms

### Room Model
```go
type Room struct {
    ID        uint      // Primary key
    Name      string    // Room name (unique)
    CreatedAt time.Time // When room was created
}
```

### Default Room
- A default "General" room (ID: 1) is created on first startup
- Users connect to room ID 1 by default if no room_id is specified

### REST API Endpoints

#### GET /api/rooms
List all available rooms.

**Response**:
```json
[
  {
    "id": 1,
    "name": "General",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

#### POST /api/rooms
Create a new room.

**Request**:
```json
{
  "name": "Tech Discussion"
}
```

**Response** (201 Created):
```json
{
  "id": 2,
  "name": "Tech Discussion",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Errors**:
- 400: Invalid request or missing name
- 409: Room name already exists

#### GET /api/rooms/{id}
Get a specific room by ID.

**Response**:
```json
{
  "id": 1,
  "name": "General",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Errors**:
- 400: Invalid room ID
- 404: Room not found

### WebSocket Connection with Rooms
Connect to a specific room using the `room_id` query parameter:

```javascript
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}&room_id=1`);
```

- If `room_id` is not provided, defaults to room 1 (General)
- Room must exist before connecting (404 error if not)
- Messages are only broadcast to users in the same room
- Each client can only be in one room at a time per connection

### Room Isolation
- Messages sent in room A are only received by users in room A
- Typing indicators are only sent to users in the same room
- Join/leave messages are scoped to the room

---

## 3. Typing Indicators

### Message Format
Send a typing indicator via WebSocket:

```json
{
  "type": "typing",
  "is_typing": true
}
```

**Fields**:
- `type`: Must be "typing"
- `is_typing`: true when user starts typing, false when they stop

### Behavior
- Typing indicators are NOT persisted to the database
- Only sent to other users in the same room (not echoed back to sender)
- Real-time only - no history is maintained

### Received Format
Other users in the room receive:

```json
{
  "type": "typing",
  "user_id": "user-uuid",
  "username": "alice",
  "room_id": 1,
  "is_typing": true
}
```

### Best Practices
- Send `is_typing: true` when user starts typing
- Send `is_typing: false` after 3 seconds of inactivity or when message is sent
- Use debouncing to avoid flooding the server with typing events

---

## Implementation Details

### Database Package
- **Location**: `database/db.go`
- **Function**: `InitDB(dbPath string)` - Initialize SQLite connection
- **Function**: `AutoMigrate(models...)` - Run GORM migrations

### New Stores
1. **RoomStore** (`store/room_store.go`):
   - Manages room CRUD operations
   - Tracks which clients are in which rooms (in-memory)
   
2. **MessageStore** (`store/message_store.go`):
   - `Save(message)` - Persist message to database
   - `GetByRoom(roomID, limit)` - Retrieve recent messages

### Hub Changes
- Now tracks `roomStore` and `messageStore`
- Messages include `RoomID` field
- Broadcasting is room-scoped
- Separate channel for typing indicators
- Async message persistence (non-blocking)

### Client Changes
- Added `roomID` field to Client struct
- Added `clientID()` method for unique identification
- Message reading now distinguishes between text and typing messages
- Automatically populates message metadata (userID, username, roomID, timestamp)

### WebSocket Handler Changes
- Accepts `room_id` query parameter
- Validates room exists before allowing connection
- Passes roomStore to handler

---

## Testing

### Test Message Persistence
1. Start the server: `./chatapp`
2. Register/login two users
3. Connect to WebSocket and send messages
4. Disconnect and reconnect
5. Verify you receive the last 50 messages

### Test Rooms
1. Create a room: `POST /api/rooms` with `{"name": "Test Room"}`
2. Connect two clients to different rooms
3. Send messages - verify isolation (messages don't cross rooms)
4. List rooms: `GET /api/rooms`

### Test Typing Indicators
1. Connect two clients to the same room
2. From client 1, send: `{"type": "typing", "is_typing": true}`
3. Client 2 should receive the typing indicator
4. Client 1 should NOT receive their own typing indicator

---

## Migration Notes

### Breaking Changes
- Message ID changed from string (UUID) to uint (auto-increment)
- Message structure now requires RoomID
- Hub constructor requires roomStore and messageStore parameters
- WSHandler requires roomStore parameter

### Database File
- Created automatically on first run
- Located at `chatapp.db` in project root
- Added to `.gitignore`

### Dependencies Added
- `gorm.io/gorm` - ORM
- `gorm.io/driver/sqlite` - SQLite driver
- `github.com/mattn/go-sqlite3` - SQLite3 C bindings (indirect)
