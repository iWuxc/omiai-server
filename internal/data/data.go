package data

import (
	"omiai-server/internal/conf"
	"omiai-server/pkg/db2"

	"github.com/google/wire"
	"github.com/iWuxc/go-wit/database"
	"github.com/iWuxc/go-wit/utils"
)

var ProviderDataSet = wire.NewSet(
	NewDB,
	NewStorage,
)

type DB struct {
	db2.DB
}

// NewDB init gorm .
func NewDB() (db *DB, f func(), e error) {
	var d db2.DB
	dbConf := new(database.DataConf)
	if err := utils.Copy(conf.GetConfig().Database, &dbConf); err != nil {
		return nil, nil, err
	}
	dbConf.Driver = conf.GetConfig().Database.Default.Driver
	dbConf.Source = conf.GetConfig().Database.Default.Source
	d, f, e = db2.NewDataBase(dbConf)
	db = &DB{d}
	return
}
