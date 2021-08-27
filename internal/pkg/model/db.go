package model

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

var DB = Database{}

func (d *Database) TestInit(uri string) {
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		//TODO handle error
		fmt.Println(err)
	}

	d.db = db

	if dropErr := d.db.Migrator().DropTable(&Subscription{}); dropErr != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "Init",
		}).Error(dropErr)
	}
	if dropErr := d.db.Migrator().DropTable(&User{}); dropErr != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "Init",
		}).Error(dropErr)
	}

	//create all the table
	if !d.db.Migrator().HasTable(&User{}) {
		err = d.db.Migrator().CreateTable(&User{})
	} else {
		err = d.db.Migrator().AutoMigrate(&User{})
	}

	if !d.db.Migrator().HasTable(&Subscription{}) {
		err = d.db.Migrator().CreateTable(&Subscription{})
	} else {
		err = d.db.Migrator().AutoMigrate(&Subscription{})
	}

}

func (d *Database) Init(uri string) {

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		//TODO handle error
	}

	d.db = db

	//drop old table
	if dropErr := d.db.Migrator().DropTable(&Subscription{}); dropErr != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "Init",
		}).Error(dropErr)
	}
	if dropErr := d.db.Migrator().DropTable(&User{}); dropErr != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "Init",
		}).Error(dropErr)
	}

	//create all the table
	if !d.db.Migrator().HasTable(&User{}) {
		err = d.db.Migrator().CreateTable(&User{})
	} else {
		err = d.db.Migrator().AutoMigrate(&User{})
	}

	if !d.db.Migrator().HasTable(&Subscription{}) {
		err = d.db.Migrator().CreateTable(&Subscription{})
	} else {
		err = d.db.Migrator().AutoMigrate(&Subscription{})
	}

}
