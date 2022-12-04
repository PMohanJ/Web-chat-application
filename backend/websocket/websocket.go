package websocket

import (
	"log"

	"github.com/gin-gonic/gin"
)

type WebSockets struct {
	Clients   map[string][]*Client
	Broadcast chan map[string]interface{}
}

func CreateWebSocketsServer() *WebSockets {
	return &WebSockets{
		Clients:   make(map[string][]*Client),
		Broadcast: make(chan map[string]interface{}),
	}
}

func (ws *WebSockets) WSEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := Upgrade(c.Writer, c.Request)
		if err != nil {
			log.Panic(err)
		}

		client := &Client{
			Conn:       conn,
			WebSockets: ws,
		}

		for {
			var data map[string]interface{}
			err := client.Conn.ReadJSON(&data)
			if err != nil {
				log.Panic(err)
			}

			log.Printf("data received is %+v", data)
			err = ws.HandleClientMessage(client, data)
			if err != nil {
				log.Panic(err)
			}
		}

	}
}

// HandleClientMessage adds client to the respective chats, and also sends client messages
// to websocket broadcast channel
func (ws *WebSockets) HandleClientMessage(clientObj *Client, data map[string]interface{}) error {

	// check if client is initiating the connection
	if data["messageType"] == "setup" {
		// create that chat and add the client
		chatId := data["chat"].(string)

		var clients []*Client
		clients = append(clients, clientObj)
		ws.Clients[chatId] = clients

		log.Println("Client added to list")
	} else {
		// check if chat already exist in our pool
		chatId := data["chat"].(string)
		cl, exists := ws.Clients[chatId]
		if exists {
			var clientExists bool
			for _, client := range cl {
				if clientObj == client {
					clientExists = true

					// client already exists, so send the data
					clientObj.WebSockets.Broadcast <- data
					break
				}
			}

			// if client doesn't exist in that chat so far, add that client and broadcast the msg
			if !clientExists {
				clients := ws.Clients[chatId]
				clients = append(clients, clientObj)
				ws.Clients[chatId] = clients
				log.Println("Client added to list")
				clientObj.WebSockets.Broadcast <- data
			}
		} else {
			// create that chat and add the client
			var clients []*Client
			clients = append(clients, clientObj)
			ws.Clients[chatId] = clients

			log.Println("Client added to list")
			clientObj.WebSockets.Broadcast <- data
		}
	}

	return nil
}

// SendMessage receives messages from broadcast channel and sends them to
// respective chat members aka clients
func (ws *WebSockets) SendMessage() {
	for {
		msg := <-ws.Broadcast
		chatId := msg["chat"].(string)
		log.Printf("msg received from broad %+v", msg)

		// get the chat to which the msg should be send
		clientsOfThisChat := ws.Clients[chatId]
		log.Printf("size of chat %v and chatid is: %v", len(clientsOfThisChat), chatId)
		for _, client := range clientsOfThisChat {
			err := client.Conn.WriteJSON(msg)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}
