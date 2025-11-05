class ChatApp {
    constructor() {
        this.ws = null;
        this.username = '';
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectInterval = 3000;
        
        this.initializeElements();
        this.setupEventListeners();
        this.showUsernameModal();
    }

    initializeElements() {
        this.elements = {
            messages: document.getElementById('messages'),
            messageInput: document.getElementById('messageInput'),
            sendButton: document.getElementById('sendButton'),
            usernameModal: document.getElementById('usernameModal'),
            usernameInput: document.getElementById('usernameInput'),
            joinChatButton: document.getElementById('joinChat'),
            usernameDisplay: document.getElementById('username-display'),
            changeUsernameButton: document.getElementById('change-username'),
            userCount: document.getElementById('userCount')
        };
    }

    setupEventListeners() {
        // Send message on button click
        this.elements.sendButton.addEventListener('click', () => this.sendMessage());
        
        // Send message on Enter key
        this.elements.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });

        // Join chat button
        this.elements.joinChatButton.addEventListener('click', () => this.joinChat());
        
        // Username input enter key
        this.elements.usernameInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                this.joinChat();
            }
        });

        // Change username button
        this.elements.changeUsernameButton.addEventListener('click', () => {
            this.disconnect();
            this.showUsernameModal();
        });

        // Handle page visibility change for reconnection
        document.addEventListener('visibilitychange', () => {
            if (!document.hidden && (!this.ws || this.ws.readyState !== WebSocket.OPEN)) {
                this.connect();
            }
        });
    }

    showUsernameModal() {
        this.elements.usernameModal.classList.add('show');
        this.elements.usernameInput.focus();
    }

    hideUsernameModal() {
        this.elements.usernameModal.classList.remove('show');
    }

    joinChat() {
        const username = this.elements.usernameInput.value.trim();
        
        if (!username) {
            alert('Please enter a username');
            return;
        }

        if (username.length > 20) {
            alert('Username must be 20 characters or less');
            return;
        }

        this.username = username;
        this.elements.usernameDisplay.textContent = username;
        this.hideUsernameModal();
        this.connect();
    }

    connect() {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            return;
        }

        this.showConnectionStatus('connecting', 'Connecting...');
        
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws?username=${encodeURIComponent(this.username)}`;
        
        try {
            this.ws = new WebSocket(wsUrl);
            this.setupWebSocketEvents();
        } catch (error) {
            console.error('WebSocket connection error:', error);
            this.handleConnectionError();
        }
    }

    setupWebSocketEvents() {
        this.ws.onopen = () => {
            console.log('Connected to chat server');
            this.reconnectAttempts = 0;
            this.showConnectionStatus('connected', 'Connected');
            this.elements.messageInput.disabled = false;
            this.elements.sendButton.disabled = false;
            
            // Hide connection status after 2 seconds
            setTimeout(() => this.hideConnectionStatus(), 2000);
        };

        this.ws.onmessage = (event) => {
            try {
                // Handle multiple messages separated by newlines
                const messages = event.data.split('\n').filter(msg => msg.trim());
                messages.forEach(msgStr => {
                    const message = JSON.parse(msgStr);
                    this.displayMessage(message);
                });
            } catch (error) {
                console.error('Error parsing message:', error);
            }
        };

        this.ws.onclose = (event) => {
            console.log('Disconnected from chat server');
            this.showConnectionStatus('disconnected', 'Disconnected');
            this.elements.messageInput.disabled = true;
            this.elements.sendButton.disabled = true;
            
            if (!event.wasClean) {
                this.attemptReconnect();
            }
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.handleConnectionError();
        };
    }

    handleConnectionError() {
        this.showConnectionStatus('disconnected', 'Connection failed');
        this.elements.messageInput.disabled = true;
        this.elements.sendButton.disabled = true;
        this.attemptReconnect();
    }

    attemptReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            this.showConnectionStatus('disconnected', 'Connection lost. Please refresh the page.');
            return;
        }

        this.reconnectAttempts++;
        this.showConnectionStatus('connecting', `Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
        
        setTimeout(() => {
            this.connect();
        }, this.reconnectInterval);
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    sendMessage() {
        const content = this.elements.messageInput.value.trim();
        
        if (!content) {
            return;
        }

        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            alert('Not connected to chat server');
            return;
        }

        const message = {
            type: 'text',
            content: content,
            username: this.username
        };

        try {
            this.ws.send(JSON.stringify(message));
            this.elements.messageInput.value = '';
            this.elements.messageInput.focus();
        } catch (error) {
            console.error('Error sending message:', error);
            alert('Failed to send message');
        }
    }

    displayMessage(message) {
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${message.type}`;
        
        if (message.type === 'text') {
            messageDiv.innerHTML = `
                <div class="message-header">
                    <span class="username">${this.escapeHtml(message.username)}</span>
                    <span class="timestamp">${this.formatTimestamp(message.timestamp)}</span>
                </div>
                <div class="message-content">${this.escapeHtml(message.content)}</div>
            `;
        } else {
            // System messages (join/leave)
            messageDiv.innerHTML = `
                <div class="message-content">${this.escapeHtml(message.content)}</div>
                <div class="timestamp">${this.formatTimestamp(message.timestamp)}</div>
            `;
        }

        this.elements.messages.appendChild(messageDiv);
        this.scrollToBottom();
        
        // Update user count (simple estimation based on join/leave messages)
        this.updateUserCount(message);
    }

    updateUserCount(message) {
        // This is a simple client-side estimation
        // In a real app, you'd want the server to send the actual count
        if (message.type === 'user_join' || message.type === 'user_left') {
            // For demo purposes, we'll show a static message
            this.elements.userCount.textContent = 'Users online';
        }
    }

    formatTimestamp(timestamp) {
        const date = new Date(timestamp);
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    scrollToBottom() {
        this.elements.messages.scrollTop = this.elements.messages.scrollHeight;
    }

    showConnectionStatus(status, message) {
        let statusElement = document.querySelector('.connection-status');
        
        if (!statusElement) {
            statusElement = document.createElement('div');
            statusElement.className = 'connection-status';
            document.body.appendChild(statusElement);
        }

        statusElement.className = `connection-status ${status}`;
        statusElement.textContent = message;
        statusElement.style.display = 'block';
    }

    hideConnectionStatus() {
        const statusElement = document.querySelector('.connection-status');
        if (statusElement) {
            statusElement.style.display = 'none';
        }
    }
}

// Initialize the chat app when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new ChatApp();
});