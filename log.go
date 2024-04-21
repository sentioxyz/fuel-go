package fuel

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Infof(template string, args ...any)
}

type simpleLogger log.Logger

func (l *simpleLogger) Infof(template string, args ...any) {
	_ = (*log.Logger)(l).Output(2, fmt.Sprintf(template+"\n", args...))
}

var SimpleLogger *simpleLogger

func init() {
	SimpleLogger = (*simpleLogger)(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile))
}
