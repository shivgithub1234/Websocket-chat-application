package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/satori/go.uuid"
)

var manager = ClientManager{
    broadcast:         make(chan []byte),
    broadcastToOthers: make(chan *BroadcastMessage),
    register:          make(chan *Client),
    unregister:        make(chan *Client),
    clients:           make(map[*Client]bool),
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Allow connections from any origin
        return true
    },
}

func wsPage(res http.ResponseWriter, req *http.Request) {
    // Add debug logging
    log.Printf("WebSocket connection attempt from: %s", req.RemoteAddr)
    
    // Upgrade HTTP connection to WebSocket
    conn, err := upgrader.Upgrade(res, req, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }
    
    log.Printf("WebSocket connection established with: %s", req.RemoteAddr)

    // Create new client
    client := &Client{
        id:     uuid.NewV4().String(),
        name:   "",
        socket: conn,
        send:   make(chan []byte, 256), // Add buffer to prevent blocking
    }

    // Register client
    manager.register <- client

    // Start client goroutines
    go client.read()
    go client.write()
}

func main() {
    fmt.Println("Starting Enhanced WebSocket chat server...")
    
    // Start the client manager
    go manager.start()

    // Serve static files
    fs := http.FileServer(http.Dir("./static/"))
    http.Handle("/", fs)
    
    // WebSocket endpoint
    http.HandleFunc("/ws", wsPage)

    // Start server
    fmt.Println("Server running on http://localhost:8080")
    fmt.Println("Press Ctrl+C to stop the server")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}