package util

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Output struct {
	io.Writer
}

func (o *Output) Out(format string, args ...any) {
	_, err := fmt.Fprintf(o.Writer, format, args...)
	if err != nil {
		log.Fatalf("write failed: %v", err)
	}
}

func (o *Output) Outf(format string, keyAndValues ...string) {
	result := format
	for i := 0; i+1 < len(keyAndValues); i += 2 {
		key, val := keyAndValues[i], keyAndValues[i+1]
		pattern := fmt.Sprintf("#{%s}", key)
		result = strings.ReplaceAll(result, pattern, val)
	}
	_, err := fmt.Fprint(o.Writer, result)
	if err != nil {
		log.Fatalf("write failed: %v", err)
	}
}

func (o *Output) OutLines(cnt, prefix string) {
	cnt = strings.TrimSpace(cnt)
	if cnt == "" {
		return
	}
	for _, line := range strings.Split(cnt, "\n") {
		o.Out("%s%s\n", prefix, line)
	}
}
