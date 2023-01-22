package logger

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
)

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

const LogTemplate = "%d\t%s\t%s\n"

type Event struct {
	Key       string
	Value     string
	EventType EventType
}

type Logger struct {
	events chan<- Event
	errors <-chan error
	file   io.ReadWriteCloser
}

func NewLogger(fileName string) (*Logger, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND, 0777)

	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	fmt.Println(file)

	return &Logger{file: file}, nil
}

func (l *Logger) WritePut(key, value string) {
	l.events <- Event{Key: key, Value: value, EventType: EventPut}
}

func (l *Logger) WriteDelete(key string) {
	l.events <- Event{Key: key, EventType: EventDelete}
}

func (l *Logger) Err() <-chan error {
	return l.errors
}

func (l *Logger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		defer close(events)
		defer close(errors)

		for event := range events {
			_, err := fmt.Fprintf(l.file, LogTemplate, event.EventType, event.Key, url.QueryEscape(event.Value))

			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *Logger) ReadEvents() (<-chan Event, <-chan error) {
	outEvents := make(chan Event, 16)
	outErrors := make(chan error, 1)
	scanner := bufio.NewScanner(l.file)

	go func() {
		var event Event

		defer close(outEvents)
		defer close(outErrors)

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Sscanf(line, LogTemplate, &event.EventType, &event.Key, &event.Value)

			uv, err := url.QueryUnescape(event.Value)
			if err != nil {
				outErrors <- err
				return
			}

			event.Value = uv
			outEvents <- event
		}

		if err := scanner.Err(); err != nil {
			outErrors <- err
			return
		}

	}()

	return outEvents, outErrors
}
