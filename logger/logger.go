package logger

import (
	"io"
	"log"
	"os"

	isatty "github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var Log zerolog.Logger
var err error
var LOG_LEVEL = "error"

func init() {
	Log, err = CreateLogger(LOG_LEVEL)
	if err != nil {
		log.Println(err)
	}
}
func CreateLogger(lvl string) (zerolog.Logger, error) {
	var logFile string
	var lumberjackLogger lumberjack.Logger
	logFile = "/tmp/auction.log"

	// create file if not exists
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lumberjackLogger = lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1024, //MB
		MaxBackups: 10,
		MaxAge:     30, //day
		Compress:   true,
	}
	newLogger, err := newLogger(lvl, io.MultiWriter(os.Stderr, &lumberjackLogger))
	if err != nil {
		return newLogger.With().Logger(), err
	}
	return newLogger, nil
}

// newLogger returns a configured logger.
func newLogger(level string, w io.Writer) (zerolog.Logger, error) {
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return pkgerrors.MarshalStack(errors.WithStack(err))
	}

	logger := zerolog.New(w).With().Stack().Bool("jira_flag", true).Timestamp().Logger()
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	logger = logger.Level(lvl)
	// pretty print during development
	if f, ok := w.(*os.File); ok {
		if isatty.IsTerminal(f.Fd()) {
			logger = logger.Output(zerolog.ConsoleWriter{Out: f})
		}
	}

	// replace standard logger with zerolog
	log.SetFlags(0)
	log.SetOutput(logger)

	return logger, nil
}
