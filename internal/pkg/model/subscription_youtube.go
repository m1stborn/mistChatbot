package model

import (
	//"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type YtSubscription struct {
	gorm.Model

	Line            string
	LineAccessToken string //line notify access token

	ChannelId string
}

func (d *Database) CreateYtSubscription(sub *YtSubscription) {
	if err := d.db.Create(sub).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "CreateYtSubscription",
		}).Error(err)
	}
}

// UpdateYtSubscription TODO: verify this function
func (d *Database) UpdateYtSubscription(sub *YtSubscription) {
	if err := d.db.Model(sub).Select([]string{"user", "channel)di"}).Updates(sub).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "UpdateYtSubscription",
		}).Error(err)
	}
}

func (d *Database) QuerySubByYtChannelId(channelId string) []string {
	var (
		subs       []YtSubscription
		subsIDList []string
	)

	if err := d.db.Where(&YtSubscription{ChannelId: channelId}).Find(&subs).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "QuerySubByYtChannelId",
		}).Error(err)
	}

	for _, sub := range subs {
		subsIDList = append(subsIDList, sub.LineAccessToken)
	}
	return subsIDList
}

func (d *Database) CheckYtChannelExist(channelId string) bool {
	var sub YtSubscription
	//TODO: verify table column "channel_id"
	if err := d.db.First(&sub, "channel_id = ?", channelId).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "CheckYtChannelExist",
		}).Error(err)
		return false
	}

	return true
}

func (d *Database) QueryYtSubByUser(accountID string) []YtSubscription {
	var subs []YtSubscription
	if err := d.db.Where(&YtSubscription{Line: accountID}).Find(&subs).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "QueryYtSubByUser",
		}).Error(err)
	}
	return subs
}

func (d *Database) DeleteSubByUserChannelId(accountID string, channelId string) error {
	var sub Subscription
	result := d.db.Where(&YtSubscription{Line: accountID, ChannelId: channelId}).Unscoped().Delete(&sub)
	if result.Error != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "DeleteSubByUserChannelId",
		}).Error(result.Error)
		return result.Error
	} else if result.RowsAffected < 1 {
		return ErrRecordNotExist
	}
	return nil
}

func (d *Database) DeleteYtSubUserUnfollow(accountID string) {
	var sub YtSubscription
	result := d.db.Where(&YtSubscription{Line: accountID}).Unscoped().Delete(&sub)
	if result.Error != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "DeleteYtSubUserUnfollow",
		}).Error(result.Error)
	} //no need to handle if user not sub anything

}
