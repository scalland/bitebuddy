package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Utils struct {
	debug    bool
	TimeZone *time.Location
}

var ut *Utils

func NewUtils() *Utils {
	if ut == nil {
		ut = &Utils{
			debug: false,
		}
	}
	return ut
}

func GetUtils() *Utils {
	return ut
}

func (u *Utils) GetDebug() bool {
	return u.debug
}

func (u *Utils) SetDebug() bool {
	var debugErr error
	debugStr := os.Getenv("DEBUG")
	if debugStr == "" {
		slog.Debug("utils.SetDebug: environment variable DEBUG has not been set or is an empty string. Reading from config file")
		u.debug = viper.GetBool("debug")
	} else {
		u.debug, debugErr = strconv.ParseBool(debugStr)
		if debugErr != nil {
			slog.Debug("utils.SetDebug: environment variable DEBUG has not been set or is an empty string. Reading from config file")
			u.debug = viper.GetBool("debug")
		}
	}

	slog.Debug("utils.SetDebug: utils.Debug is %t", u.debug)

	//if u.debug {
	//	u.Log.SetLevel(log.DebugLevel)
	//}
	return u.debug
}

func (u *Utils) SetTimeZone(tz string) {
	var tzErr error
	slog.Debug("utils.SetTimeZone: trying to load timezone %s...", tz)
	u.TimeZone, tzErr = time.LoadLocation(tz)
	if tzErr != nil {
		slog.Error(fmt.Sprintf("utils.SetTimeZone: error loading time zone: %s: %s", tz, tzErr.Error()))
		slog.Info("utils.SetTimeZone: setting timezone to \"local\"")
		u.TimeZone = time.Local
	}
}

func (u *Utils) GetTimeZone() string {
	return u.TimeZone.String()
}
