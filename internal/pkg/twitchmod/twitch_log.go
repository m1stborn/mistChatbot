package twitchmod

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger = log.New()

func init() {
	fmt.Println("init log")
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	formatter.ForceColors = true

	// Set specific colors for prefix and timestamp
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})

	logger.Formatter = formatter
	logger.Level = log.DebugLevel
	logger.Info("Success initializing logrus")
}
