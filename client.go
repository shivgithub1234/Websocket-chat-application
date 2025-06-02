package main

import (
    "github.com/gorilla/websocket"
    "encoding/json"
    "log"
    "time"
)

// Client represents a connected websocket client
type Client struct {
    id     string
    name   string
    socket *websocket.Conn
    send   chan []byte
}

// read handles reading messages from the websocket connection
func (c *Client) read() {
    defer func() {
        log.Printf("Client %s (%s) read goroutine ending", c.id, c.name)
        manager.unregister <- c
        c.socket.Close()
    }()

    // Set read deadline and pong handler for keepalive
    c.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.socket.SetPongHandler(func(string) error {
        c.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    for {
        _, message, err := c.socket.ReadMessage()
        if err != nil {
            log.Printf("Error reading message from client %s: %v", c.id, err)
            break
        }
        
        log.Printf("Received message from client %s (%s): %s", c.id, c.name, string(message))
        
        // Parse incoming message to check if it's a name setting message
        var incomingMsg Message
        if err := json.Unmarshal(message, &incomingMsg); err == nil {
            if incomingMsg.Type == "setName" {
                log.Printf("Client %s attempting to set name to: %s", c.id, incomingMsg.Content)
                
                // Check if name is already taken
                if manager.isNameTaken(incomingMsg.Content) {
                    log.Printf("Name '%s' is already taken", incomingMsg.Content)
                    // Send error message back to client
                    errorMsg, _ := json.Marshal(&Message{
                        Type: "nameError",
                        Content: "Name already taken. Please choose another name.",
                    })
                    c.send <- errorMsg
                    continue
                }
                
                // Update client name
                oldName := c.name
                c.name = incomingMsg.Content
                log.Printf("Client %s name set to: %s", c.id, c.name)
                
                // Send success message back to client
                successMsg, _ := json.Marshal(&Message{
                    Type: "nameSuccess",
                    Content: "Name set successfully",
                })
                c.send <- successMsg
                
                // Notify others about name change/join
                var notificationMsg string
                if oldName == "" {
                    notificationMsg = c.name + " joined the chat"
                } else {
                    notificationMsg = oldName + " changed name to " + c.name
                }
                
                jsonMessage, _ := json.Marshal(&Message{
                    Type: "system",
                    Content: notificationMsg,
                })
                manager.broadcast <- jsonMessage
                continue
            }
        }
        
        // Regular chat message
        if c.name == "" {
            log.Printf("Client %s tried to send message without setting name", c.id)
            continue
        }
        
        jsonMessage, _ := json.Marshal(&Message{
            Type:    "message",
            Sender:  c.name,
            Content: string(message),
        })
        manager.broadcastToOthers <- &BroadcastMessage{
            message: jsonMessage,
            sender:  c,
        }
        
        // Send back to sender with isOwnMessage flag
        ownMessage, _ := json.Marshal(&Message{
            Type:    "message",
            Sender:  c.name,
            Content: string(message),
            IsOwnMessage: true,
        })
        c.send <- ownMessage
    }
}

// write handles writing messages to the websocket connection
func (c *Client) write() {
    ticker := time.NewTicker(54 * time.Second) // Ping every 54 seconds
    defer func() {
        ticker.Stop()
        c.socket.Close()
        log.Printf("Client %s (%s) write goroutine ending", c.id, c.name)
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                c.socket.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            if err := c.socket.WriteMessage(websocket.TextMessage, message); err != nil {
                log.Printf("Error writing message to client %s: %v", c.id, err)
                return
            }
            
        case <-ticker.C:
            c.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
                log.Printf("Error sending ping to client %s: %v", c.id, err)
                return
            }
        }
    }
}