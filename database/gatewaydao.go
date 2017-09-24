package database

import (
	"database/sql"
)

type (
	GatewayDAO struct {
		existStmt    string
		retrieveStmt string
		insStmt      string
		createStmt   string
		AbstractDAO
	}
)

//
//
//
func NewGatewayDAO(db *sql.DB) *GatewayDAO {
	dao := &GatewayDAO{
		existStmt:    "SELECT GatewayID FROM gateway WHERE alias = ?",
		retrieveStmt: "select GatewayID, hostname, alias from gateway",
		insStmt:      "INSERT INTO gateway (GatewayType, Hostname, Alias) VALUES (?, ?, ?)",
		createStmt: `CREATE TABLE Gateways (
			GatewayID   int NOT NULL AUTO_INCREMENT,
			GatewayType varchar(96),
			Hostname    varchar(32),
			Alias       varchar(32),
			PRIMARY KEY (GatewayID)
		);`,
	}
	dao.AbstractDAO.Db = db
	return dao
}

//
//
//
func (d *GatewayDAO) CreateMyTable() error {
	return d.CreateTable(d.createStmt)
}

//
//
//
func (d *GatewayDAO) GatewayExists(alias string) (int, error) {
	var gwID int

	err := d.Db.QueryRow(d.existStmt, alias).Scan(&gwID)

	return gwID, err
}

//
//
//
func (d *GatewayDAO) RetrieveAll() (GatewayList, error) {

	list := GatewayList{}

	rows, err := d.Db.Query(d.retrieveStmt)
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
func (d *GatewayDAO) Insert(gw *Gateway) (int64, error) {

	stmt, err := d.Db.Prepare(d.insStmt)
	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(gw.GatewayType, gw.Hostname, gw.Alias)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

//
//
//
func (d *GatewayDAO) Update() error {
	return nil
}

//
//
//
func (d *GatewayDAO) Delete() error {
	return nil
}
