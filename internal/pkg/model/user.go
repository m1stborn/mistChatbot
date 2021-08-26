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
	Enable        bool           `gorm:"default:true"`

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
	}).Info("Create user success")
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
	}).Info("Update user success")
}

func (d *Database) CheckLineAccessTokenExist(accountID string) bool {
	var user User
	if err := d.db.First(&user, "line = ?", accountID).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "CheckLineAccessTokenExist",
		}).Error(err)
		return false
	}
	if user.LineAccessToken != "" {
		return false
	}
	return true
}

func (d *Database) GetUser(accountID string) *User {
	var user User
	if err := d.db.First(&user, "line = ?", accountID).Error; err != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "GetUser",
		}).Error(err)
		return nil
	}
	return &user
}

func (d *Database) UserUnfollow(accountID string) {
	result := d.db.Model(&User{}).Where("line = ?", accountID).Update("enable", false)
	if result.Error != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "UserUnfollow",
		}).Error(result.Error)
		return
	} else if result.RowsAffected < 1 {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "UserUnfollow",
		}).Info("No user updated")
		return
	}

	return
}

func (d *Database) UserConnectNotify(accountID string, accessToken string) error {
	result := d.db.Model(&User{}).Where("line = ?", accountID).Update("line_access_token", accessToken)
	if result.Error != nil {
		logger.WithFields(log.Fields{
			"pkg":  "model",
			"func": "UserConnectNotify",
		}).Error(result.Error)
		return result.Error
	}
	return nil
}
