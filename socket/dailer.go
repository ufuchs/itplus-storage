package socket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"hidrive.com/ufuchs/itplus/base/fcc"

	"github.com/gorilla/websocket"
)

//
// dialSocket
//
func dialSocket(ctx context.Context, URL string) (conn *websocket.Conn, err error) {

	var wg *sync.WaitGroup

	if wg, err = fcc.GetWGFromContext(ctx); err != nil {
		return nil, err
	}

	wg.Add(1)
	defer wg.Done()

	r := connInstance{}
	ch := make(chan connInstance, 1)

	var dial = func(URL string) <-chan connInstance {
		fmt.Println("==> Websocket, trying URL ", URL)
		r.conn, _, r.err = websocket.DefaultDialer.Dial(URL, nil)
		ch <- r
		return ch
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case ok := <-dial(URL):
			if ok.err == nil {
				fmt.Printf("==> Websocket, got connection to '%v'\n", URL)
				return ok.conn, ok.err
			}
			fmt.Printf("==> Websocket, got error: '%v'\n", ok.err)
			time.Sleep(500 * time.Millisecond)
		}
	}

}
