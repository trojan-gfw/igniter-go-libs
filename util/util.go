package util

import (
	"bufio"
	"log"
	"runtime"
)

const (
	// logcat is line-buffered
	LogLineBufLen int = 1024
	MaxLogBufLen  int = 16 * 1024 * 1024
)

func init() {
	log.SetPrefix("igniter-golib-log: ")
}

// LogGoRoutineCount log goroutine count to logcat on android
func LogGoRoutineCount() {
	num := runtime.NumGoroutine()
	log.Printf("goroutine num: %d\n", num)
}

// LogGoroutineStackTrace log goroutine stack trace to logcat on android
func LogGoroutineStackTrace() {
	var err error
	bufferedWriter := bufio.NewWriterSize(log.Writer(), LogLineBufLen)
	buf := make([]byte, MaxLogBufLen)
	buf = buf[:runtime.Stack(buf, true)]
	bufferedWriter.WriteString("=== BEGIN goroutine stack dump ===\n")
	total := len(buf)
	start := 0
	end := 0
	remaining := total
	if remaining >= LogLineBufLen {
		end += LogLineBufLen
	} else {
		end += remaining
	}

	log.Printf("buf len is %d\n", total)
	for {
		nn, err := bufferedWriter.Write(buf[start:end])
		if err != nil {
			panic(err)
		}
		start += nn
		remaining -= nn
		if remaining <= 0 {
			break
		}
		if remaining >= LogLineBufLen {
			end = start + LogLineBufLen
		} else {
			end = start + remaining
		}
	}

	bufferedWriter.WriteString("\n=== END goroutine stack dump ===\n")
	err = bufferedWriter.Flush()
	if err != nil {
		panic(err)
	}
}
