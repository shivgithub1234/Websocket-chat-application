package main
import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

var manager = ClientManager{
	broadcast : make(chan []byte),
	unregister : make(chan *Client),
	register : make(chan *Client),
	clients: make(map[*Client]bool),
}

func wspage(res http.ResponseWriter,req *http.Request){
	conn,err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool{
			return true
		},
	
	}).Upgrade(res,req,nil)
	if err!= nil{
		http.NotFound(res,req)
		return
	}

	// create new client
	client := &Client{
		id: uuid.NewV4().String(),
		socket: conn,
		send: make(chan []byte),
	}
	manager.register <- client // register the client
	go client.read()
	go client.write()
}

func main(){
	fmt.Println("Starting Websocket chat server...")
	go manager.start() // start the client manager
	http.Handle("/",http.FileServer(http.Dir("./static"))) // serve static files
	http.HandleFunc("/ws",wspage) // handle websocket connections
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080",nil)) // start the server
	
}
