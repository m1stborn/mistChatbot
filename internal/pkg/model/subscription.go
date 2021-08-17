package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type Subscription struct {
	gorm.Model

	LineUser string

	TwitchLoginName string
}

func (d *Database) CreateSubscription(sub *Subscription) {
	if err := d.db.Create(sub).Error; err != nil {
		//TODO handle error
		logger.WithFields(log.Fields{
			"func": "CreateSubscription",
			"pkg":  "model",
		}).Error(err)
	}

	logger.WithFields(log.Fields{
		"func": "CreateSubscription",
		"pkg":  "model",
	}).Info("Create Subscription Success")
}

func (d *Database) UpdateSubscription(sub *Subscription) {
	if err := d.db.Model(sub).Select([]string{"user", "twitch_id", "twitch_login_name"}).Update(sub).Error; err != nil {
		//TODO handle error
		logger.WithFields(log.Fields{
			"func": "UpdateSubscription",
			"pkg":  "model",
		}).Error(err)
	}
	logger.WithFields(log.Fields{
		"func": "UpdateSubscription",
		"pkg":  "model",
	}).Info("Update Subscription Success")
}

func (d *Database) QuerySubByTwitchLoginName(twitchLoginName string) []string {
	var (
		subs       []Subscription
		subsIDList []string
	)

	if err := d.db.Where(&Subscription{TwitchLoginName: twitchLoginName}).Find(&subs).Error; err != nil {
		//TODO handle error
		logger.WithFields(log.Fields{
			"func": "QuerySubByTwitchLoginName",
			"pkg":  "model",
		}).Error(err)
	}

	for _, sub := range subs {
		subsIDList = append(subsIDList, sub.LineUser)
	}

	logger.WithFields(log.Fields{
		"func": "QuerySubByTwitchLoginName",
		"pkg":  "model",
	}).Infof(fmt.Sprintf("Query Result: %+v", subsIDList))

	return subsIDList
}
