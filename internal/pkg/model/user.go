package model

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type User struct {
	gorm.Model

	Line            string `gorm:"primary_key"` //line userID
	LineAccessToken string //for notify usage

	Subscriptions []Subscription `gorm:"foreign_key:LineUser;references:Line"`

	//Type    string //`json:"type,omitempty"`
	//Email   string //`json:"email"`

	//Subscribes []string `json:"Subscribes"`
}

func (d *Database) CreateUser(user *User) {
	if err := d.db.Create(user).Error; err != nil {
		//TODO handle error
		logger.WithFields(log.Fields{
			"func": "CreateUser",
			"pkg":  "model",
		}).Error(err)
	}

	logger.WithFields(log.Fields{
		"func": "CreateUser",
		"pkg":  "model",
	}).Info("Create User Success")
}

func (d *Database) UpdateUser(user *User) {
	if err := d.db.Model(user).Select([]string{"line", "line_access_token"}).Update(user).Error; err != nil {
		//TODO handle error
		logger.WithFields(log.Fields{
			"func": "UpdateUser",
			"pkg":  "model",
		}).Error(err)
	}

	logger.WithFields(log.Fields{
		"func": "UpdateUser",
		"pkg":  "model",
	}).Info("Update User Success")
}