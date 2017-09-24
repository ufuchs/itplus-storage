package database

import (
	"database/sql"
)

//
//
//
type GatewayService struct {
	db *sql.DB
}

//
//
//
func NewGatewayService(db *sql.DB) *GatewayService {
	return &GatewayService{
		db: db,
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
func (s *GatewayService) AddGateway(hostname string) (int64, error) {

	var dao = NewGatewayDAO(s.db)

	dao.CreateMyTable()

	return dao.Insert(&Gateway{
		GatewayType: "",
		Hostname:    hostname,
		Alias:       "",
	})

}

//
//
//
func (s *GatewayService) GatewayExists(host string) (int, error) {
	dao := NewGatewayDAO(s.db)
	return dao.GatewayExists(host)
}
