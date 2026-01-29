package db2

import (
	"github.com/iWuxc/go-wit/database"
	"github.com/iWuxc/go-wit/log"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var ErrInvalidDBDriver = errors.New("unsupported sql driver .")

type DB = *gorm.DB

// NewDB init gorm .
func NewDataBase(conf *database.DataConf) (DB, func(), error) {
	driver, err := newDBDriver(conf.Driver, conf.Source)
	if err != nil {
		return nil, nil, err
	}

	logLevel := gormLog.Warn
	if conf.Debug {
		logLevel = gormLog.Info
	}

	db, err := gorm.Open(driver, &gorm.Config{
		Logger: NewWithLog(log.GetInstance(), gormLog.Config{
			LogLevel:                  logLevel,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
		}),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxOpenConns(conf.DBConf.MaxOpenConn)
	sqlDB.SetMaxIdleConns(conf.DBConf.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Minute * conf.DBConf.ConnMaxLifeTime)

	return db, func() {
		_ = sqlDB.Close()
	}, nil
}

// newDBDriver get gorm dialect .
func newDBDriver(dbDriver, dbDSN string) (gorm.Dialector, error) {
	switch dbDriver {
	case "mysql":
		return mysql.Open(dbDSN), nil
	case "sqlite":
		return sqlite.Open(dbDSN), nil
	default:
		return nil, ErrInvalidDBDriver
	}
}
