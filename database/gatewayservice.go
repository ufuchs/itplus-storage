package database

import (
	"database/sql"
)

//
//
//
type GatewayService struct {
	db   *sql.DB
	list map[string]*Gateway
}

//
//
//
func NewGatewayService(db *sql.DB) *GatewayService {
	return &GatewayService{
		db:   db,
		list: map[string]*Gateway{},
	}
}

//
//
//
func CreateTableGateway(db *sql.DB) error {
	return NewGatewayDAO(db).CreateMyTable()
}

//
//
//
func (s *GatewayService) Insert(g *Gateway) (int64, error) {
	return NewGatewayDAO(s.db).Insert(g)
}

//
//
//
func (s *GatewayService) GatewayExists(host string) (int64, error) {
	return NewGatewayDAO(s.db).GatewayExists(host)
}

func (s *GatewayService) RetrieveAll() error {
	var err error
	s.list, err = NewGatewayDAO(s.db).RetrieveAll()
	return err
}
