package pubsub

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const ReadTimeout = 1 * time.Minute
const WriteTimeout = 1 * time.Minute

type Client struct {
	conn *websocket.Conn
	input chan []byte
	errChan chan error
}


func (c *Client) Writeloop(){
	resetDeadline := func() error {
		dl := time.Now().Add(WriteTimeout)
		if err := c.conn.SetWriteDeadline(dl); err != nil {
			return err
		}
		return nil
	}

	for b := range c.input {
		resetDeadline()
		c.conn.WriteMessage(websocket.BinaryMessage, b)
	}
}

func (c *Client) Readloop()chan []byte{
	resetDeadline := func() error {
		dl := time.Now().Add(ReadTimeout)
		if err := c.conn.SetReadDeadline(dl); err != nil {
			return err
		}
		return nil
	}

	output := make(chan []byte)

	go func(){
		defer close(output)
		for {
			resetDeadline()
			mtype, b, err := c.conn.ReadMessage()
			if err != nil {
				c.errChan <- err
				return
			}
			
			if mtype == websocket.CloseMessage {
				c.errChan <- err
				return
			}

			output <- b
		}
	}()

	return output
}

func (c *Client) Close() {
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
    if err != nil {
		//appear when client close 
		if !errors.Is(err, websocket.ErrCloseSent){
			slog.Error("write close", "err", err)
		}
    }	
	close(c.input)
	err = c.conn.Close()
    if err != nil {
    	slog.Error("conn close", "err", err)
    }

}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{ 
		conn: conn,
		input: make(chan []byte), 
		errChan: make(chan error, 1),
	}
}
