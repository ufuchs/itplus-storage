package database

import (
	"database/sql"
)

//
//
//
type GatewayService struct {
	db    *sql.DB
	hubID int
}

//
//
//
func NewGatewayService(db *sql.DB, hubID int) *GatewayService {
	return &GatewayService{
		db:    db,
		hubID: hubID,
	}
}

//
//
//
func (s *GatewayService) AddGateway(hubID int, hostname string) (int64, error) {

	var dao = NewGatewayDAO(s.db)

	dao.CreateMyTable()

	return dao.Insert(hubID, &Gateway{
		HubID:       hubID,
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
