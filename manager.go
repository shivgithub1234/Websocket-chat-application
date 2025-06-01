package main
import (
	"encoding/json"
)

type ClientManager struct {
	clients map[*Client]bool
	broadcast chan []byte
	unregister chan *Client
	register chan *Client
}

func (manager *ClientManager)start(){
	for{
		select {
		case conn := <- manager.register:
			manager.clients[conn] = true
			jsonMessage,_ := json.Marshal(&Message{
				Content: "A new user has joined the chat",
			})
			manager.send(jsonMessage, conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				delete(manager.clients, conn)
				close(conn.send)
				jsonMessage, _ := json.Marshal(&Message{
					Content: "A user has left the chat",
				})
				manager.send(jsonMessage, conn)
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
		}
	}
}

func (manager *ClientManager) send(message []byte,ignore *Client){
	for conn := range manager.clients{
		if conn != ignore{
			conn.send <- message
		}
	}
}