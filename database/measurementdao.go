package database

import (
	"database/sql"
)

//
//
//
type MeasurementExDAO struct {
	insStmt string
	db      *sql.DB
}

//
//
//
func NewMeasurementExDAO(db *sql.DB) *MeasurementExDAO {

	var insStmt = `INSERT INTO measurements (
		GatewayID,
		Num, 
		Alias,
		PhenomenonTime,   
		Lon, 
		Lat,  
		Alt,  
		Temp,  
		Pressure,  
		Humidity,  
		LowBattery    
	) VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	)
	`

	return &MeasurementExDAO{
		insStmt: insStmt,
		db:      db,
	}
}

//
// {"host":"rigaer28-v3","num":5,"alias":"Wohnzimmer","phenomenontime":1505763276,"lon":5.1,"lat":5.2,"alt":5.3,"temp":18.7,"pressure":56,"humidity":-999,"lowbattery":false}
//
func (d *MeasurementExDAO) Insert(gwID int, m *MeasurementEx) error {

	stmt, err := d.db.Prepare(d.insStmt)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(gwID, m.Num, m.Alias, m.PhenomenonTime, m.Lon, m.Lat,
		m.Alt, m.Temp, m.Pressure, m.Humidity, m.LowBattery)
	if err != nil {
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
