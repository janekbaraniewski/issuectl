package issuectl

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// Logger struct holds the verbosity level
type Logger struct {
	verbosity int
}

// NewLogger parses command line argument for verbosity level and returns a new Logger
func NewLogger() *Logger {
	v := flag.Int("v", 1, "verbosity level")
	flag.Parse()

	return &Logger{verbosity: *v}
}

// Infof formats the log message and calls output function with default level 1
func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(1, fmt.Sprintf(format, v...))
}

// V returns a verbosityLogger with the specified verbosity level
func (l *Logger) V(level int) *verbosityLogger {
	return &verbosityLogger{parent: l, level: level}
}

// verbosityLogger is a Logger that has a specific verbosity level
type verbosityLogger struct {
	parent *Logger
	level  int
}

// Infof formats the log message and calls the parent logger's output function
func (v *verbosityLogger) Infof(format string, args ...interface{}) {
	v.parent.output(v.level, fmt.Sprintf(format, args...))
}

// output logs the message if the logger's verbosity level is greater than or equal to the level
func (l *Logger) output(level int, message string) {
	if l.verbosity >= level {
		_, file, _, _ := runtime.Caller(2)
		file = file[strings.LastIndex(file, "/")+1:]
		log.SetOutput(os.Stdout)
		log.SetFlags(0)
		log.Println(fmt.Sprintf("%v|%s| %s", time.Now().Format(time.RFC3339), file, message))
	}
}

// Log is a global logger
var Log = NewLogger()
