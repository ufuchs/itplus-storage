package database

import (
	"database/sql"

	"hidrive.com/ufuchs/itplus/base/fcc"
)

var DSN = "itplus:abc12345@tcp(192.168.178.12:3306)/itp_home"

//
//
//
type Gateway struct {
	GatewayID   int
	HubID       int
	GatewayType string
	Hostname    string
	Alias       string
}

//
//
//
type GatewayList []*Gateway

//
//
//
type MeasurementEx struct {
	MeasureID int
	GatewayID int
	fcc.MeasurementDTO
}

//
// NewGateway
//
func NewGateway(rs *sql.Rows) (*Gateway, error) {

	var gw = &Gateway{}

	return gw, rs.Scan(&gw.GatewayID, &gw.HubID, &gw.Hostname, &gw.Alias)

}

//
// Equals
//
func (h *Gateway) Equals(o *Gateway) bool {
	if o == nil {
		return false
	}

	if o == h {
		return true
	}

	return h.HubID == o.HubID && h.Alias == o.Alias && h.Hostname == o.Hostname
}

//
// Contains
//
func (l GatewayList) Contains(o *Gateway) bool {

	if len(l) == 0 {
		return false
	}

	for _, g := range l {

		if g.Equals(o) {
			return true
		}

	}

	return false

}
