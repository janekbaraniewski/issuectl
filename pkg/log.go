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

type Logger struct {
	verbosity int
}

func NewLogger() *Logger {
	v := flag.Int("v", 1, "verbosity level")
	flag.Parse()

	return &Logger{verbosity: *v}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(1, fmt.Sprintf(format, v...))
}

func (l *Logger) V(level int) *verbosityLogger {
	return &verbosityLogger{parent: l, level: level}
}

type verbosityLogger struct {
	parent *Logger
	level  int
}

func (v *verbosityLogger) Infof(format string, args ...interface{}) {
	v.parent.output(v.level, fmt.Sprintf(format, args...))
}

func (l *Logger) output(level int, message string) {
	if l.verbosity >= level {
		_, file, _, _ := runtime.Caller(2)
		file = file[strings.LastIndex(file, "/")+1:]
		log.SetOutput(os.Stdout)
		log.SetFlags(0)
		log.Println(fmt.Sprintf("%v|%s| %s", time.Now().Format(time.RFC3339), file, message))
	}
}

var Log = NewLogger()
