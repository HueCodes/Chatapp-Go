# Authentication Documentation

## Overview

The ChatApp-Go now includes JWT-based authentication to secure user access and WebSocket connections.

## Features

- **User Registration**: New users can create accounts with username, email, and password
- **User Login**: Existing users can authenticate with username and password
- **JWT Tokens**: Secure, stateless authentication using JSON Web Tokens
- **Protected WebSocket**: WebSocket connections require valid JWT tokens
- **Session Persistence**: Tokens are stored in localStorage for persistent sessions
- **Password Security**: Passwords are hashed using bcrypt before storage

## API Endpoints

### Register a New User

```bash
POST /api/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword123"
}

# Response (201 Created)
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "john_doe",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Validation:**
- Username, email, and password are required
- Password must be at least 6 characters
- Username must be unique

### Login

```bash
POST /api/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "securepassword123"
}

# Response (200 OK)
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "john_doe",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### WebSocket Connection

```javascript
// Connect with token as query parameter
const ws = new WebSocket('ws://localhost:8080/ws?token=YOUR_JWT_TOKEN');

// Or use Authorization header (when supported by client)
const ws = new WebSocket('ws://localhost:8080/ws', {
  headers: {
    'Authorization': 'Bearer YOUR_JWT_TOKEN'
  }
});
```

## Token Details

- **Algorithm**: HS256 (HMAC with SHA-256)
- **Expiration**: 24 hours from issuance
- **Claims**:
  - `user_id`: Unique user identifier
  - `username`: User's username
  - `exp`: Expiration timestamp
  - `iat`: Issued at timestamp

## Configuration

Set the JWT secret in your environment:

```bash
# .env file
JWT_SECRET=your-super-secret-key-change-this-in-production
```

Or export as environment variable:

```bash
export JWT_SECRET="your-super-secret-key-change-this-in-production"
```

**Important**: Use a strong, random secret in production!

## Security Best Practices

1. **JWT Secret**: Always use a strong, random secret (at least 32 characters)
2. **HTTPS**: Use HTTPS/WSS in production to prevent token interception
3. **Token Storage**: Tokens are stored in localStorage (consider HttpOnly cookies for enhanced security)
4. **Password Requirements**: Minimum 6 characters (consider stricter requirements for production)
5. **Token Expiration**: Tokens expire after 24 hours (adjust as needed)

## User Storage

Currently, users are stored **in-memory** for simplicity. This means:
- User data is lost when the server restarts
- Not suitable for production use

### For Production

Consider implementing persistent storage:
- PostgreSQL
- MySQL
- MongoDB
- SQLite (for smaller deployments)

Example migration path:
1. Implement `UserStore` interface with your chosen database
2. Replace `store.NewUserStore()` in `main.go`
3. Add database connection configuration

## Frontend Integration

The frontend automatically handles authentication:

1. **First Visit**: Shows login/register modal
2. **Successful Auth**: Stores token in localStorage
3. **WebSocket Connection**: Automatically includes token
4. **Logout**: Clears token and disconnects
5. **Token Persistence**: Remains logged in across page refreshes

## Testing

### Test Registration

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

### Test Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### Test WebSocket with Token

```javascript
// In browser console after logging in
const token = localStorage.getItem('token');
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
ws.onopen = () => console.log('Connected!');
ws.onmessage = (e) => console.log('Message:', e.data);
```

## Error Handling

- **401 Unauthorized**: Invalid or expired token
- **400 Bad Request**: Invalid request body or missing required fields
- **409 Conflict**: Username already exists (registration)
- **500 Internal Server Error**: Server-side error

## Future Enhancements

Potential improvements:
- Email verification
- Password reset functionality
- Refresh tokens for extended sessions
- Rate limiting on authentication endpoints
- Two-factor authentication (2FA)
- OAuth integration (Google, GitHub, etc.)
- User profile management
- Password strength requirements
