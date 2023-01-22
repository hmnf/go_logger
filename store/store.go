package store

import (
	"bufio"
	"fmt"
	"net/url"
	"sync"

	"github.com/hmnf/go_logger/logger"
)

type Event struct {
	Key   string
	Value string
}

type Store struct {
	values map[string]string
	events chan Event
	l      *logger.Logger
	errors chan error
	sync.RWMutex
}

func NewStorageService() (*Store, error) {
	l, err := logger.NewLogger()

	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return &Store{
		values: make(map[string]string),
		l:      l,
	}, nil
}

func (s *Store) Put(key, value string) {
	s.Lock()
	fmt.Fprintf(s.l.File, "%s %s %s\n", "Put", key, value)
	s.values[key] = value
	fmt.Println(s.values)
	s.Unlock()
}

func (s *Store) Get(key string) string {
	s.RLock()
	defer s.RUnlock()
	return s.values[key]
}

func (s *Store) Delete(key string) {
	s.Lock()
	fmt.Fprintf(s.l.File, "%s %s %v\n", "Delete", key, nil)
	fmt.Println(s.values)
	delete(s.values, key)
	fmt.Println(s.values)
	s.Unlock()
}

func (s *Store) Restore() {
	scanner := bufio.NewScanner(s.l.File)

	for scanner.Scan() {
		line := scanner.Text()
		var method, key, value string
		fmt.Sscanf(line, "%s %s %s\n", &method, &key, &value)
		value, _ = url.QueryUnescape(value)
		switch method {
		case "Put":
			s.Put(key, value)
		case "Delete":
			s.Delete(key)
		default:
			fmt.Println("error")
		}

	}
}

func (s *Store) Run() {
	events := make(chan Event, 16)
	errors := make(chan error, 1)

	s.events = events
	s.errors = errors

	go func() {

	}()
}
