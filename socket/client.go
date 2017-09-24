package socket

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/ufuchs/itplus/base/fcc"

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

const (
	BUFSIZE = 8196
)

type (
	Client struct {
		ID     int
		URL    string
		Conn   *websocket.Conn
		ctx    context.Context
		In     chan []byte
		Out    chan []byte
		closed bool
	}

	connInstance struct {
		err  error
		conn *websocket.Conn
	}
)

//
// NewClient
//
func NewClient(ctx context.Context, addr string, ID int) *Client {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	c := &Client{
		ID:  ID,
		URL: u.String(),
		Out: make(chan []byte, 512),
	}

	go c.Run(ctx)

	return c

}

//
//
//
func (c *Client) readPump() {
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fmt.Println("==> Server closed connection:", c.Conn.RemoteAddr())
			} else {
				fmt.Println("==> Server socket err: ", err)
			}
			c.closed = true
			break
		}
		//fmt.Println(string(p))
		c.Out <- append([]byte(nil), p...)
	}
}

//
//
//
func (c *Client) writePump(ctx context.Context) {

	ticker := time.NewTicker(pingPeriod)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
		case message, ok := <-c.In:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.closed = true
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.closed = true
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				c.closed = true
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.closed = true
				return
			}
		}
	}
}

//
//
//
func (c *Client) writeMeasurement(ctx context.Context, conn *websocket.Conn) error {

	defer c.Conn.Close()

	go c.writePump(ctx)

	go c.readPump()

	select {
	case <-ctx.Done():

		// To cleanly close a connection, a client should send a close
		// frame and wait for the server to close the connection.
		if !c.closed {
			err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Storage server shutting down"))
			if err != nil {
				fmt.Println("write close:", err)
			}
		}
		fmt.Println("==> Closing web socket")
		return ctx.Err()
	}

}

//
// http://dahernan.github.io/2015/02/04/context-and-cancellation-of-goroutines/
//
func (c *Client) Run(ctx context.Context) (err error) {

	var wg *sync.WaitGroup

	var random = func(min, max int) time.Duration {
		rand.Seed(time.Now().Unix())
		return time.Duration(rand.Intn(max-min) + min)
	}

	if wg, err = fcc.GetWGFromContext(ctx); err != nil {
		return err
	}

	wg.Add(1)
	defer wg.Done()

	for {

		c.Conn, err = dialSocket(ctx, c.URL)

		if err != nil {
			if err.Error() == "context canceled" {
				return err
			}
		}

		err = c.writeMeasurement(ctx, c.Conn)
		if err.Error() == "context canceled" {
			return err
		}

		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			fmt.Println("----------------------------------")
			// close() ??
			return err
		}

		// https://tools.ietf.org/html/rfc6455#section-7.2.3
		time.Sleep(random(0, 5) * time.Second)

	}

}
