package gnomon

import (
	"errors"
	"testing"
	"time"
)

var logDir = "./log"

func TestLog(t *testing.T) {
	Log().Set(debugLevel, false)
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Info("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Warn("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
}

func TestLogError(t *testing.T) {
	Log().Set(debugLevel, false)
	err := errors.New("err err")
	Log().Error("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true), Log().Err(err))
	time.Sleep(time.Second)
}

func TestLogPanic(t *testing.T) {
	Log().Set(debugLevel, false)
	err := errors.New("err err")
	Log().Panic("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true), Log().Err(err))
	time.Sleep(time.Second)
}

func TestLogFatal(t *testing.T) {
	Log().Set(debugLevel, false)
	err := errors.New("err err")
	Log().Fatal("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true), Log().Err(err))
	time.Sleep(time.Second)
}

func TestLogWithStorage(t *testing.T) {
	Log().Init(logDir, 1, 1, false)
	Log().Set(debugLevel, false)
	for i := 0; i < 10000; i++ {
		Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	}
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Debug("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Info("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	Log().Warn("test", Log().Field("1", "2"), Log().Field("2", 3), Log().Field("3", true))
	time.Sleep(time.Second)
}
