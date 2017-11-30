package janitor

import (
	"log"
	"time"
)

type Janitor struct {
	Name string
	fn   func() error
}

func NewJanitor(name string, fn func() error) *Janitor {
	return &Janitor{
		Name: name,
		fn:   fn,
	}
}

func (j *Janitor) Run() error {
	log.Printf("[DEBUG] %s Janitor: Starting Run", j.Name)
	return j.fn()
}

func (j *Janitor) RunEvery(d time.Duration) *time.Ticker {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			if err := j.Run(); err != nil {
				log.Printf("[ERROR] %s Janitor: %v", j.Name, err)
			}
		}
	}()

	return ticker
}
