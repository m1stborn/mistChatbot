package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	db *gorm.DB
}

var DB = Database{}

//func NewDatabase(uri string) Database {
//	db, err := gorm.Open("postgres", uri)
//	if err != nil {
//		//handle error
//	}
//	return Database{db}
//}

func (d *Database) Init(uri string) {

	db, err := gorm.Open("postgres", uri)
	if err != nil {
		//TODO handle error
	}

	d.db = db

	//create all the table
	if !d.db.HasTable(&User{}) {
		err = d.db.CreateTable(&User{}).Error
	} else {
		err = d.db.AutoMigrate(&User{}).Error
	}

	if !d.db.HasTable(&Subscription{}) {
		err = d.db.CreateTable(&Subscription{}).Error
	} else {
		err = d.db.AutoMigrate(&Subscription{}).Error
	}

}
