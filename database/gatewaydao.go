package database

import (
	"database/sql"
)

type (
	HubDAO struct {
		db *sql.DB
	}
)

//
//
//
func NewHubDAO(db *sql.DB) *HubDAO {
	return &HubDAO{
		db: db,
	}
}

//
//
//
func (d *HubDAO) GatewayExists(alias string) (int, error) {
	var gwID int

	err := d.db.QueryRow("SELECT GatewayID FROM gateway WHERE alias = ?", alias).Scan(&gwID)

	return gwID, err
}

//
//
//
func (d *HubDAO) RetrieveAll(hubID int) (GatewayList, error) {

	list := GatewayList{}

	rows, err := d.db.Query("select GatewayID, HubID, hostname, alias from gateway where HubID = ?", hubID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		gw, err := NewGateway(rows)

		if err != nil {
			return nil, err
		}

		list = append(list, gw)

	}

	return list, rows.Err()

}

//
//
//
func (d *HubDAO) Insert(hubID int, gw *Gateway) (int64, error) {

	var insStmt = "INSERT INTO gateway (HubID, GatewayType, Hostname, Alias) VALUES (?, ?, ?, ?)"

	stmt, err := d.db.Prepare(insStmt)
	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(hubID, gw.GatewayType, gw.Hostname, gw.Alias)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

//
//
//
func (d *HubDAO) Update() error {
	return nil
}

//
//
//
func (d *HubDAO) Delete() error {
	return nil
}
