package database

import "github.com/ufuchs/itplus/base/fcc"

type Hub struct {
	HubID    int
	Hostname string
	Alias    string
}

//
//
//
type MeasurementEx struct {
	MeasureID int
	GatewayID int
	fcc.MeasurementDTO
}
