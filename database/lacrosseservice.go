package database

import "database/sql"

//
//
//
type LacrosseService struct {
	db             *sql.DB
	gwID           int
	schema         string
	existingTables []string
}

//
//
//
func NewLacrosseService(db *sql.DB, gwID int, schema string) (*LacrosseService, error) {

	var err error

	res := &LacrosseService{
		db:     db,
		gwID:   gwID,
		schema: schema,
	}

	res.existingTables, err = res.getExistingTables("Lacrosse_")

	return res, err

}

//
//
//
func (s *LacrosseService) getExistingTables(alias string) ([]string, error) {
	dao := NewLacrosseDAO(s.db, s.schema)
	return dao.GetKnownTables(alias)
}
