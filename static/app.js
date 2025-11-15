class ChatApp {
    constructor() {
        this.ws = null;
        this.username = '';
        this.token = '';
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectInterval = 3000;
        
        this.initializeElements();
        this.setupEventListeners();
        this.checkAuth();
    }

    initializeElements() {
        this.elements = {
            messages: document.getElementById('messages'),
            messageInput: document.getElementById('messageInput'),
            sendButton: document.getElementById('sendButton'),
            authModal: document.getElementById('authModal'),
            loginForm: document.getElementById('loginForm'),
            registerForm: document.getElementById('registerForm'),
            loginUsername: document.getElementById('loginUsername'),
            loginPassword: document.getElementById('loginPassword'),
            loginBtn: document.getElementById('loginBtn'),
            registerUsername: document.getElementById('registerUsername'),
            registerEmail: document.getElementById('registerEmail'),
            registerPassword: document.getElementById('registerPassword'),
            registerBtn: document.getElementById('registerBtn'),
            showRegister: document.getElementById('showRegister'),
            showLogin: document.getElementById('showLogin'),
            authTitle: document.getElementById('authTitle'),
            authError: document.getElementById('authError'),
            usernameDisplay: document.getElementById('username-display'),
            logoutBtn: document.getElementById('logout-btn'),
            userCount: document.getElementById('userCount')
        };
    }

    setupEventListeners() {
        this.elements.sendButton.addEventListener('click', () => this.sendMessage());
        
        this.elements.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });

        this.elements.loginBtn.addEventListener('click', () => this.login());
        this.elements.registerBtn.addEventListener('click', () => this.register());
        
        this.elements.showRegister.addEventListener('click', (e) => {
            e.preventDefault();
            this.showRegisterForm();
        });
        
        this.elements.showLogin.addEventListener('click', (e) => {
            e.preventDefault();
            this.showLoginForm();
        });

        this.elements.logoutBtn.addEventListener('click', () => this.logout());

        this.elements.loginPassword.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                this.login();
            }
        });

        this.elements.registerPassword.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                this.register();
            }
        });

        document.addEventListener('visibilitychange', () => {
            if (!document.hidden && (!this.ws || this.ws.readyState !== WebSocket.OPEN)) {
                if (this.token) {
                    this.connect();
                }
            }
        });
    }

    checkAuth() {
        const token = localStorage.getItem('token');
        const username = localStorage.getItem('username');
        
        if (token && username) {
            this.token = token;
            this.username = username;
            this.elements.usernameDisplay.textContent = username;
            this.elements.logoutBtn.style.display = 'inline-block';
            this.hideAuthModal();
            this.connect();
        } else {
            this.showAuthModal();
        }
    }

    showAuthModal() {
        this.elements.authModal.classList.add('show');
        this.showLoginForm();
    }

    hideAuthModal() {
        this.elements.authModal.classList.remove('show');
    }

    showLoginForm() {
        this.elements.loginForm.style.display = 'block';
        this.elements.registerForm.style.display = 'none';
        this.elements.authTitle.textContent = 'Login';
        this.elements.authError.textContent = '';
        this.elements.loginUsername.focus();
    }

    showRegisterForm() {
        this.elements.loginForm.style.display = 'none';
        this.elements.registerForm.style.display = 'block';
        this.elements.authTitle.textContent = 'Register';
        this.elements.authError.textContent = '';
        this.elements.registerUsername.focus();
    }

    async login() {
        const username = this.elements.loginUsername.value.trim();
        const password = this.elements.loginPassword.value;

        if (!username || !password) {
            this.showError('Please enter username and password');
            return;
        }

        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });

            if (!response.ok) {
                const error = await response.text();
                throw new Error(error || 'Login failed');
            }

            const data = await response.json();
            this.handleAuthSuccess(data);
        } catch (error) {
            this.showError(error.message);
        }
    }

    async register() {
        const username = this.elements.registerUsername.value.trim();
        const email = this.elements.registerEmail.value.trim();
        const password = this.elements.registerPassword.value;

        if (!username || !email || !password) {
            this.showError('Please fill in all fields');
            return;
        }

        if (password.length < 6) {
            this.showError('Password must be at least 6 characters');
            return;
        }

        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, email, password })
            });

            if (!response.ok) {
                const error = await response.text();
                throw new Error(error || 'Registration failed');
            }

            const data = await response.json();
            this.handleAuthSuccess(data);
        } catch (error) {
            this.showError(error.message);
        }
    }

    handleAuthSuccess(data) {
        this.token = data.token;
        this.username = data.username;
        
        localStorage.setItem('token', data.token);
        localStorage.setItem('username', data.username);
        
        this.elements.usernameDisplay.textContent = data.username;
        this.elements.logoutBtn.style.display = 'inline-block';
        this.hideAuthModal();
        this.connect();
    }

    logout() {
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        
        this.disconnect();
        this.token = '';
        this.username = '';
        
        this.elements.usernameDisplay.textContent = '';
        this.elements.logoutBtn.style.display = 'none';
        this.elements.messages.innerHTML = '';
        
        this.showAuthModal();
    }

    showError(message) {
        this.elements.authError.textContent = message;
    }

    connect() {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            return;
        }

        this.showConnectionStatus('connecting', 'Connecting...');
        
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws?token=${encodeURIComponent(this.token)}`;
        
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
            
            setTimeout(() => this.hideConnectionStatus(), 2000);
        };

        this.ws.onmessage = (event) => {
            try {
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
            messageDiv.innerHTML = `
                <div class="message-content">${this.escapeHtml(message.content)}</div>
                <div class="timestamp">${this.formatTimestamp(message.timestamp)}</div>
            `;
        }

        this.elements.messages.appendChild(messageDiv);
        this.scrollToBottom();
        
        this.updateUserCount(message);
    }

    updateUserCount(message) {
        if (message.type === 'user_join' || message.type === 'user_left') {
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

document.addEventListener('DOMContentLoaded', () => {
    new ChatApp();
});