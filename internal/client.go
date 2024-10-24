package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	conn   *websocket.Conn
	userId int64
	hub    *Hub
	sendTo chan *Message
}

type Message struct {
	From    int64  `json:"from"`
	To      int64  `json:"to"`
	Content string `json:"Content"`
}

func (c *Client) SendError(errorMsg string) {

	c.sendTo <- &Message{
		From:    0, //from server
		To:      c.userId,
		Content: fmt.Sprintf(`{"error": true, "message": "%s"}`, errorMsg),
	}
	// Add a sleep time to allow the channel to be emptied
	time.Sleep(2 * time.Second)
	c.conn.Close()

}

func (c *Client) readPump() {

	defer func() {
		c.conn.Close()
		c.hub.unregister <- c
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		log.Println(msg)

		//Here Client.Id comes from the token and msg.from comes from the payload both of them should be same so that we can say we are sending the message by a credible user
		if c.userId != msg.From {
			c.SendError("Unauthorized !! Auth user and Sent User are Different")
			return
		}
		c.hub.broadcast <- &msg
	}
}
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.sendTo:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Marshal the message to JSON
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(message)
			finalPayload := reqBodyBytes.Bytes()

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// Write the initial message
			w.Write(finalPayload)

			// Write any remaining queued messages
			n := len(c.sendTo)
			for i := 0; i < n; i++ {
				json.NewEncoder(reqBodyBytes).Encode(<-c.sendTo)
				w.Write(reqBodyBytes.Bytes())
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

		//send to the ssocket listener
	}
}

func CreateNewSocketUser(hub *Hub, connection *websocket.Conn, userID int64) {

	client := &Client{
		hub:    hub,
		conn:   connection,
		userId: userID,
		sendTo: make(chan *Message),
	}

	go client.readPump()
	go client.writePump()

	client.hub.register <- client
}
