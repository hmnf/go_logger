package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	msg    chan<- string
	errors <-chan error
	File   io.ReadWriteCloser
}

func NewLogger() (*Logger, error) {
	file, err := os.OpenFile("./test.log", os.O_APPEND, 0777)

	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	fmt.Println(file)

	return &Logger{File: file}, nil
}

func (l *Logger) WritePut(msg string) {
	l.msg <- msg
}

func (l *Logger) WriteDelete(msg string) {
	l.msg <- msg
}

func (l *Logger) Err() <-chan error {
	return l.errors
}

func (l *Logger) Run() {
	msg := make(chan string, 16)
	l.msg = msg

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		defer close(msg)
		defer close(errors)

		for m := range msg {
			log.Println("Gorountine message: ", m)
		}

	}()

}
