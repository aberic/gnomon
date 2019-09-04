package gnomon

import (
	"errors"
	"testing"
	"time"
)

var logDir = "./log"

func TestLog(t *testing.T) {
	Log().Init(logDir, 1, 1, 1, false, false)
	Log().SetLevel(DebugLevel, false)
	for i := 0; i < 10000; i++ {
		Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	}
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Info("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Warn("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	err := errors.New("err err")
	Log().Error("test", LogField("1", "2"), LogField("2", 3), LogField("3", true), LogErr(err))
	Log().Panic("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Fatal("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	time.Sleep(3 * time.Second)
}
