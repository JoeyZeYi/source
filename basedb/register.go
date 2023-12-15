package basedb

var dbMaps = map[string]IDB{}

func init() {
	dbMaps = make(map[string]IDB)
}
func RegisterDB(idb IDB, dbName string) {
	dbMaps[dbName] = idb
}

func GetDB(dbName string) IDB {
	return dbMaps[dbName]
}
