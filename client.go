package main
import (
	"github.com/gorilla/websocket"
	"encoding/json"
)

type Client struct {
	id string
	socket *websocket.Conn
	send chan[]byte
}

// reading function
func (c *Client) read() {
	defer func(){
		manager.unregister <- c
		c.socket.Close()
	}()

	for{
		_,message,err := c.socket.ReadMessage()
		if(err != nil) {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		// here we are trying to comnvert the meassege to proper format
		jsonMessage,_ := json.Marshal(&Message{
			Sender:c.id,
			Content:string(message),
		})
		manager.broadcast <- jsonMessage // it sends the message to all except the sender

	}

}

// wrting function
func (c *Client) write(){
	defer func(){
		c.socket.Close()
	}()
	for{
		select {
		case message,ok := <-c.send:
			if(!ok){
				c.socket.WriteMessage(websocket.CloseMessage,[]byte{})
				return 
			}
			c.socket.WriteMessage(websocket.TextMessage,message)
		}
	}
		
}