//
// Copyright(c) 2017 Uli Fuchs <ufuchs@gmx.com>
// MIT Licensed
//

// [ Zeit ist das, was man an der Uhr abliest              ]
// [                                   - Albert Einstein - ]

// see: https://github.com/fatih/vim-go-tutorial/blob/master/vimrc

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ufuchs/itplus/base/zvous"
	//	"github.com/ufuchs/itplus/storage/app"
	"github.com/ufuchs/itplus/storage/database"
	"github.com/ufuchs/itplus/storage/socket"

	"github.com/ufuchs/zeroconf"
	//_ "net/http/pprof"
)

const PORT = ":8080"

//
// https://stackoverflow.com/questions/30652577/go-doing-a-get-request-and-building-the-querystring
//
// func (r *ReconfigurableDispatcher) RequestConfiguration() {

// 	client := &http.Client{}

// 	req, _ := http.NewRequest("GET", "http://api.themoviedb.org/3/tv/popular", nil)
// 	req.Header.Add("Accept", "application/json")
// 	resp, err := client.Do(req)

// 	if err != nil {
// 		fmt.Println("Errored when sending request to the server")
// 		return
// 	}

// 	defer resp.Body.Close()
// 	resp_body, _ := ioutil.ReadAll(resp.Body)

// 	fmt.Println(resp.Status)
// 	fmt.Println(string(resp_body))

// }

//
//
//
func handleSignals(cancel context.CancelFunc, sigch <-chan os.Signal) {
	for {
		select {
		case <-sigch:
			cancel()
		}
	}
}

//
//
//
func main() {

	var (
		// err         error
		// discoverer  *discovery.TCPDiscoverer
		// discovered  discovery.ServiceEntryler
		sigs = make(chan os.Signal, 2)
		wg   sync.WaitGroup
	)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	ctxWg := context.WithValue(ctx, 0, &wg)
	go handleSignals(cancel, sigs)

	hubs := zvous.NewZCBrowserService(zvous.AVAHI_MEASUREMENT, zeroconf.IPv4, 4)

	go func() {

		for {
			select {
			case in := <-hubs.Out:
				for _, e := range in {
					fmt.Println(e)
				}
			}
		}

	}()

	client := socket.NewClient("192.168.178.12:8080", 1)

	go client.Run(ctxWg)

	db, err := database.NewService(ctx, database.DSN, "salata")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	db.In = client.Out

	go db.Run(ctxWg)

	select {
	case <-ctx.Done():
		wg.Wait()
		fmt.Println("==> Service stopped...")
	}

}
