package main

import (
	"fmt"
	"sync"
	"time"
)

type Scheduler struct {
	sync.Mutex
	started  bool
	function func() error
	ticker   *time.Ticker
}

func NewScheduler() *Scheduler {
	status := &Scheduler{}
	return status
}

func (s *Scheduler) Start(refreshInterval time.Duration, f func() error) {
	s.Lock()
	defer s.Unlock()

	if s.started {
		return
	}
	defer func() { s.started = true }()

	s.ticker = time.NewTicker(refreshInterval)
	s.function = f

	go func() {
		for {
			select {
			case <-s.ticker.C:
				if !func() bool {
					s.Lock()
					defer s.Unlock()
					return s.started
				}() {
					return
				}

				err := s.function()
				if err != nil {
					fmt.Println(err)
				}
			default:
				// stop goroutine if stopped
				if !func() bool {
					s.Lock()
					defer s.Unlock()
					return s.started
				}() {
					fmt.Println("stopped ticker.")
					return
				}

			}
		}

	}()
}

func (s *Scheduler) Stop() {
	s.Lock()
	defer s.Unlock()

	if !s.started {
		return
	}

	defer func() { s.started = false }()

	s.ticker.Stop()
}
