package fuel

import (
	"fmt"
	"time"
)

type Logger interface {
	Infof(template string, args ...any)
}

type SimpleStdoutLogger struct{}

func (l SimpleStdoutLogger) Infof(template string, args ...any) {
	fmt.Printf(time.Now().String()+" "+template+"\n", args...)
}
