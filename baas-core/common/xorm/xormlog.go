package xorm

import (
	"github.com/go-xorm/core"
	"fmt"
	"github.com/op/go-logging"
)

//xorm 日志实体
type OrmLogger struct {
	showSQL bool
	level   core.LogLevel
	logger  *logging.Logger
}

// Error implement core.ILogger
func (s *OrmLogger) Error(v ...interface{}) {
	if s.level <= core.LOG_ERR {
		s.logger.Error(fmt.Sprint(v...))
	}
	return
}

// Errorf implement core.ILogger
func (s *OrmLogger) Errorf(format string, v ...interface{}) {
	if s.level <= core.LOG_ERR {
		s.logger.Errorf(format, v)
	}
	return
}

// Debug implement core.ILogger
func (s *OrmLogger) Debug(v ...interface{}) {
	if s.level <= core.LOG_DEBUG {
		s.logger.Debug(fmt.Sprint(v...))
	}
	return
}

// Debugf implement core.ILogger
func (s *OrmLogger) Debugf(format string, v ...interface{}) {
	if s.level <= core.LOG_DEBUG {
		s.logger.Debugf(format, v)
	}
	return
}

// Info implement core.ILogger
func (s *OrmLogger) Info(v ...interface{}) {
	if s.level <= core.LOG_INFO {
		s.logger.Info(fmt.Sprint(v...))
	}
	return
}

// Infof implement core.ILogger
func (s *OrmLogger) Infof(format string, v ...interface{}) {
	if s.level <= core.LOG_INFO {
		s.logger.Infof(format, v)
	}
	return
}

// Warn implement core.ILogger
func (s *OrmLogger) Warn(v ...interface{}) {
	if s.level <= core.LOG_WARNING {
		s.logger.Warning(fmt.Sprint(v...))
	}
	return
}

// Warnf implement core.ILogger
func (s *OrmLogger) Warnf(format string, v ...interface{}) {
	if s.level <= core.LOG_WARNING {
		s.logger.Warningf(format, v...)
	}
	return
}

// Level implement core.ILogger
func (s *OrmLogger) Level() core.LogLevel {
	return s.level
}

// SetLevel implement core.ILogger
func (s *OrmLogger) SetLevel(l core.LogLevel) {
	s.level = l
	return
}

// ShowSQL implement core.ILogger
func (s *OrmLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (s *OrmLogger) IsShowSQL() bool {
	return s.showSQL
}
