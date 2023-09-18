package cmd

import (
	"testing"
	"time"
)

func TestShowHealth(t *testing.T) {
	if err := showHealth(); err != nil {
		t.Error(err)
	}
}

func TestEstopRecovery(t *testing.T) {
	c := make(chan string)
	go scan(c)
	go watcher(catRecovery)
	timeout := time.After(time.Second * 5)
	select {
	case <-c:
		t.Error("watcher failed")
	case <-timeout:
		t.Log("watcher passed")
	}
}
