package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ufuchs/itplus/base/fcc"
	"github.com/ufuchs/itplus/storage/app"

	_ "github.com/go-sql-driver/mysql"
)

type (
	Service struct {
		ctx            context.Context
		Db             *sql.DB
		GatewayService *GatewayService
		hostname       string
		In             chan []byte
		seenGateways   map[string]bool
	}
)

// https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin

//
//
//
func openDatabase(dsn string) (db *sql.DB, err error) {

	if db, err = sql.Open("mysql", dsn); err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}

	return
}

//
// NewService
//
func NewService(ctx context.Context, dsn, hostname string) (*Service, error) {

	var (
		err error
	)

	s := &Service{
		hostname:     hostname,
		seenGateways: map[string]bool{},
	}

	if s.Db, err = openDatabase(dsn); err != nil {
		return nil, err
	}

	s.GatewayService = NewGatewayService(s.Db)

	return s, err

}

//
// Prepare
//
func (s *Service) Prepare() error {

	err := s.GatewayService.RetrieveAll()
	if err != nil {
		return err
	}

	/*
		for _, gw := range gwList {
			s.seenGateways[gw.Hostname] = true
		}
	*/

	return nil

}

//
//
//
func (s *Service) CreateTableGateway() error {
	return NewGatewayDAO(s.Db).CreateMyTable()
}

//
//
//
func (s *Service) PopulateTableGateway(gateways []string) error {
	var err error

	for _, hostname := range app.Gateways {

		g := &Gateway{
			GatewayType: "",
			Hostname:    hostname,
			Alias:       "",
		}

		fmt.Println(hostname)

		if g.GatewayID, err = s.GatewayService.Insert(g); err != nil {
			return err
		}

	}

	return nil
}

// {"host":"rigaer28-v3","num":5,"alias":"Wohnzimmer","phenomenontime":1505763276,"lon":5.1,"lat":5.2,"alt":5.3,"temp":18.7,"pressure":56,"humidity":-999,"lowbattery":false}

//
// Close
//
func (s *Service) Close() {
	err := s.Db.Close()
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
		hubID = 1
		err   error
		wg    *sync.WaitGroup
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

			var err error

			m := &fcc.MeasurementDTO{}
			if err = json.Unmarshal(in, m); err != nil {
				fmt.Println(err)
				continue
			}

			if s.seenGateways[m.Host] {
				fmt.Printf("==> Writing to DB - %v\n", m.Alias)
				continue
			}

			_, err = s.exists(m.Host)
			if err != nil {
				s.addGateway(hubID, m)
			}

		}
	}

}

func (s *Service) addGateway(hubID int, m *fcc.MeasurementDTO) {

	var err error

	fmt.Printf("==> '%v' doesn't exist - ", m.Host)
	var h = NewGatewayDAO(s.Db)
	var gw = &Gateway{
		GatewayType: "",
		Hostname:    m.Host,
		Alias:       "",
	}
	if _, err = h.Insert(gw); err != nil {
		fmt.Printf("adding failed - '%v'\n", err)
	} else {
		fmt.Println("added")
		s.seenGateways[m.Host] = true
	}

}

//
//
//
func (s *Service) exists(alias string) (int64, error) {
	return NewGatewayDAO(s.Db).GatewayExists(alias)
}
