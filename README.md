# WebSocket Chat Application

A real-time chat application built with Go backend and modern JavaScript frontend, featuring beautiful UI design and seamless WebSocket communication.

## 🚀 Features

### 💬 **Real-time Messaging**
- Instant message delivery using WebSocket connections
- Support for multiple concurrent users
- System notifications for user join/leave events
- Name change notifications

### 🎨 **Modern UI Design**
- Beautiful gradient backgrounds with purple/blue color scheme
- Smooth animations and transitions
- Clean, minimalist interface with message bubbles
- Responsive design that works on all devices

### 👤 **User Management**
- Dynamic username setting with validation
- Real-time name changes with notifications
- Guest user fallback system
- Current user display in header

### 🔧 **Technical Features**
- WebSocket-based real-time communication
- Go backend with Gorilla WebSocket
- Clean separation of concerns with multiple Go files
- Static file serving
- Error handling and connection management

## 📁 Project Structure

```
websocket-chat/
├── go.mod                  # Go module file
├── main.go                 # Main server file
├── client.go               # WebSocket client handling
├── manager.go              # Client connection management
├── message.go              # Message structure definitions
└── static/
    ├── index.html          # Main HTML interface
    ├── style.css           # Modern CSS styling
    └── script.js           # Frontend JavaScript logic
```

## 🛠️ Installation & Setup

### Prerequisites
- Go 1.21 or higher
- Modern web browser
- VS Code (recommended)

### Step 1: Clone and Initialize
```bash
mkdir websocket-chat
cd websocket-chat
go mod init websocket-chat
```

### Step 2: Install Dependencies
```bash
go get github.com/gorilla/websocket
go get github.com/satori/go.uuid
```

### Step 3: Create Project Files
Copy all the provided files into their respective locations according to the project structure above.

### Step 4: Run the Application
```bash
go run .
```

### Step 5: Access the Chat
Open your browser and navigate to:
```
http://localhost:8080
```

## 🎯 Usage

1. **Join the Chat**: Enter your name when prompted
2. **Send Messages**: Type in the input field and press Enter or click Send
3. **Change Name**: Click the "Change Name" button in the header
4. **Multiple Users**: Open multiple browser tabs to test with different users

## 🏗️ Architecture

### Backend (Go)
- **main.go**: HTTP server setup and WebSocket upgrade handling
- **client.go**: Individual client connection management and message handling
- **manager.go**: Centralized client management and message broadcasting
- **message.go**: Message structure definitions for JSON communication

### Frontend (JavaScript)
- **Modular Design**: Clean separation of UI and WebSocket logic
- **Event-Driven**: Responsive to user interactions and server messages
- **Error Handling**: Graceful handling of connection issues

### Message Types
- **setName**: User name setting/changing
- **message**: Regular chat messages
- **system**: Server notifications (join/leave/name changes)

## 🎨 UI Components

### Chat Interface
- **Header**: App title, current username, and change name button
- **Message Area**: Scrollable message history with different message types
- **Input Area**: Text input field with send functionality

### Modal System
- **Name Input Modal**: Clean popup for username entry
- **Form Validation**: Name length and content validation
- **Smooth Animations**: Fade and slide effects

### Message Styling
- **User Messages**: Blue gradient bubbles with sender names
- **System Messages**: Orange-accented notifications
- **Responsive Layout**: Adapts to different screen sizes

## 🔧 Development

### VS Code Setup
Recommended extensions:
- Go (by Google)
- Go Outliner
- GitLens
- Thunder Client

### Debugging Configuration
Create `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
            "console": "integratedTerminal"
        }
    ]
}
```

## 🧪 Testing

1. **Single User**: Test basic functionality with one browser window
2. **Multiple Users**: Open several tabs to test real-time communication
3. **Name Changes**: Test name changing functionality
4. **Connection Issues**: Test reconnection handling

## 🐛 Troubleshooting

### Common Issues

**Port Already in Use:**
```bash
lsof -ti:8080 | xargs kill -9
```

**Module Dependencies:**
```bash
go mod tidy
```

**WebSocket Connection Issues:**
- Check firewall settings
- Ensure port 8080 is available
- Verify WebSocket URL protocol (ws:// vs wss://)

**CORS Issues:**
The application includes CORS handling in the WebSocket upgrader.

## 🚀 Deployment

### Local Development
```bash
go run .
```

### Production Build
```bash
go build -o chat-server
./chat-server
```

## 🔐 Security Considerations

- Input validation for usernames and messages
- WebSocket origin checking
- Rate limiting (can be added)
- Message sanitization (recommended for production)

## 📈 Performance

- Efficient message broadcasting to multiple clients
- Memory management with proper connection cleanup
- Scalable architecture with goroutines

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

**Happy Chatting! 💬**
