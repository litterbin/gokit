package main

import (
	kitlog "github.com/go-kit/kit/log"
	"os"
)

func main() {
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger.Log("question", "What is the meaning of life?", "answer", 42)
	logger = kitlog.NewContext(logger).With("zz", 1)
	logger = kitlog.NewContext(logger).With("ts", kitlog.DefaultTimestampUTC)
	logger = kitlog.NewContext(logger).With("a", 1)
	logger.Log("question", "What is the meaning of life?", "answer", 42)
}
