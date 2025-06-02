class ChatClient {
    constructor() {
        this.socket = null;
        this.userName = '';
        this.messages = document.getElementById('messages');
        this.messageInput = document.getElementById('messageInput');
        this.nameModal = document.getElementById('nameModal');
        this.chatContainer = document.getElementById('chatContainer');
        this.currentNameSpan = document.getElementById('currentName');
        this.nameError = document.getElementById('nameError');
        this.joinBtn = document.getElementById('joinBtn');
        this.sendBtn = document.getElementById('sendBtn');
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        
        this.showNameModal();
        this.setupEventListeners();
    }

    showNameModal() {
        this.nameModal.style.display = 'flex';
        this.chatContainer.style.display = 'none';
        const nameInput = document.getElementById('nameInput');
        nameInput.value = '';
        nameInput.focus();
        this.nameError.textContent = '';
        this.joinBtn.disabled = false;
        this.joinBtn.textContent = 'Join Chat';
        
        // Handle Enter key in name input
        nameInput.onkeypress = (e) => {
            if (e.key === 'Enter' && !this.joinBtn.disabled) {
                this.setName();
            }
        };
    }

    setName() {
        const nameInput = document.getElementById('nameInput');
        const name = nameInput.value.trim();
        
        if (name === '') {
            this.showNameError('Please enter a valid name');
            return;
        }
        
        if (name.length > 20) {
            this.showNameError('Name must be 20 characters or less');
            return;
        }
        
        if (name.length < 2) {
            this.showNameError('Name must be at least 2 characters');
            return;
        }
        
        // Disable button and show loading
        this.joinBtn.disabled = true;
        this.joinBtn.textContent = 'Connecting...';
        this.nameError.textContent = '';
        
        this.userName = name;
        
        // Connect to WebSocket and set name
        this.connect();
    }

    showNameError(message) {
        this.nameError.textContent = message;
        this.nameError.style.color = '#e74c3c';
        
        // Reset button state
        this.joinBtn.disabled = false;
        this.joinBtn.textContent = 'Join Chat';
    }

    connect() {
        console.log('Attempting to connect to WebSocket...');
        
        // Close existing connection if any
        if (this.socket) {
            this.socket.close();
        }
        
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        
        console.log('Connecting to:', wsUrl);
        
        try {
            this.socket = new WebSocket(wsUrl);
            
            this.socket.onopen = (event) => {
                console.log('WebSocket connection opened');
                this.reconnectAttempts = 0;
                this.addMessage('Connected to chat server', 'system');
                this.sendNameToServer();
            };
            
            this.socket.onclose = (event) => {
                console.log('WebSocket connection closed:', event.code, event.reason);
                this.addMessage('Disconnected from chat server', 'system');
                this.sendBtn.disabled = true;
                this.messageInput.disabled = true;
                this.messageInput.placeholder = "Disconnected...";
                
                // Attempt to reconnect if not intentional disconnect
                if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
                    setTimeout(() => {
                        this.reconnectAttempts++;
                        console.log(`Reconnection attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}`);
                        this.addMessage(`Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`, 'system');
                        this.connect();
                    }, 2000 * this.reconnectAttempts);
                }
            };
            
            this.socket.onmessage = (event) => {
                console.log('Received message:', event.data);
                
                let message;
                try {
                    message = JSON.parse(event.data);
                } catch (e) {
                    console.error('Error parsing message:', e);
                    return;
                }
                
                if (message.type === 'nameError') {
                    this.showNameError(message.content);
                } else if (message.type === 'nameSuccess') {
                    console.log('Name set successfully');
                    // Hide modal and show chat
                    this.nameModal.style.display = 'none';
                    this.chatContainer.style.display = 'flex';
                    this.currentNameSpan.textContent = this.userName;
                    
                    // Enable input and button
                    this.messageInput.disabled = false;
                    this.sendBtn.disabled = false;
                    this.messageInput.placeholder = "Type your message...";
                    
                    // Focus on input
                    setTimeout(() => {
                        this.messageInput.focus();
                    }, 100);
                } else if (message.type === 'system') {
                    this.addMessage(message.content, 'system');
                } else if (message.type === 'message') {
                    if (message.isOwnMessage) {
                        this.addMessage(message.content, 'own', message.sender);
                    } else {
                        this.addMessage(message.content, 'other', message.sender);
                    }
                }
            };
            
            this.socket.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.addMessage('Connection error occurred', 'system');
                this.showNameError('Connection failed. Please try again.');
            };
            
        } catch (error) {
            console.error('Failed to create WebSocket connection:', error);
            this.showNameError('Failed to connect. Please try again.');
        }
    }

    sendNameToServer() {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            const nameMessage = {
                type: 'setName',
                content: this.userName
            };
            console.log('Sending name to server:', nameMessage);
            this.socket.send(JSON.stringify(nameMessage));
        } else {
            console.error('WebSocket not ready for sending name');
        }
    }

    addMessage(content, type = 'other', sender = '') {
        const messageDiv = document.createElement('div');
        
        if (type === 'system') {
            messageDiv.className = 'message system-message';
            messageDiv.innerHTML = `<span class="content">${this.escapeHtml(content)}</span>`;
        } else {
            messageDiv.className = `message ${type}`;
            
            if (sender) {
                messageDiv.innerHTML = `
                    <span class="sender">${this.escapeHtml(sender)}</span>
                    <span class="content">${this.escapeHtml(content)}</span>
                `;
            } else {
                messageDiv.innerHTML = `<span class="content">${this.escapeHtml(content)}</span>`;
            }
        }
        
        this.messages.appendChild(messageDiv);
        this.messages.scrollTop = this.messages.scrollHeight;
    }
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    sendMessage() {
        const message = this.messageInput.value.trim();
        if (message && this.socket && this.socket.readyState === WebSocket.OPEN) {
            console.log('Sending message:', message);
            this.socket.send(message);
            this.messageInput.value = '';
        } else if (this.socket && this.socket.readyState !== WebSocket.OPEN) {
            this.addMessage('Connection lost. Message not sent.', 'system');
        }
    }

    setupEventListeners() {
        this.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !this.sendBtn.disabled && !this.messageInput.disabled) {
                this.sendMessage();
            }
        });
        
        // Initially disable send button and input
        this.sendBtn.disabled = true;
        this.messageInput.disabled = true;
        this.messageInput.placeholder = "Join chat first...";
    }
}

// Initialize chat client when page loads
let chatClient;
window.addEventListener('DOMContentLoaded', () => {
    console.log('DOM loaded, initializing chat client...');
    chatClient = new ChatClient();
});

// Global functions for button onclick
function sendMessage() {
    if (chatClient) {
        chatClient.sendMessage();
    }
}

function setName() {
    if (chatClient) {
        chatClient.setName();
    }
}

function showNameModal() {
    if (chatClient) {
        chatClient.showNameModal();
    }
}