package routine

import (
	"fmt"
	"github.com/v2pro/plz/logging"
	"runtime"
	"runtime/debug"
	"time"
)

var panicLogger = logging.LoggerOf("metric", "counter", "panic", "routine")

func Go(oneOff func(), kv ...interface{}) error {
	err := Spi.BeforeStart(kv)
	if err != nil {
		return err
	}
	_, callerFile, callerLine, _ := runtime.Caller(2)
	go func() {
		defer func() {
			recovered := recover()
			if recovered != nil {
				Spi.AfterPanic(recovered, append(kv, "caller", fmt.Sprintf("%s:%d", callerFile, callerLine)))
			}
		}()
		oneOff()
		Spi.AfterFinish(kv)
	}()
	return nil
}

func GoLongRunning(longRunning func(), kv ...interface{}) error {
	err := Spi.BeforeStart(kv)
	if err != nil {
		return err
	}
	go func() {
		for restartedTimes := 0; goLongRunningOnce(longRunning, kv); restartedTimes++ {
			shouldRestartAgain := Spi.BeforeRestart(restartedTimes, kv)
			if !shouldRestartAgain {
				break
			}
		}
		Spi.AfterFinish(kv)
	}()
	return nil
}

func goLongRunningOnce(longRunning func(), kv []interface{}) (notDone bool) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			Spi.AfterPanic(recovered, kv)
			notDone = true
		}
	}()
	longRunning()
	return false
}

type Config struct {
	AfterPanic    func(recovered interface{}, kv []interface{})
	BeforeRestart func(restartedTimes int, kv []interface{}) bool
	BeforeStart   func(kv []interface{}) error
	AfterFinish   func(kv []interface{})
}

var Spi = Config{
	AfterPanic: func(recovered interface{}, kv []interface{}) {
		panicLogger.Error("goroutine panic", append(kv, "recovered", recovered, "stack", string(debug.Stack()))...)
	},
	BeforeRestart: func(restartedTimes int, kv []interface{}) bool {
		time.Sleep(100 * time.Microsecond)
		return true
	},
	BeforeStart: func(kv []interface{}) error {
		return nil
	},
	AfterFinish: func(kv []interface{}) {
		// no op
	},
}

func (cfg *Config) Append(newCfg Config) {
	if newCfg.AfterPanic != nil {
		oldAfterPanic := cfg.AfterPanic
		cfg.AfterPanic = func(recovered interface{}, kv []interface{}) {
			oldAfterPanic(recovered, kv)
			newCfg.AfterPanic(recovered, kv)
		}
	}
	if newCfg.BeforeRestart != nil {
		oldBeforeRestart := cfg.BeforeRestart
		cfg.BeforeRestart = func(restartedTimes int, kv []interface{}) bool {
			if !oldBeforeRestart(restartedTimes, kv) {
				return false
			}
			return newCfg.BeforeRestart(restartedTimes, kv)
		}
	}
	if newCfg.BeforeStart != nil {
		oldBeforeStart := cfg.BeforeStart
		cfg.BeforeStart = func(kv []interface{}) error {
			err := oldBeforeStart(kv)
			if err != nil {
				return err
			}
			return newCfg.BeforeStart(kv)
		}
	}
	if newCfg.AfterFinish != nil {
		oldAfterFinish := cfg.AfterFinish
		cfg.AfterFinish = func(kv []interface{}) {
			oldAfterFinish(kv)
			newCfg.AfterFinish(kv)
		}
	}
}
