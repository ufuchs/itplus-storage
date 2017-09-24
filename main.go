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
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-sql-driver/mysql"
	"github.com/ufuchs/itplus/base/fcc"
	"github.com/ufuchs/itplus/base/zvous"
	"github.com/ufuchs/itplus/storage/app"
	"github.com/ufuchs/itplus/storage/database"
	"github.com/ufuchs/itplus/storage/socket"

	"github.com/ufuchs/zeroconf"
	//_ "net/http/pprof"
)

//
//
//
func prepare() {

	var err error

	if app.BaseDir, err = os.Getwd(); err != nil {
		fcc.Fatal(err)
	}

	svc := app.NewConfigService().
		RetrieveAll()

	if svc.LastErr != nil {
		fcc.Fatal(svc.LastErr)
	}

}

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
		err error
		dbs *database.Service
		// discoverer  *discovery.TCPDiscoverer
		// discovered  discovery.ServiceEntryler
		sigs = make(chan os.Signal, 2)
		wg   sync.WaitGroup
	)

	prepare()

	ctx, cancel := context.WithCancel(context.Background())
	ctxWg := context.WithValue(ctx, 0, &wg)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
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

	if dbs, err = database.NewService(ctx, app.DSN, "salata"); err != nil {
		fcc.Fatal(err)
	}

	defer dbs.Close()

Retry:

	// for err = dbs.Prepare()

	if err = dbs.Prepare(); err != nil {

		me, ok := err.(*mysql.MySQLError)
		if !ok {
			fcc.Fatal(err)
		}

		fmt.Println(me.Number)

		if me.Number != 1146 {
			fcc.Fatal(err)
		}

		// table doesn't exist
		if err = dbs.CreateTableGateway(); err != nil {
			fcc.Fatal(err)
		}

		goto Retry

	}

	client := socket.NewClient(ctxWg, app.Hub, 1)

	dbs.In = client.Out
	go dbs.Run(ctxWg)

	select {
	case <-ctx.Done():
		wg.Wait()
		fmt.Println("==> Service stopped...")
	}

}
