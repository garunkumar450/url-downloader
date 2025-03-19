package src

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Helpful guide: https://betterstack.com/community/guides/logging/zerolog/
func initLogger(filePath string) (err error, logger zerolog.Logger) {
	// Open the log file for writing
	logPath := strings.Split(filePath, ".")[0]
	err = os.MkdirAll(logPath, os.ModePerm)
	if err != nil {
		return err, logger
	}

	logFile := logPath + "/" + getFileName(filePath) + ".log"
	file := &lumberjack.Logger{
		Filename:  logFile,
		MaxSize:   50,
		MaxAge:    7,
		Compress:  true,
		LocalTime: true,
	}

	zerolog.SetGlobalLevel(zerolog.Level(zerolog.DebugLevel))
	zerolog.TimeFieldFormat = "02-Jan-2006 15:04:05.000"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string { // this specifies the file which called the logging statement.
		short := file
		short = strings.Split(filepath.Base(file), ".")[0]
		file = short
		return "[" + file + ":" + strconv.Itoa(line) + "]"
	}

	// Create a console output writer
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	consoleWriter.TimeFormat = zerolog.TimeFieldFormat
	logger = zerolog.New(zerolog.MultiLevelWriter(consoleWriter, file)).With().Timestamp().Logger()
	return nil, logger
}
