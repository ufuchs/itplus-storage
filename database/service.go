package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"hidrive.com/ufuchs/itplus/base/fcc"

	_ "github.com/go-sql-driver/mysql"
)

type (
	Service struct {
		ctx          context.Context
		db           *sql.DB
		hostname     string
		In           chan []byte
		seenGateways map[string]bool
	}
)

// https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin

//
//
//
func mergeDSN(user, pwd, dbhost, dbname string) string {
	return fmt.Sprint("%s:%s@tcp(%s)/%s", user, pwd, dbhost, dbname)
}

//
// NewService
//
func NewService(ctx context.Context, dsn, hostname string) (*Service, error) {

	var (
		err error
	)

	s := &Service{
		hostname: hostname,
	}

	if s.db, err = sql.Open("mysql", dsn); err != nil {
		return nil, err
	}

	if err = s.db.Ping(); err != nil {
		return nil, err
	}

	return s, nil

}

// {"host":"rigaer28-v3","num":5,"alias":"Wohnzimmer","phenomenontime":1505763276,"lon":5.1,"lat":5.2,"alt":5.3,"temp":18.7,"pressure":56,"humidity":-999,"lowbattery":false}

//
// Close
//
func (s *Service) Close() {
	err := s.db.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("==> Database closed")

}

//
//
//
func WriteDB(in []byte, s *Service) error {

	var err error

	dto := &fcc.MeasurementDTO{}

	if err = json.Unmarshal(in, dto); err != nil {
		return err
	}

	return err
}

//
//
//
func (s *Service) Run(ctx context.Context) error {

	var (
		err error
		wg  *sync.WaitGroup
	)

	if wg, err = fcc.GetWGFromContext(ctx); err != nil {
		return err
	}

	wg.Add(1)
	defer wg.Done()

	for {
		select {

		case <-ctx.Done():
			s.Close()
			fmt.Println("==> Leaving database")
			return ctx.Err()
		case in := <-s.In:

			m := &fcc.MeasurementDTO{}

			json.Unmarshal(in, m)

			//fmt.Println(m)

			if s.seenGateways[m.Host] {
				continue
			}

			if !s.exists(m.Host) {
				fmt.Println("==> Doesn't exist", m.Host)
			}

			//			s.seenGateways[m.Alias] = true

		}
	}

}

//
//
//
func (s *Service) exists(alias string) bool {

	h := NewHubDAO(s.db)

	b, err := h.RowExists(alias)
	if err != nil {
		fmt.Println("-->", err)
	}

	return b
}
