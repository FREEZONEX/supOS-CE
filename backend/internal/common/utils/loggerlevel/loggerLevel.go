package loggerlevel

import "strings"

const (
	levelAlert  = "alert"
	levelInfo   = "info"
	levelError  = "error"
	levelSevere = "severe"
	levelFatal  = "fatal"
	levelSlow   = "slow"
	levelStat   = "stat"
	levelDebug  = "debug"
)

var CurrentLevel = "INFO"
var DoStats = false

func IsDebug() bool {
	return strings.ToLower(CurrentLevel) == levelDebug
}
