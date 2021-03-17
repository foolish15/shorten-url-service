package writer

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Hook is a hook that writes logs of specified LogLevels to specified Writer
type Hook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *Hook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		entry.Data["requestId"] = entry.Context.Value(echo.HeaderXRequestID)
	}
	line, err := entry.Bytes()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *Hook) Levels() []logrus.Level {
	return hook.LogLevels
}
