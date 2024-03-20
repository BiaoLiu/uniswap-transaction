package orm

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Config mysql config.
type Config struct {
	MaxIdleConns  int           // pool
	MaxOpenConns  int           // pool
	MaxLifeTime   time.Duration // connect max lifetime
	SlowThreshold time.Duration // slow threshold
	Source        string        // data source name
	LogLevel      string        // log level
	Logger        log.Logger
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) (db *gorm.DB) {
	logger := NewLogger(
		WithSlowThreshold(c.SlowThreshold),
		WithLogLevel(c.LogLevel),
		WithLogger(c.Logger),
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      c.Source,
		DisableDatetimePrecision: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxOpenConns(c.MaxOpenConns)
	sqlDb.SetMaxIdleConns(c.MaxIdleConns)
	sqlDb.SetConnMaxLifetime(c.MaxLifeTime)
	return
}
