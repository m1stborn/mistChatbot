package model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func (d *Database) TestInit(uri string) {
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		//TODO handle error
	}

	d.db = db
}

func (d *Database) Init(uri string) {

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		//TODO handle error
	}

	d.db = db

	//drop old table
	if !d.db.Migrator().HasTable(&Subscription{}) {
		err = d.db.Migrator().DropTable(&Subscription{})
	}
	if !d.db.Migrator().HasTable(&User{}) {
		err = d.db.Migrator().DropTable(&User{})
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
