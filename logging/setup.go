package logging

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sakokazuki/simplegrpc/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// DammyCloser ...
type DammyCloser struct{}

// Close is dammy
func (f *DammyCloser) Close() error {
	return nil
}

// Setup zerolog
func Setup(config config.Config) io.Closer {
	//set output level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Debug == true {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	//output writer settings
	var closer io.Closer
	var writer io.Writer
	if config.Debug == true {
		closer = &DammyCloser{}
		writer = os.Stdout

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: writer})
	} else {
		file := &lumberjack.Logger{
			Filename:   config.LogRotate.Filename,
			MaxSize:    config.LogRotate.MaxSize, // megabytes
			MaxBackups: config.LogRotate.MaxBackups,
			MaxAge:     config.LogRotate.MaxAge,   //days
			Compress:   config.LogRotate.Compress, // disabled by default
		}

		closer = file
		writer = io.MultiWriter(file, os.Stdout)

		//nopretty
		log.Logger = log.Output(writer)
	}

	return closer
}
