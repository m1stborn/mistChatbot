package model

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type User struct {
	gorm.Model

	Line            string `gorm:"primaryKey;unique"` //line userID
	LineAccessToken string //for notify usage

	Subscriptions []Subscription `gorm:"foreignKey:Line;references:Line"`

	//Type    string //`json:"type,omitempty"`
	//Email   string //`json:"email"`

	//Subscribes []string `json:"Subscribes"`
}

func (d *Database) CreateUser(user *User) {
	if err := d.db.Create(user).Error; err != nil {
		//TODO handle error
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "CreateUser",
		}).Error(err)
	}

	logger.WithFields(log.Fields{
		"pkg":  "model",
		"func": "CreateUser",
	}).Info("Create User Success")
}

func (d *Database) UpdateUser(user *User) {
	if err := d.db.Model(user).Select([]string{"line", "line_access_token"}).Updates(user).Error; err != nil {
		//TODO handle error
		logger.WithFields(log.Fields{
			"func": "UpdateUser",
			"pkg":  "model",
		}).Error(err)
	}

	logger.WithFields(log.Fields{
		"pkg":  "model",
		"func": "UpdateUser",
	}).Info("Update User Success")
}
