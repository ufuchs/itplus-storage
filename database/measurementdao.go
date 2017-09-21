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
	schema  string
}

//
//
//
func NewMeasurementExDAO(db *sql.DB, schema string) *MeasurementExDAO {

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
		schema:  schema,
	}
}

//
//
//
func (d *MeasurementExDAO) CreateTable(tableNmame string) error {

	var txt = `CREATE TABLE ? (
		MeasureID int NOT NULL AUTO_INCREMENT,
		GatewayID int NOT NULL,
		Num            int, 
		Alias          varchar(32),
		PhenomenonTime bigint,   
		Lon            float(18), 
		Lat            float(18),  
		Alt            float(18),  
		Temp           float(18),  
		Pressure       float(18),  
		Humidity       float(18),  
		LowBattery     boolean,    
		PRIMARY KEY (MeasureID),
		CONSTRAINT FK_Measurements_Gateway FOREIGN KEY (GatewayID) REFERENCES gateway(GatewayID)
	);
	`

	stmt, err := d.db.Prepare(txt)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(tableNmame)
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

	return err

}

//
//
//
func (d *MeasurementExDAO) TableExists(tableNmame string) error {

	var s = "SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name LIKE ?"

	rows, err := d.db.Query(s, d.schema, tableNmame+"%")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {

		gw, err := NewGateway(rows)

		if err != nil {
			return nil, err
		}

		list = append(list, gw)

	}

	return nil
}

//
// {"host":"rigaer28-v3","num":5,"alias":"Wohnzimmer","phenomenontime":1505763276,"lon":5.1,"lat":5.2,"alt":5.3,"temp":18.7,"pressure":56,"humidity":-999,"lowbattery":false}
//
func (d *MeasurementExDAO) Insert(gwID int, m *MeasurementEx) error {

	stmt, err := d.db.Prepare(d.insStmt)
	if err != nil {
		return err
	}

	defer stmt.Close()

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
