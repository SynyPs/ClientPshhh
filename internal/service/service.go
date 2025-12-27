package service

import (
	"context"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type ChatClient struct {
	conn        *websocket.Conn
	IncomingMsg chan []byte
}

func NewChatClient(url string) (*ChatClient, error) {
	ctx := context.Background()
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	client := &ChatClient{
		conn:        c,
		IncomingMsg: make(chan []byte),
	}

	go client.listen()

	return client, nil
}

func (c *ChatClient) listen() {
	defer c.conn.Close(websocket.StatusInternalError, "Internal error")

	for {
		var v interface{}
		err := wsjson.Read(context.Background(), c.conn, &v)
		if err != nil {
			close(c.IncomingMsg)
			return
		}
	}
}
