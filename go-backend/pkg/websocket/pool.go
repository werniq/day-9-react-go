package websocket

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true
			fmt.Println("size of connection pool", len(p.Clients))
			for client, _ := range p.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined"})
			}
			break
		case client := <-p.Unregister:
			delete(p.Clients, client)
			fmt.Println("size of connection pool", len(p.Clients))
			for client, _ := range p.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})

			}
			break
		case messsage := <-p.Broadcast:
			fmt.Println("Sending message to all clients in the pool")
			for client, _ := range p.Clients {
				if err := client.Conn.WriteJSON(messsage); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
