package model

import (
	"gorm.io/gorm"

	. "github.com/m1stborn/mistChatbot/internal/pkg/logger"
	log "github.com/sirupsen/logrus"
)

type YtVideo struct {
	gorm.Model

	VideoId string
}

func (d *Database) CreateYtVideo(vid *YtVideo) {
	if err := d.db.Create(vid).Error; err != nil {
		Log.WithFields(log.Fields{
			"pkg":  "model",
			"func": "CreateYtVideo",
		}).Error(err)
	}
}

func (d *Database) DeleteYtVideo(videoId string) {
	var vid YtVideo
	result := d.db.Where(&YtVideo{VideoId: videoId}).Unscoped().Delete(&vid)
	if result.Error != nil {
		Log.WithFields(log.Fields{
			"pkg":  "model",
			"func": "DeleteYtVideo",
		}).Error(result.Error)
	} //no need to handle if user not vid anything
}

func (d *Database) QueryAllVideoIds() []string {
	var (
		videos   []YtVideo
		videoIds []string
	)
	if err := d.db.Find(&videos).Error; err != nil {
		Log.WithFields(log.Fields{
			"pkg":  "model",
			"func": "QueryAllVideo",
		}).Error(err)
	}
	for _, vid := range videos {
		videoIds = append(videoIds, vid.VideoId)
	}
	return videoIds
}
