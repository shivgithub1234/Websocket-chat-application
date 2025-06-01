class ChatClient {
    constructor() {
        this.socket = null;
        this.messages = document.getElementById('messages');
        this.messageInput = document.getElementById('messageInput');
        this.connect();
        this.setupEventListeners();
    }

    connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        
        this.socket = new WebSocket(wsUrl);
        
        this.socket.onopen = (event) => {
            this.addMessage('Connected to chat server', true);
        };
        
        this.socket.onclose = (event) => {
            this.addMessage('Disconnected from chat server', true);
        };
        
        this.socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            let content = message.content;
            
            if (message.sender) {
                content = `${message.sender}: ${content}`;
            }
            
            this.addMessage(content, content.startsWith('A new socket') || content.startsWith('A socket has'));
        };
        
        this.socket.onerror = (error) => {
            this.addMessage('Connection error occurred', true);
        };
    }

    addMessage(content, isSystem = false) {
        const messageDiv = document.createElement('div');
        messageDiv.className = isSystem ? 'message system-message' : 'message';
        messageDiv.textContent = content;
        
        this.messages.appendChild(messageDiv);
        this.messages.scrollTop = this.messages.scrollHeight;
    }

    sendMessage() {
        const message = this.messageInput.value.trim();
        if (message && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(message);
            this.messageInput.value = '';
        }
    }

    setupEventListeners() {
        this.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.sendMessage();
            }
        });
    }
}

// Initialize chat client when page loads
let chatClient;
window.addEventListener('DOMContentLoaded', () => {
    chatClient = new ChatClient();
});

// Global function for button onclick
function sendMessage() {
    chatClient.sendMessage();
}