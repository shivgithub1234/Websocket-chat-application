// package main
// import (
// 	"encoding/json"
// )

// type ClientManager struct {
// 	clients map[*Client]bool
// 	broadcast chan []byte
// 	unregister chan *Client
// 	register chan *Client
// }

// func (manager *ClientManager)start(){
// 	for{
// 		select {
// 		case conn := <- manager.register:
// 			manager.clients[conn] = true
// 			jsonMessage,_ := json.Marshal(&Message{
// 				Content: "A new user has joined the chat",
// 			})
// 			manager.send(jsonMessage, conn)
// 		case conn := <-manager.unregister:
// 			if _, ok := manager.clients[conn]; ok {
// 				delete(manager.clients, conn)
// 				close(conn.send)
// 				jsonMessage, _ := json.Marshal(&Message{
// 					Content: "A user has left the chat",
// 				})
// 				manager.send(jsonMessage, conn)
// 			}
// 		case message := <-manager.broadcast:
// 			for conn := range manager.clients {
// 				select {
// 				case conn.send <- message:
// 				default:
// 					close(conn.send)
// 					delete(manager.clients, conn)
// 				}
// 			}
// 		}
// 	}
// }

// func (manager *ClientManager) send(message []byte,ignore *Client){
// 	for conn := range manager.clients{
// 		if conn != ignore{
// 			conn.send <- message
// 		}
// 	}
// }


package main

import (
    "encoding/json"
    "strings"
)

// BroadcastMessage represents a message to broadcast with sender info
type BroadcastMessage struct {
    message []byte
    sender  *Client
}

// ClientManager manages all connected clients
type ClientManager struct {
    clients          map[*Client]bool
    broadcast        chan []byte
    broadcastToOthers chan *BroadcastMessage
    register         chan *Client
    unregister       chan *Client
}

// isNameTaken checks if a name is already in use
func (manager *ClientManager) isNameTaken(name string) bool {
    normalizedName := strings.ToLower(strings.TrimSpace(name))
    for client := range manager.clients {
        if strings.ToLower(strings.TrimSpace(client.name)) == normalizedName && client.name != "" {
            return true
        }
    }
    return false
}

// start runs the main client manager loop
func (manager *ClientManager) start() {
    for {
        select {
        case conn := <-manager.register:
            manager.clients[conn] = true
            // Don't send connection message immediately - wait for name to be set

        case conn := <-manager.unregister:
            if _, ok := manager.clients[conn]; ok {
                close(conn.send)
                delete(manager.clients, conn)
                
                // Send disconnection message with name if available
                if conn.name != "" {
                    jsonMessage, _ := json.Marshal(&Message{
                        Type:    "system",
                        Content: conn.name + " left the chat",
                    })
                    manager.send(jsonMessage, conn)
                }
            }

        case message := <-manager.broadcast:
            for conn := range manager.clients {
                select {
                case conn.send <- message:
                default:
                    close(conn.send)
                    delete(manager.clients, conn)
                }
            }

        case broadcastMsg := <-manager.broadcastToOthers:
            for conn := range manager.clients {
                if conn != broadcastMsg.sender {
                    select {
                    case conn.send <- broadcastMsg.message:
                    default:
                        close(conn.send)
                        delete(manager.clients, conn)
                    }
                }
            }
        }
    }
}

// send broadcasts a message to all clients except the ignored one
func (manager *ClientManager) send(message []byte, ignore *Client) {
    for conn := range manager.clients {
        if conn != ignore {
            conn.send <- message
        }
    }
}