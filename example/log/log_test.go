package main

import (
	"errors"
	"github.com/aberic/gnomon/log"
	"testing"
	"time"
)

func logDo() {
	log.Debug("test", log.Field("1", "2"), log.Field("2", 3), log.Field("3", true))
	log.Debug("test", nil)
	log.Info("test", log.Field("1", "2"), log.Field("2", 3), log.Field("3", true))
	log.Warn("test", log.Field("1", "2"), log.Field("2", 3), log.Field("3", true))
	log.Error("test", log.Field("1", "2"), log.Field("2", 3), log.Field("3", true), log.Err(errors.New("yes")))
	log.DebugSkip(1, "test", log.Field("1", "2"), log.Field("2", 3), log.Field("3", true))
	log.DebugSkip(1, "test", nil)
	log.InfoSkip(1, "test", log.Field("1", "2"), log.Field("2", 3), log.Field("3", true))
	log.WarnSkip(1, "test", log.Field("1", "2"), log.Field("2", 3), log.Field("3", true))
	log.ErrorSkip(1, "test", log.Field("1", "2"), log.Field("2", 3), log.Field("3", true), log.Err(errors.New("yes")))
}

func TestLogDebug(t *testing.T) {
	log.Set(log.DebugLevel(), "tmp/log", 1, 1, false, false)
	logDo()
	time.Sleep(3 * time.Second)
}
