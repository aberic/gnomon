package gnomon

import (
	"errors"
	"testing"
	"time"
)

var logDir = "./log"

func TestLog(t *testing.T) {
	Log().Set(debugLevel, false)
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Debug("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Info("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	Log().Warn("test", LogField("1", "2"), LogField("2", 3), LogField("3", true))
	time.Sleep(time.Second)
}

func TestLogError(t *testing.T) {
	Log().Set(debugLevel, false)
	err := errors.New("err err")
	Log().Error("test", LogField("1", "2"), LogField("2", 3), LogField("3", true), LogErr(err))
	time.Sleep(time.Second)
}

func TestLogPanic(t *testing.T) {
	Log().Set(debugLevel, false)
	err := errors.New("err err")
	Log().Panic("test", LogField("1", "2"), LogField("2", 3), LogField("3", true), LogErr(err))
	time.Sleep(time.Second)
}

func TestLogFatal(t *testing.T) {
	Log().Set(debugLevel, false)
	err := errors.New("err err")
	Log().Fatal("test", LogField("1", "2"), LogField("2", 3), LogField("3", true), LogErr(err))
	time.Sleep(time.Second)
}

func TestLogWithStorage(t *testing.T) {
	Log().Init(logDir, 1, 1, false)
	Log().Set(debugLevel, false)
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
	time.Sleep(time.Second)
}
