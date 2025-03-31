package db

import (
	"baseTemp/common/config"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func DB() *gorm.DB {
	return db
}
func Initgorm() {
	var err error
	if(config.Conf().Pgsql==nil){
		return
	}
	dns := config.Conf().Pgsql.GetHost()

	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {

		panic("failed to connect database")
	}
	db.Session(&gorm.Session{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
}

func MigrateTable() {
	if config.Conf().InitTable == false {
		return
	}
	err := db.AutoMigrate()
	if err != nil {
		panic(err.Error())
	}
}
