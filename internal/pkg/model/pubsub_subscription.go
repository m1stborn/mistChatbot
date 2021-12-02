package model

import (
	"gorm.io/gorm"

	. "github.com/m1stborn/mistChatbot/internal/pkg/logger"
	log "github.com/sirupsen/logrus"
)

// TO restore previous PubSub subscription after app shutdown

type PubSubSubscription struct {
	gorm.Model

	Topic      string
	CallbackId int
}

func (d *Database) CreatePubSubSubscription(pubsub *PubSubSubscription) {
	if err := d.db.Create(pubsub).Error; err != nil {
		Log.WithFields(log.Fields{
			"pkg":  "model",
			"func": "CreatePubSubSubscription",
		}).Error(err)
	}
}

func (d *Database) DeletePubSubSubscription(topic string) {
	var pubsub PubSubSubscription
	result := d.db.Where(&PubSubSubscription{Topic: topic}).Unscoped().Delete(&pubsub)
	if result.Error != nil {
		Log.WithFields(log.Fields{
			"pkg":  "model",
			"func": "DeletePubSubSubscription",
		}).Error(result.Error)
	} //no need to handle if user not sub anything
}

func (d *Database) QueryAllPubSub() []PubSubSubscription {
	var pubsubs []PubSubSubscription
	if err := d.db.Find(&pubsubs).Error; err != nil {
		Log.WithFields(log.Fields{
			"pkg":  "model",
			"func": "QueryAllPubSub",
		}).Error(err)

	}
	return pubsubs
}
