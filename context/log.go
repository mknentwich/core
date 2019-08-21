package context

import (
	"log"
	"os"
)

//Defines all available log levels.
const (
	LOG_ERROR   LogLevel = "ERROR"
	LOG_INFO    LogLevel = "INFO"
	LOG_WARNING LogLevel = "WARNING"
)

//Represents a log level which is also used in the log.
type LogLevel string

//Defines the function type for a single log.
type Log func(level LogLevel, format string, a interface{})

//The central logger.
var logger = log.New(os.Stdout, "", log.LstdFlags)
