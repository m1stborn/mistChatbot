package model

import (
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger = log.New()

func init() {
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	formatter.ForceColors = true
	formatter.DisableTimestamp = true // Heroku logging already have timestamp

	// Set specific colors for prefix and timestamp
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})

	logger.Formatter = formatter
	logger.Level = log.DebugLevel
	logger.Info("Successfully initializing logrus with color")
}
