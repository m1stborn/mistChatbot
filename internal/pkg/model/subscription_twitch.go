package model

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model

	Line            string
	LineAccessToken string //line notify access token

	TwitchLoginName string
	//TODO TwitchEventSubID
}

func (d *Database) CreateSubscription(sub *Subscription) {
	if err := d.db.Create(sub).Error; err != nil {
		//TODO handle error
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "CreateSubscription",
		}).Error(err)
	}
}

func (d *Database) UpdateSubscription(sub *Subscription) {
	if err := d.db.Model(sub).Select([]string{"user", "twitch_id", "twitch_login_name"}).Updates(sub).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "UpdateSubscription",
		}).Error(err)
	}
}

func (d *Database) QuerySubByTwitchLoginName(twitchLoginName string) []string {
	var (
		subs       []Subscription
		subsIDList []string
	)

	if err := d.db.Where(&Subscription{TwitchLoginName: twitchLoginName}).Find(&subs).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "QuerySubByTwitchLoginName",
		}).Error(err)
	}

	for _, sub := range subs {
		subsIDList = append(subsIDList, sub.LineAccessToken)
	}
	return subsIDList
}

func (d *Database) CheckStreamerExist(twitchLoginName string) bool {
	var sub Subscription
	if err := d.db.First(&sub, "twitch_login_name = ?", twitchLoginName).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "CheckStreamerExist",
		}).Error(err)
		return false
	}

	return true
}

func (d *Database) QuerySubByUser(accountID string) []Subscription {
	var subs []Subscription
	if err := d.db.Where(&Subscription{Line: accountID}).Find(&subs).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "QuerySubByUser",
		}).Error(err)
	}
	return subs
}

var ErrRecordNotExist = errors.New("wrong streamer name or not sub yet")

func (d *Database) DeleteSubByUserBroadcaster(accountID string, broadcaster string) error {
	var sub Subscription
	result := d.db.Where(&Subscription{Line: accountID, TwitchLoginName: broadcaster}).Unscoped().Delete(&sub)
	if result.Error != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "DeleteSubByUserBroadcaster",
		}).Error(result.Error)
		return result.Error
	} else if result.RowsAffected < 1 {
		return ErrRecordNotExist
	}
	return nil
}

func (d *Database) DeleteSubUserUnfollow(accountID string) {
	var sub Subscription
	result := d.db.Where(&Subscription{Line: accountID}).Unscoped().Delete(&sub)
	if result.Error != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "DeleteSubUserUnfollow",
		}).Error(result.Error)
	} //no need to handle if user not sub anything

}
