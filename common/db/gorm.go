package db

import (
	"time"

	"codeup.aliyun.com/67c7c688484ca2f0a13acc04/baseTemp/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var db *gorm.DB

func DB() *gorm.DB {
	return db
}
func Initgorm() {
	var err error
	if config.Conf().Pgsql == nil {
		return
	}
	source := make([]gorm.Dialector, 0)
	replicas := make([]gorm.Dialector, 0)
	dnsList := make([]string, 0)
	primaryDB := ""
	for _, v := range *config.Conf().Pgsql {
		dns := v.GetHost()
		dnsList = append(dnsList, dns)
		if v.Primary {
			if primaryDB == "" {
				primaryDB = dns
			} else {

				source = append(source, postgres.Open(dns))
			}
		} else {
			replicas = append(replicas, postgres.Open(dns))
		}
	}
	if primaryDB == "" {
		primaryDB = dnsList[0]
	}
	db, err = gorm.Open(postgres.Open(primaryDB), &gorm.Config{})
	if err != nil {

		panic("failed to connect database")
	}
	db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  source,
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}))
	db.Session(&gorm.Session{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
}
