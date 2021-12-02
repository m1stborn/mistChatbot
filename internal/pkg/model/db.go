package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	TestChannelIds = []string{
		"UC1DCedRgGHBdm81E1llLhOQ",
		"UC-hM6YJuNYVAmUWxeIr9FeA",
		"UC1opHUrw8rvnsadT-iGp7Cg",
		"UCCzUftO8KOVkV4wQG1vkUvg",
		"UCl_gCybOJRIgOXw6Qb4qJzQ",
		"UCiEm9noegBIb-AzjqpxKffA", //羅傑
		"UCqm3BQLlJfvkTsX_hvm0UmA", //WTM
		"UCMwGHR0BTZuLsmjY_NT5Pwg", //Ina
		"UChgTyjG-pdNvxxhdsXfHQ5Q", //Pavolia
		"UCD8HOxPs4Xvsm8H0ZxXGiBw", //Mel
		"UC_vMYWcDjmfdpH6r4TTn1MQ", //Iroha
		"UCvInZx9h3jC2JzsIzoOebWg", //Flare
		"UC4G-xDOf5U9luBcfpyaqF3Q", //My Channel
		"UCZlDXzGoo7d44bwdNObFacg", //Katana
		"UCK9V2B22uJYu3N7eR_BT9QA", //polka
	}

	TestVideoIds = []string{
		"6hZ-kf1aQ1M",
		"omgSWqwVTjY",
		"IwlECRC8c0E",
		"SHJgH64VN9g",
		"nU63cC_brTo",
	}
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
	if dropErr := d.db.Migrator().DropTable(&YtSubscription{}); dropErr != nil {
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

	if !d.db.Migrator().HasTable(&YtSubscription{}) {
		err = d.db.Migrator().CreateTable(&YtSubscription{})
	} else {
		err = d.db.Migrator().AutoMigrate(&YtSubscription{})
	}

}

func (d *Database) Init(uri string) {

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		//TODO handle error
	}

	d.db = db

	//drop old table
	//TODO: should remove in the future
	//if dropErr := d.db.Migrator().DropTable(&Subscription{}); dropErr != nil {
	//	logger.WithFields(log.Fields{
	//		"pkg":  "model",
	//		"func": "Init",
	//	}).Error(dropErr)
	//}
	//if dropErr := d.db.Migrator().DropTable(&User{}); dropErr != nil {
	//	logger.WithFields(log.Fields{
	//		"pkg":  "model",
	//		"func": "Init",
	//	}).Error(dropErr)
	//}
	//if dropErr := d.db.Migrator().DropTable(&YtSubscription{}); dropErr != nil {
	//	logger.WithFields(log.Fields{
	//		"pkg":  "model",
	//		"func": "Init",
	//	}).Error(dropErr)
	//}

	//create all the table
	if !d.db.Migrator().HasTable(&User{}) {
		err = d.db.Migrator().CreateTable(&User{})
	} else {
		//err = d.db.Migrator().AutoMigrate(&User{})
	}
	if !d.db.Migrator().HasTable(&Subscription{}) {
		err = d.db.Migrator().CreateTable(&Subscription{})
	} else {
		//err = d.db.Migrator().AutoMigrate(&Subscription{})
	}
	if !d.db.Migrator().HasTable(&YtSubscription{}) {
		err = d.db.Migrator().CreateTable(&YtSubscription{})
	} else {
		//err = d.db.Migrator().AutoMigrate(&YtSubscription{})
	}

	if !d.db.Migrator().HasTable(&PubSubSubscription{}) {
		err = d.db.Migrator().CreateTable(&PubSubSubscription{})
		for i, channelId := range TestChannelIds {
			d.CreatePubSubSubscription(&PubSubSubscription{
				Topic:      "https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + channelId,
				CallbackId: i,
			})
		}
	}

	if !d.db.Migrator().HasTable(&YtVideo{}) {
		err = d.db.Migrator().CreateTable(&YtVideo{})
		for _, id := range TestVideoIds {
			d.CreateYtVideo(&YtVideo{
				VideoId: id,
			})
		}
	}

}
