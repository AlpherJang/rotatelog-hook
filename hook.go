package rotateloghook

import (
	"fmt"
	"io"
	"log"
	"path"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var (
	// We are logging to file, strip coloirs to make the output more readable.
	defaultFormat = &logrus.TextFormatter{DisableColors: true}
	// We are rotate log file every day
	defaultRotateTime = 24 * time.Hour
	// If log file size is more than 200MB, we will rotate file
	defaultRotateSize = 200 * 1024
)

// RotateLogHook use to append log to file
type RotateLogHook struct {
	writer io.Writer
	sync.Mutex
	logDir      string
	logFileName string
	formatter   logrus.Formatter
}

func (h *RotateLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook,err:= NewRotateLogHook("/log","example",defaultFormat)`
// `if err == nil {log.Hooks.Add(hook)}`
func NewRotateLogHook(logDir, filePrefix string, formatter logrus.Formatter) (*RotateLogHook, error) {
	hook := &RotateLogHook{logDir: logDir, logFileName: filePrefix}
	err := hook.SetWriter(nil)
	hook.SetFormatter(formatter)
	return hook, err
}

// SetWriter to custom writer, if send nil, will create log file by rotatelogs libary
// and default directory is logDir with init, so user need to init directroy before.
func (h *RotateLogHook) SetWriter(writer io.Writer) error {
	h.Lock()
	defer h.Unlock()
	var err error
	if writer == nil {
		fullFileName := path.Join(h.logDir, h.logFileName)
		writer, err = rotatelogs.New(
			fmt.Sprintf("%s.%%Y%%m%%d%%H%%M", fullFileName),
			rotatelogs.WithLinkName(fmt.Sprintf("%s.log", fullFileName)),
			rotatelogs.WithMaxAge(-1),
			rotatelogs.WithRotationTime(defaultRotateTime),
			rotatelogs.WithRotationSize(int64(defaultRotateSize)),
		)
	}
	h.writer = writer
	return err
}

// SetFormatter custom formatter to writer , if nil use default
func (h *RotateLogHook) SetFormatter(formatter logrus.Formatter) {
	h.Lock()
	defer h.Unlock()
	if formatter == nil {
		formatter = defaultFormat
	} else {
		switch formatter.(type) {
		case *logrus.TextFormatter:
			textFormatter := formatter.(*logrus.TextFormatter)
			textFormatter.DisableColors = true
		}
	}
	h.formatter = formatter
}

// Fire append to writer
func (h *RotateLogHook) Fire(entry *logrus.Entry) error {
	h.Lock()
	defer h.Unlock()
	line, err := h.formatter.Format(entry)
	if err != nil {
		log.Println("failed to format string for entry: ", err)
		return err
	}
	_, err = h.writer.Write(line)
	return err
}
