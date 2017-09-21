package database

import (
	"github.com/ufuchs/itplus/base/fcc"
)

var DSN = "itplus:abc12345@tcp(192.168.178.12:3306)/itp_home"

//
//
//
type MeasurementEx struct {
	MeasureID int
	GatewayID int
	fcc.MeasurementDTO
}
